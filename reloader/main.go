package ssentr

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"gopkg.in/fsnotify.v1"

	"github.com/alexandrevicenzi/go-sse"
	"go.uber.org/zap"
)

func init() {
	caddy.RegisterModule(Middleware{})
	httpcaddyfile.RegisterHandlerDirective("reloader", parseCaddyfile)
}

// Middleware implements an HTTP handler that serves
// an SSE notification on reloads
type Middleware struct {
	watcher *fsnotify.Watcher
	listMap map[string]bool
	dirsMap map[string]bool
	w       io.Writer
	log     *zap.Logger
	s       *sse.Server
}

// CaddyModule returns the Caddy module information.
func (Middleware) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.reloader",
		New: func() caddy.Module { return new(Middleware) },
	}
}

// Provision implements caddy.Provisioner.
func (m *Middleware) Provision(ctx caddy.Context) error {
	m.log = ctx.Logger(m)
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		m.log.Fatal("cannot start watcher", zap.Error(err))
	}
	m.watcher = watcher

	m.listMap = make(map[string]bool)
	m.dirsMap = make(map[string]bool)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		n := scanner.Text()
		if fn, err := filepath.Abs(n); err == nil {
			m.listMap[fn] = true
			m.dirsMap[filepath.Dir(fn)] = true
		}
	}

	for template, _ := range m.dirsMap {
		err = watcher.Add(template)
		if err != nil {
			m.log.Fatal("cannot add template", zap.Error(err))
		}
	}
	m.s = sse.NewServer(&sse.Options{
		Logger: zap.NewStdLog(m.log),
	})
	m.log.Info("init", zap.Reflect("list", m.listMap))
	m.log.Info("init", zap.Reflect("dirs", m.dirsMap))
	return nil
}

// Validate implements caddy.Validator.
func (m *Middleware) Validate() error {
	if m.watcher == nil {
		return fmt.Errorf("no watcher")
	}
	if m.s == nil {
		return fmt.Errorf("no sse server")
	}
	return nil
}

// ServeHTTP implements caddyhttp.MiddlewareHandler.
func (m Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	m.log.Debug("setting up SSE for", zap.String("RequestURI", r.RequestURI))
	notify := m.watchFiles()
	done := make(chan bool)
	go func() {
		defer m.s.CloseChannel(r.RequestURI)
	L:
		for {
			select {
			case event := <-notify:
				m.log.Debug(event)
				m.s.SendMessage(r.RequestURI, sse.SimpleMessage(event))
			case <-done:
				break L
			}
		}
	}()
	m.s.ServeHTTP(w, r)
	m.log.Debug("closing SSE for", zap.String("RequestURI", r.RequestURI))
	done <- true
	return nil
}

// parseCaddyfile unmarshals tokens from h into a new Middleware.
func parseCaddyfile(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	var m Middleware
	return m, nil
}

// Interface guards
var (
	_ caddy.Provisioner           = (*Middleware)(nil)
	_ caddy.Validator             = (*Middleware)(nil)
	_ caddyhttp.MiddlewareHandler = (*Middleware)(nil)
)

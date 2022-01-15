package ssentr

import (
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
)

func glob(dir string, ext string) ([]string, error) {
	files := []string{}
	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if filepath.Ext(path) == ext {
			files = append(files, path)
		}
		return nil
	})

	return files, err
}
func (m Middleware) watchFiles() chan string {

	ch := make(chan string)
	go func() {
		for event := range m.watcher.Events {
		L:
			for {
				m.log.Debug("received", zap.Any("event", event))
				select {
				case event = <-m.watcher.Events:
				case err, ok := <-m.watcher.Errors:
					if !ok {
						return
					}
					m.log.Fatal("fsnotify had error", zap.Error(err))

				// Debounce and hack to deal with how some editors remove
				// replace files instead of only a write
				case <-time.After(100 * time.Millisecond):
					if _, exists := m.listMap[event.Name]; exists {
						ch <- event.Name
					}
					break L
				}
			}
		}
		close(ch)
	}()
	return ch
}

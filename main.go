package main

import (
	caddycmd "github.com/caddyserver/caddy/v2/cmd"

	// plug in Caddy modules here
	_ "github.com/caddyserver/caddy/v2/modules/standard"
	_ "github.com/greenpau/caddy-security"
	_ "github.com/greenpau/caddy-trace"
	_ "github.com/tomberek/ssentr/reloader"
)

func main() {
	caddycmd.Main()
}

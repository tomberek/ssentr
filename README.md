# ssentr

Send SSE when files change.

Setup:
```
$ cat Caddyfile
http://localhost:8085
route /reload {
	header Access-Control-Allow-Origin *
	header Access-Control-Request-Method GET
    reloader
}
```

Usage:
```
find -iname '*.html' | ssentr run
```

Add the following hyperscript to your project
```
<script type="text/hyperscript">
eventsource Reloader from http://localhost:8085/reload
    on message
        call location.reload()
    end
    on error
        log "error with connecting to reload SSE"
    end
    on close
        log "closing reload SSE"
    end
end
</script>
```

When any of the files change, a message will be sent to any connection to http://localhost:8085/reloader

# Credit
Based on and is an extension of [Caddy](https://caddyserver.com/)
Inspired by [entr(1)](https://eradman.com/entrproject/)
Used along with [hyperscript](https://hyperscript.org/)

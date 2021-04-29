package main

import (
	"flag"
	"log"
	"strings"

	"github.com/juztin/statictls/pkg/auth"
	"github.com/juztin/statictls/pkg/session"
	"github.com/juztin/statictls/web/app"
)

func main() {
	var (
		contentPath      = flag.String("content", "content", "path to static content to serve")
		cachePath        = flag.String("cache", ".cache", "path to TLS cert cache")
		authTemplatePath = flag.String("login", "", "path to login template")
		usersPath        = flag.String("users", "users.json", "path to users data")
		hostsString      = flag.String("hosts", "localhost", "hosts for autocert tls, comma separated")
		httpAddr         = flag.String("httpAddr", ":8000", "http address to listen on")
		tlsAddr          = flag.String("tlsAddr", ":8443", "tls address to listen on")
	)
	flag.Parse()

	a := auth.NewJson(*usersPath)
	s := session.NewMemory()
	app := app.New(s, a, *contentPath, *cachePath, *authTemplatePath, strings.Split(*hostsString, ",")...)
	log.Fatalln(app.Serve(*httpAddr, *tlsAddr))
}

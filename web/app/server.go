package app

import (
	"log"
	"net/http"

	"github.com/juztin/statictls/pkg/auth"
	"github.com/juztin/statictls/pkg/session"
	"github.com/juztin/statictls/web/middleware"
	"golang.org/x/crypto/acme/autocert"
)

type Server struct {
	*http.ServeMux
	session          session.Manager
	auth             auth.Authenticator
	authTemplatePath string
	cachePath        string
	contentPath      string
	hosts            []string
}

type loginParam struct {
	Username string
	HasError bool
	Message  string
}

func newSessionCookie(session string) *http.Cookie {
	return &http.Cookie{
		Name:     "session",
		Value:    session,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
}

func redirectInsecure(w http.ResponseWriter, r *http.Request) {
	target := "https://" + r.Host + r.URL.Path
	if len(r.URL.RawQuery) > 0 {
		target += "?" + r.URL.RawQuery
	}
	log.Printf("app: redirected insecure: %s\n", target)
	http.Redirect(w, r, target, http.StatusPermanentRedirect)
}

func (s *Server) Serve(httpAddr, tlsAddr string) error {
	fs := http.FileServer(http.Dir(s.contentPath))
	s.Handle("/user", http.HandlerFunc(authHandler(s.auth, s.session, s.authTemplatePath)))
	s.Handle("/", middleware.Auth(s.session, middleware.NoCache(fs)))
	if len(s.hosts) == 1 && s.hosts[0] == "localhost" {
		log.Printf("Listening on %s...\n", httpAddr)
		return http.ListenAndServe(httpAddr, s)
	} else {
		log.Printf("acme hosts: %s\n", s.hosts)
		m := &autocert.Manager{
			Cache:      autocert.DirCache(s.cachePath),
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(s.hosts...),
		}
		go func() {
			log.Printf("Listening on %s\n", httpAddr)
			log.Fatal(http.ListenAndServe(httpAddr, m.HTTPHandler(http.HandlerFunc(redirectInsecure))))
		}()
		svr := &http.Server{
			Addr:      tlsAddr,
			TLSConfig: m.TLSConfig(),
		}
		svr.Handler = s
		log.Printf("Listening on %s\n", tlsAddr)
		return (svr.ListenAndServeTLS("", ""))
	}
	return nil
}

func New(s session.Manager, a auth.Authenticator, contentPath, cachePath, authTemplatePath string, hosts ...string) *Server {
	return &Server{
		ServeMux:         http.NewServeMux(),
		session:          s,
		auth:             a,
		authTemplatePath: authTemplatePath,
		cachePath:        cachePath,
		contentPath:      contentPath,
		hosts:            hosts,
	}
}

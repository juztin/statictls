package app

import (
	"log"
	"net/http"
	"time"

	"github.com/juztin/statictls/internal/middleware"
	"github.com/juztin/statictls/pkg/auth"
	"github.com/juztin/statictls/pkg/session"
	"golang.org/x/crypto/acme/autocert"
)

type Server struct {
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

func (s *Server) Serve() error {
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir(s.contentPath))
	mux.Handle("/user", http.HandlerFunc(authHandler(s.auth, s.session, s.authTemplatePath)))
	mux.Handle("/", middleware.Auth(s.session, middleware.NoCache(fs)))

	go func() {
		for {
			<-time.After(5 * time.Minute)
			s.session.Remove(15 * time.Minute)
		}
	}()

	if len(s.hosts) == 1 && s.hosts[0] == "localhost" {
		log.Println("Listening on :3000...")
		return http.ListenAndServe(":3000", mux)
	} else {
		go func() {
			log.Println("Listening on :80")
			log.Fatal(http.ListenAndServe(":80", http.HandlerFunc(redirectInsecure)))
		}()
		log.Println("Listening on :443")
		m := &autocert.Manager{
			Cache:      autocert.DirCache(s.cachePath),
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(s.hosts...),
		}
		s := &http.Server{
			Addr:      ":https",
			TLSConfig: m.TLSConfig(),
		}
		s.Handler = mux
		return (s.ListenAndServeTLS("", ""))
	}
	return nil
}

func New(s session.Manager, a auth.Authenticator, contentPath, cachePath, authTemplatePath string, hosts []string) *Server {
	return &Server{
		session:          s,
		auth:             a,
		authTemplatePath: authTemplatePath,
		cachePath:        cachePath,
		contentPath:      contentPath,
		hosts:            hosts,
	}
}

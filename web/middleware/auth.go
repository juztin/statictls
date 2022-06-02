package middleware

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/juztin/statictls/pkg/session"
)

func Auth(s session.Manager, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isValid, err := s.Check(r)
		if err != nil {
			log.Printf("middleware: failed session check: %q\n", err)
			w.WriteHeader(http.StatusInternalServerError)
		} else if isValid == false {
			newURL := fmt.Sprintf("/user?redirect=%s", url.QueryEscape(r.URL.Path))
			http.Redirect(w, r, newURL, http.StatusSeeOther)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

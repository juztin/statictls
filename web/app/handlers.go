package app

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/juztin/statictls/pkg/auth"
	"github.com/juztin/statictls/pkg/session"
)

const login = `
<!doctype html>
<html>
<head>
<style>
* {
    box-sizing: border-box;
    margin: 0;
    padding: 0;
	font-family: 'Courier';
}
body {
	height: 100vh;
	background: rgb(20,20,40);
	background: linear-gradient(to bottom right,
								rgb(20, 20, 20) 0%,
								rgb(50, 50, 50) 100%);
}
body
body > div {
	position: absolute;
	top: 50%;
	left: 0;
	width: 100%;
	height: 400px;
	margin-top: -200px;
	overflow: hidden;
}
body > div > div {
	max-width: 600px;
	margin: 0 auto;
	padding: 80px 0;
	height: 400px;
	text-align: center;
}
body > div > div > h1 {
	color: rgb(240, 240, 240);
	font-size: 40px;
	font-weight: 200;
}

form input {
    -webkit-appearance: none;
       -moz-appearance: none;
            appearance: none;
    outline: 0;
    border: 1px solid rgba(255, 255, 255, 0.4);
    background-color: rgba(255, 255, 255, 0.2);
    width: 250px;
    border-radius: 3px;
    padding: 10px 15px;
    margin: 0 auto 10px auto;
    display: block;
    text-align: center;
    font-size: 18px;
    color: white;
    transition-duration: 0.25s;
    font-weight: 300;
}
form input:hover {
    background-color: rgba(255, 255, 255, 0.4);
}
form button {
    -webkit-appearance: none;
       -moz-appearance: none;
            appearance: none;
    outline: 0;
    background-color: white;
    border: 0;
    padding: 10px 15px;
    color: rgb(20,20,20);
    border-radius: 3px;
    width: 250px;
    cursor: pointer;
    font-size: 18px;
    transition-duration: 0.25s;
}
form button:hover {
    background-color: rgb(220, 220, 220);
}

form + div {
	padding: 15px 0;
	color: rgb(255, 166, 166);
}
form.error > input {
	border-color: rgb(255, 166, 166);
}

</style>
</head>
<body>
	<div>
		<div>
			<h1>Login</h1>
			<form{{if .HasError}} class="error"{{end}} action="/user" method="POST">
				<input name="username" placeholder="username" value="{{.Username}}">
				<input name="password" type="password" placeholder="password">
				<button type="submit">Login</button>
			</form>
			<div>{{.Message}}</div>
		</div>
	</div>
</body>
</html>`

var defaultAuthTemplate = template.Must(template.New("login").Parse(login))

func authHandler(a auth.Authenticator, s session.Manager, authTemplatePath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var l loginParam
		if r.Method == http.MethodPost {
			err := a.Authenticate(r.FormValue("username"), r.FormValue("password"))
			if err == nil {
				log.Printf("authHandler: authenticated user %q\n", r.FormValue("username"))
				newURL := r.URL.Query().Get("redirect")
				if newURL == "" {
					newURL = "/"
				}
				session, err := s.New()
				if err != nil {
					log.Printf("authHandler: failed to create session: %q\n", err)
					w.WriteHeader(http.StatusInternalServerError)
				} else {
					http.SetCookie(w, newSessionCookie(session))
					http.Redirect(w, r, newURL, http.StatusSeeOther)
				}
				return
			}
			log.Printf("authHandler: authentication failure from %s: %q\n", r.RemoteAddr, err)
			l = loginParam{
				Username: r.FormValue("username"),
				HasError: true,
				Message:  "Invalid username/password?",
			}
		}

		var authTemplate *template.Template
		if authTemplatePath == "" {
			authTemplate = defaultAuthTemplate
		} else if _, err := os.Stat(authTemplatePath); os.IsNotExist(err) {
			log.Printf("authHandler: login template missing %q\n", authTemplatePath)
			authTemplate = defaultAuthTemplate
		} else if authTemplate, err = template.ParseFiles(authTemplatePath); err != nil {
			log.Printf("authHandler: template parse failure: %q\n", err)
			authTemplate = defaultAuthTemplate
		}
		err := authTemplate.Execute(w, l)
		if err != nil {
			log.Printf("authHandler: template execution failure: %q\n", err)
			log.Printf("authHandler: %q\n", err)
		}
	}
}

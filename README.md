# statictls

Simple static site server with auto TLS redirection and authentication


Install:

```
go get github.com/juztin/statictls/cmd/statictls
```

Run:

```
statictls
```

Params:

```
  -cache string
    	path to TLS cert cache (default ".cache/")
  -content string
    	path to static content to serve (default "content/")
  -hosts string
    	hosts for autocert tls, comma separated (default "localhost")
  -login string
    	path to login template
  -users string
    	path to users data (default "users.json")
```

_If no hosts are specified, or `localhost` is used, it will only serve on http, with not auto TLS redirection._  
_If no login template is specified, a default one is used._

#### Users

For authenticaion, you'll need a `users.json` file.  
It should be structed like:

```json
{
  "username": "bcrypt password"
}
```

To quickly generate a `bcrypt` hash:

```go
package main

import (
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	password := os.Args[1]
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(string(hash))
}
```

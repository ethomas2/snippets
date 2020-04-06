package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"
)

func main() {

	var port int
	var cors bool
	flag.IntVar(&port, "port", 8000, "The port to listen on (default 8000)")
	flag.BoolVar(&cors, "cors", false, "Whether to respond to CORS headers")
	flag.Parse()

	s := http.Server{
		Addr: fmt.Sprintf(":%d", port),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if cors && strings.ToLower(r.Method) == "options" {
				origins := r.Header["Origin"]
				var origin string
				// try to avoid setting origin to * directly bc
				// chrome complains about setting origin to *
				// if request.credentials has been set to
				// 'include'
				// https://developer.mozilla.org/en-US/docs/Web/API/Request/credentials
				if len(origins) > 0 {
					origin = origins[0]
				} else {
					origin = "*"
				}
				w.Header().Set("Access-Control-Allow-Origin", origin)
				// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Credentials
				w.Header().Set("Access-Control-Allow-Credentials", "true")

				reqHeaders := r.Header["Access-Control-Request-Headers"]
				if len(reqHeaders) > 0 {
					w.Header().Set(
						"Access-Control-Allow-Headers",
						strings.Join(reqHeaders, ","))
				}
			}

			bytes, err := httputil.DumpRequest(r, true)
			if err != nil {
				fmt.Println("\033[31moh no, error\033[m")
			}
			fmt.Println(string(bytes))
		}),
	}
	fmt.Println(s.ListenAndServe())
}

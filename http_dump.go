package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/http/httputil"
)

func main() {

	var port int
	flag.IntVar(&port, "port", 8000, "The port to listen on (default 8000)")
	flag.Parse()

	s := http.Server{
		Addr: fmt.Sprintf(":%d", port),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			bytes, err := httputil.DumpRequest(r, true)
			if err != nil {
				fmt.Println("\033[31moh no, error\033[m")
			}
			fmt.Println(string(bytes))
		}),
	}
	fmt.Println(s.ListenAndServe())
}

package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
)

func main() {

	s := http.Server{
		Addr: ":8000",
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

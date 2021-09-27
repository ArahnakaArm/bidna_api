package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/apix/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World")
	})

	http.HandleFunc("/apix/greet/", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Path[len("/greet/"):]
		fmt.Fprintf(w, "Hello %s\n", name)
	})

	http.ListenAndServe(":3334", nil)
}

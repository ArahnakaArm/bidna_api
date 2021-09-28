package routes

import (
	"fmt"
	"net/http"
)

func AddCommonRoute() {
	http.HandleFunc("/apix/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World")
	})
}

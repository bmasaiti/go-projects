package checker

import (
	"fmt"
	"html"
	"net/http"
)

func main() {

	fmt.Println("Startinf Server...")
	router := http.NewServeMux()

	//basic routing
	router.HandleFunc("users/{user}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello %s", html.EscapeString(r.PathValue("user")))
	})

	v2 := http.NewServeMux()
	v2.HandleFunc("/users/{user}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hell v2 %s", html.EscapeString(r.PathValue("user")))
	})

	router.Handle("/v2", http.StripPrefix("/v2",v2))
	server := http.Server{
		Addr:    ":8000",
		Handler: router,
	}
	server.ListenAndServe()

}

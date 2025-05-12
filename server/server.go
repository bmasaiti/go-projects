package main

import (
	//"crypto/tls"
	"encoding/json"
	"fmt"
	"html"
	"net/http"
)

type Site struct {
	Id      int      `json:"id"`
	Url     string   `json:"url"`
	Headers []string `json:"headers"`
	Method  string   `json:"method"`
}

func main() {

	fmt.Println("Startinf Server...")
	router := http.NewServeMux()

	//basic routing
	router.HandleFunc("/users/{user}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello  basic %s", html.EscapeString(r.PathValue("user")))
	})

	//method bases routing takes precedense over the basic one fyi the method requires space btwn method and path
	router.HandleFunc("GET /users/{user}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello method based %s", html.EscapeString(r.PathValue("user")))
	})

	//host based routing
	router.HandleFunc("GET /some.host.com", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello method based %s", html.EscapeString(r.PathValue("user")))
	})

	router.HandleFunc("POST /sites/upload",
		func(w http.ResponseWriter, r *http.Request) {

			var sites Site
			err := json.NewDecoder(r.Body).Decode(&sites)
			defer r.Body.Close()
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			fmt.Printf("Printing received sites:%+v", sites)

		})

	v2 := http.NewServeMux()
	v2.HandleFunc("/users/{user}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hell v2 %s", html.EscapeString(r.PathValue("user")))
	})

	//router.Handle("/v2", http.StripPrefix("/v2", v2))
	server := http.Server{
		Addr:    ":8000",
		Handler: router,
	}

	//certFilePath := "path to the certificate store "
	// keyFilePath := "again similar to above"

	//serverTLScert ,err := tls.LoadX509KeyPair(certFilePath,keyFilePath)

	// if err!=nil {
	// 	fmt.Errorf("Failed to validate cert")
	// }

	server.ListenAndServe()
	// tlsConfig := &tls.Config{
	// Certificates: []tls.Certificate{serverTLScert},
}

// tls setup (ListenandServeTLS + the port is now 443)
// server2 := http.Server{
// 	Addr:    ":80",
// 	Handler: router,
// }

//server2.ListenAndServeTLS()

//}

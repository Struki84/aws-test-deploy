package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/GoLangWebSDK/rest"
)

var (
	MainHost = os.Getenv("HOST")
	WSHost   = os.Getenv("WS_HOST")
)

func main() {
	if MainHost == "" {
		MainHost = "http://localhost"
	}

	if WSHost == "" {
		WSHost = "http://ws-app:8080"
	}
	router := rest.NewRouter()
	ctrl := rest.NewController(router)

	// Serve static files
	fs := http.FileServer(http.Dir("static"))
	router.Mux.Handle("/static/", http.StripPrefix("/static/", fs))

	router.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	})

	ctrl.Get("/health", func(session *rest.Session) {
		session.Response.WriteHeader(200)
		session.Response.Write([]byte("Main container is up..."))
	})

	ctrl.Get("/health/ws", func(session *rest.Session) {
		reqUrl := WSHost + "/health"

		req, err := http.NewRequest("GET", reqUrl, nil)
		if err != nil {
			fmt.Println(err)
			session.SetStatus(http.StatusInternalServerError)
			return
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println(err)
			session.Response.Write([]byte("Failed to reach WS container via GET request to " + reqUrl))
			session.SetStatus(http.StatusInternalServerError)
			return
		}

		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			fmt.Println(err)
			session.Response.Write([]byte("Failed to reach WS container via GET request to " + reqUrl))
			return
		}

		session.Response.Write([]byte("WS container is up, tested via GET request to " + reqUrl))
	})

	port := ":" + os.Getenv("PORT")
	if port == ":" {
		port = ":80"
	}

	fmt.Println("Listening on port " + port)
	http.ListenAndServe(port, router)
}

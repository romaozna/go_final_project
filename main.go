package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"log"
	"main/src/store"
	"net/http"
	"os"
	"strconv"

	_ "modernc.org/sqlite"
)

func getPort() int {
	port := 7540
	envPort := os.Getenv("TODO_PORT")
	if len(envPort) > 0 {
		if pport, err := strconv.ParseInt(envPort, 10, 32); err == nil {
			port = int(pport)
		}
	}

	return port
}

func main() {
	store.CreateDatabase()
	webDir := "./web"

	r := chi.NewRouter()
	r.Mount("/", http.FileServer(http.Dir(webDir)))

	serverPort := getPort()
	log.Println(fmt.Sprintf("Адрес сервера: %d", serverPort))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", serverPort), r))
}

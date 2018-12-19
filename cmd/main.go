package main

import (
	"flag"
	web "github.com/oktadeveloper/okta-go-vue-example/pkg/http"
	"github.com/oktadeveloper/okta-go-vue-example/pkg/storage"
	"log"
	"net/http"
	"os"
)

func main() {
	httpPort := os.Getenv("PORT")
	if httpPort == "" {
		httpPort = ":4444"
	} else {
		httpPort = ":" + httpPort
	}

	flag.StringVar(&httpPort, "b", httpPort, "bind on port")
	flag.Parse()

	repo := storage.NewMongoRepository()
	webService := web.New(repo)

	log.Printf("Running on port %s\n", httpPort)
	log.Fatal(http.ListenAndServe(httpPort, webService.Router))
}

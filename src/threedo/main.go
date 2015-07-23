package main

import (
	"flag"
	"fmt"
	"lib"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	_ "lib/api/v1"
)

var (
	WEB_PORT      = flag.Uint64("port", 8080, "Port for the webserver to listen on")
	FRONTEND_PATH = flag.String("frontend", "frontend", "Path to frontend files")
	DEBUG         = flag.Bool("debug", false, "Debug mode")
	DB_SEED       = flag.Bool("db-seed", false, "Seed the DB with initial values")
	DB_URI        = flag.String("db-uri", ":memory:", "Path/URI to store DB at")
	FRONTEND_URL  = flag.String(
		"frontend-url",
		"http://localhost:4200",
		"In debug mode reverse proxy is used instead of static file server",
	)
)

func initDB() {
	log.Printf("Initialising database at %s...\n", *DB_URI)
	db := lib.InitDB(*DB_URI)
	lib.CreateDBSchema(db)

	if *DB_SEED {
		log.Println("Seeding database...")
		if err := lib.SeedDB(db); err != nil {
			log.Fatalln(err)
		}
	}
}

func addStaticRoute() {
	if *DEBUG {
		u, err := url.Parse(*FRONTEND_URL)
		if err != nil {
			log.Fatal("Couldn't parse frontend URL")
		}
		lib.Routes.PathPrefix("/").Handler(httputil.NewSingleHostReverseProxy(u))
	} else {
		lib.Routes.PathPrefix("/").Handler(http.FileServer(http.Dir(*FRONTEND_PATH)))
	}
}

func startHTTPServer() {
	if *DEBUG {
		log.Printf("Frontend located at %s\n", *FRONTEND_URL)
	} else {
		log.Printf("Frontend located at %s\n", *FRONTEND_PATH)
	}
	log.Printf("Listening on :%d\n", *WEB_PORT)

	routes := lib.GetRouteHandler()
	addr := fmt.Sprintf(":%d", *WEB_PORT)
	if err := http.ListenAndServe(addr, routes); err != nil {
		log.Fatal(err)
	}
}

func main() {
	flag.Parse()

	initDB()
	addStaticRoute()
	startHTTPServer()
}

package main

import (
	"fmt"
	"github.com/3onyc/3do/api/v1"
	"github.com/3onyc/3do/appinit"
	"github.com/3onyc/3do/util"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/namsral/flag"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	_ "github.com/3onyc/3do/api/v1"
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

func initDB() *sqlx.DB {
	log.Printf("Initialising database at %s...\n", *DB_URI)
	db := sqlx.MustConnect("sqlite", *DB_URI)
	appinit.CreateDBSchema(db)

	return db
}

func seedDB(ctx *util.Context) {
	if *DB_SEED {
		log.Println("Seeding database...")
		if err := appinit.SeedDB(ctx); err != nil {
			log.Fatalln(err)
		}
	}
}

func addStaticRoute(router *mux.Router) {
	if *DEBUG {
		u, err := url.Parse(*FRONTEND_URL)
		if err != nil {
			log.Fatal("Couldn't parse frontend URL")
		}
		router.PathPrefix("/").Handler(httputil.NewSingleHostReverseProxy(u))
	} else {
		router.PathPrefix("/").Handler(http.FileServer(http.Dir(*FRONTEND_PATH)))
	}
}

func initRoutes(ctx *util.Context) {
	v1.NewListsAPI(ctx).Register()
	v1.NewGroupsAPI(ctx).Register()
	v1.NewItemsAPI(ctx).Register()
}

func startHTTPServer(router *mux.Router) {
	if *DEBUG {
		log.Printf("Frontend located at %s\n", *FRONTEND_URL)
	} else {
		log.Printf("Frontend located at %s\n", *FRONTEND_PATH)
	}
	log.Printf("Listening on :%d\n", *WEB_PORT)

	addr := fmt.Sprintf(":%d", *WEB_PORT)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal(err)
	}
}

func main() {
	flag.Parse()

	router := mux.NewRouter()
	db := initDB()

	ctx := util.NewContext(router, db)
	initRoutes(ctx)
	seedDB(ctx)

	addStaticRoute(router)
	startHTTPServer(router)
}

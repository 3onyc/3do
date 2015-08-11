package main

import (
	"fmt"
	"github.com/3onyc/3do/api/v1"
	"github.com/3onyc/3do/appinit"
	"github.com/3onyc/3do/middleware"
	"github.com/3onyc/3do/module"
	"github.com/3onyc/3do/util"
	"github.com/GeertJohan/go.rice"
	"github.com/codegangsta/negroni"
	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/namsral/flag"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var (
	CONFIG_FILE  = flag.String("config", "", "Path to config file")
	WEB_PORT     = flag.Uint64("port", 8080, "Port for the webserver to listen on")
	DEBUG        = flag.Bool("debug", false, "Debug mode")
	DB_SEED      = flag.Bool("db-seed", false, "Seed the DB with initial values")
	DB_URI       = flag.String("db-uri", ":memory:", "Path/URI to store DB at")
	FRONTEND_URL = flag.String(
		"frontend-url",
		"http://localhost:4200",
		"In debug mode reverse proxy is used instead of embedded files",
	)
)

func initDB(l log.Logger) *sqlx.DB {
	l.Log(
		"action", "init-db",
		"driver", "sqlite3",
		"dburi", *DB_URI,
	)

	db := sqlx.MustConnect("sqlite3", *DB_URI)
	appinit.CreateDBSchema(db)

	return db
}

func seedDB(ctx *util.Context) {
	if *DB_SEED {
		if err := appinit.SeedDB(ctx); err != nil {
			ctx.Log(
				"action", "seed-db",
				"result", "false",
				"message", err,
			)
		}
	}
}

func addStaticRoute(ctx *util.Context) {
	if *DEBUG {
		u, err := url.Parse(*FRONTEND_URL)
		if err != nil {
			ctx.Log(
				"action", "add-static-route",
				"result", false,
				"message", "Couldn't parse frontend URL",
			)
		}
		ctx.Router.PathPrefix("/").Handler(httputil.NewSingleHostReverseProxy(u))
	} else {
		ctx.Router.PathPrefix("/").Handler(
			util.CatchAllFileServer(
				rice.MustFindBox("frontend/dist").HTTPBox(),
				"/index.html",
				ctx.Logger,
			),
		)
	}
}

func initRoutes(ctx *util.Context) {
	v1.NewListsAPI(ctx).Register()
	v1.NewGroupsAPI(ctx).Register()
	v1.NewItemsAPI(ctx).Register()
}

func initModules(ctx *util.Context) {
	if err := module.NewHabitica(ctx).Init(); err != nil {
		ctx.Log(
			"action", "init-modules",
			"result", false,
			"message", err,
		)
	}
}

func startHTTPServer(ctx *util.Context) {
	if *DEBUG {
		ctx.Log(
			"action", "debug",
			"message", fmt.Sprintf("Frontend located at %s", *FRONTEND_URL),
		)
	}

	addr := fmt.Sprintf(":%d", *WEB_PORT)
	ctx.Log(
		"action", "listen",
		"address", addr,
	)

	n := negroni.New(middleware.NewLogger(ctx.Logger))
	n.UseHandler(ctx.Router)

	if err := http.ListenAndServe(addr, n); err != nil {
		ctx.Log(
			"action", "listen",
			"result", false,
			"message", err,
		)
	}
}

func main() {
	flag.Parse()

	log := log.NewJSONLogger(os.Stderr)

	router := mux.NewRouter()
	db := initDB(log)

	ctx := util.NewContext(router, db, log)
	initRoutes(ctx)
	initModules(ctx)
	seedDB(ctx)

	addStaticRoute(ctx)
	startHTTPServer(ctx)
}

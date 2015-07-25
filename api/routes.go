package api

import (
	"github.com/3onyc/threedo-backend/middleware"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"net/http"
)

var (
	Routes = mux.NewRouter()
)

func GetRouteHandler() http.Handler {
	n := negroni.New(
		middleware.NewLogger(),
		negroni.NewRecovery(),
	)

	n.UseHandler(Routes)
	return n
}

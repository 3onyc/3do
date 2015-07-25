package threedo

import (
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"net/http"
)

var (
	Routes = mux.NewRouter()
)

func GetRouteHandler() http.Handler {
	n := negroni.New(
		NewLogger(),
		negroni.NewRecovery(),
	)

	n.UseHandler(Routes)
	return n
}

// +build nobuild
package api1

import (
	"github.com/gorilla/mux"
	"lib"
	"lib/model"
	"lib/util"
	"net/http"
	"strconv"
)

type GroupListResponse struct {
	TodoGroups []model.TodoGroup `json:"todoGroups"`
}

type GroupGetResponse struct {
	TodoGroup model.TodoGroup `json:"todoGroup"`
}

func groupsList(w http.ResponseWriter, r *http.Request) {
	if groups, err := model.GetAllTodoGroups(lib.GetDB()); err != nil {
		http.Error(w, err.Error(), 500)
	} else {
		util.WriteJSONResponse(w, &GroupListResponse{
			TodoGroups: groups,
		})
	}
}

func groupGet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	if group, err := model.FindTodoGroup(lib.GetDB(), id); err != nil {
		http.Error(w, err.Error(), 500)
	} else if group == nil {
		http.Error(w, "Group not found", 404)
	} else {
		util.WriteJSONResponse(w, &GroupGetResponse{
			TodoGroup: *group,
		})
	}
}

func init() {
	lib.Routes.HandleFunc("/api/v1/todoGroups", groupsList).Methods("GET")
	lib.Routes.HandleFunc("/api/v1/todoGroups/{id}", groupGet).Methods("GET")
}

// +build nobuild
package api1

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"lib"
	"lib/model"
	"lib/util"
	"net/http"
	"strconv"
)

type GroupListResponse struct {
	TodoGroups []*model.TodoGroup `json:"todoGroups"`
}

type GroupGetResponse struct {
	TodoGroup *model.TodoGroup `json:"todoGroup"`
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
			TodoGroup: group,
		})
	}
}

func groupPut(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	group, err := model.FindTodoGroup(lib.GetDB(), id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if group == nil {
		http.Error(w, "Item not found", 404)
		return
	}

	var payload = &GroupGetResponse{}
	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding JSON (%s)", err), 500)
		return
	}

	group.Title = payload.TodoGroup.Title
	group.List = payload.TodoGroup.List

	if err := model.UpdateTodoGroup(lib.GetDB(), group); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func groupPost(w http.ResponseWriter, r *http.Request) {
	var payload = &GroupGetResponse{}
	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding JSON (%s)", err), 500)
		return
	}

	postItem := payload.TodoGroup
	_, err := model.InsertTodoGroup(lib.GetDB(), &model.TodoGroup{
		Title: postItem.Title,
		List:  postItem.List,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func init() {
	lib.Routes.HandleFunc("/api/v1/todoGroups", groupsList).Methods("GET")
	lib.Routes.HandleFunc("/api/v1/todoGroups/{id}", groupGet).Methods("GET")
	lib.Routes.HandleFunc("/api/v1/todoGroups/{id}", groupPut).Methods("PUT")
	lib.Routes.HandleFunc("/api/v1/todoGroups", groupPost).Methods("POST")
}

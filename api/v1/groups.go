package api1

import (
	"encoding/json"
	"fmt"
	"github.com/3onyc/threedo-backend/api"
	"github.com/3onyc/threedo-backend/model"
	"github.com/3onyc/threedo-backend/util"
	"github.com/gorilla/mux"
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
	if groups, err := model.GetAllTodoGroups(model.GetDB()); err != nil {
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

	if group, err := model.FindTodoGroup(model.GetDB(), id); err != nil {
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

	group, err := model.FindTodoGroup(model.GetDB(), id)
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

	if err := model.UpdateTodoGroup(model.GetDB(), group); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	util.WriteJSONResponse(w, &GroupGetResponse{
		TodoGroup: group,
	})
}

func groupPost(w http.ResponseWriter, r *http.Request) {
	var payload = &GroupGetResponse{}
	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding JSON (%s)", err), 500)
		return
	}

	postGroup := payload.TodoGroup
	group := &model.TodoGroup{
		Title: postGroup.Title,
		List:  postGroup.List,
	}

	if err := model.InsertTodoGroup(model.GetDB(), group); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	util.WriteJSONResponse(w, &GroupGetResponse{
		TodoGroup: group,
	})
}

func groupDelete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if err := model.DeleteTodoGroup(model.GetDB(), id); err == model.GroupNotFound {
		http.Error(w, err.Error(), 404)
		return
	} else if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func init() {
	api.Routes.HandleFunc("/api/v1/todoGroups", groupsList).Methods("GET")
	api.Routes.HandleFunc("/api/v1/todoGroups/{id}", groupGet).Methods("GET")
	api.Routes.HandleFunc("/api/v1/todoGroups/{id}", groupPut).Methods("PUT")
	api.Routes.HandleFunc("/api/v1/todoGroups/{id}", groupDelete).Methods("DELETE")
	api.Routes.HandleFunc("/api/v1/todoGroups", groupPost).Methods("POST")
}

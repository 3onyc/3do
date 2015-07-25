package v1

import (
	"encoding/json"
	"fmt"
	"github.com/3onyc/threedo-backend/model"
	"github.com/3onyc/threedo-backend/util"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type GroupsAPI struct {
	*util.Context
}

func NewGroupsAPI(ctx *util.Context) *GroupsAPI {
	return &GroupsAPI{ctx}
}

type GroupListResponse struct {
	TodoGroups []*model.TodoGroup `json:"todoGroups"`
}

type GroupGetResponse struct {
	TodoGroup *model.TodoGroup `json:"todoGroup"`
}

func (g *GroupsAPI) Register() {
	g.Router.HandleFunc("/api/v1/todoGroups", g.list).Methods("GET")
	g.Router.HandleFunc("/api/v1/todoGroups/{id}", g.get).Methods("GET")
	g.Router.HandleFunc("/api/v1/todoGroups/{id}", g.put).Methods("PUT")
	g.Router.HandleFunc("/api/v1/todoGroups/{id}", g.delete).Methods("DELETE")
	g.Router.HandleFunc("/api/v1/todoGroups", g.post).Methods("POST")
}

func (g *GroupsAPI) list(w http.ResponseWriter, r *http.Request) {
	if groups, err := model.GetAllTodoGroups(g.DB); err != nil {
		http.Error(w, err.Error(), 500)
	} else {
		util.WriteJSONResponse(w, &GroupListResponse{
			TodoGroups: groups,
		})
	}
}

func (g *GroupsAPI) get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	if group, err := model.FindTodoGroup(g.DB, id); err != nil {
		http.Error(w, err.Error(), 500)
	} else if group == nil {
		http.Error(w, "Group not found", 404)
	} else {
		util.WriteJSONResponse(w, &GroupGetResponse{
			TodoGroup: group,
		})
	}
}

func (g *GroupsAPI) put(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	group, err := model.FindTodoGroup(g.DB, id)
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

	if err := model.UpdateTodoGroup(g.DB, group); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	util.WriteJSONResponse(w, &GroupGetResponse{
		TodoGroup: group,
	})
}

func (g *GroupsAPI) post(w http.ResponseWriter, r *http.Request) {
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

	if err := model.InsertTodoGroup(g.DB, group); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	util.WriteJSONResponse(w, &GroupGetResponse{
		TodoGroup: group,
	})
}

func (g *GroupsAPI) delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if err := model.DeleteTodoGroup(g.DB, id); err == model.GroupNotFound {
		http.Error(w, err.Error(), 404)
		return
	} else if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

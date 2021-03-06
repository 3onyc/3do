package v1

import (
	"encoding/json"
	"fmt"
	"github.com/3onyc/3do/model"
	"github.com/3onyc/3do/util"
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
	TodoGroup *model.TodoGroup  `json:"todoGroup"`
	TodoItems []*model.TodoItem `json:"todoItem,omitempty"`
}

func (g *GroupsAPI) Register() {
	g.Router.HandleFunc("/api/v1/todoGroups", g.list).Methods("GET")
	g.Router.HandleFunc("/api/v1/todoGroups/{id}", g.get).Methods("GET")
	g.Router.HandleFunc("/api/v1/todoGroups/{id}", g.put).Methods("PUT")
	g.Router.HandleFunc("/api/v1/todoGroups/{id}", g.delete).Methods("DELETE")
	g.Router.HandleFunc("/api/v1/todoGroups", g.post).Methods("POST")
}

func (g *GroupsAPI) list(w http.ResponseWriter, r *http.Request) {
	if groups, err := g.Groups.GetAll(); err != nil {
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
		return
	}

	group, err := g.Groups.Find(id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	} else if group == nil {
		http.Error(w, "Group not found", 404)
		return
	}

	items, err := g.Items.GetAllForGroup(id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	util.WriteJSONResponse(w, &GroupGetResponse{
		TodoGroup: group,
		TodoItems: items,
	})
}

func (g *GroupsAPI) put(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	group, err := g.Groups.Find(id)
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
	group.ListID = payload.TodoGroup.ListID

	if err := g.Groups.Update(group); err != nil {
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
		Title:  postGroup.Title,
		ListID: postGroup.ListID,
	}

	if err := g.Groups.Insert(group); err != nil {
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

	if err := g.Groups.Delete(id); err == model.GroupNotFound {
		http.Error(w, err.Error(), 404)
		return
	} else if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

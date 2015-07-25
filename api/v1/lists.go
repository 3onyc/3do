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

type ListsAPI struct {
	*util.Context
}

func NewListsAPI(ctx *util.Context) *ListsAPI {
	return &ListsAPI{ctx}
}

func (l *ListsAPI) Register() {
	l.Router.HandleFunc("/api/v1/todoLists", l.list).Methods("GET")
	l.Router.HandleFunc("/api/v1/todoLists/{id}", l.get).Methods("GET")
	l.Router.HandleFunc("/api/v1/todoLists/{id}", l.put).Methods("PUT")
	l.Router.HandleFunc("/api/v1/todoLists/{id}", l.delete).Methods("DELETE")
	l.Router.HandleFunc("/api/v1/todoLists", l.post).Methods("POST")
}

type ListListResponse struct {
	TodoLists []*model.TodoList `json:"todoLists"`
}

type ListGetResponse struct {
	TodoList *model.TodoList `json:"todoList"`
}

func (l *ListsAPI) list(w http.ResponseWriter, r *http.Request) {
	if lists, err := l.Lists.GetAll(); err != nil {
		http.Error(w, err.Error(), 500)
	} else {
		util.WriteJSONResponse(w, &ListListResponse{
			TodoLists: lists,
		})
	}
}

func (l *ListsAPI) get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	if list, err := l.Lists.Find(id); err != nil {
		http.Error(w, err.Error(), 500)
	} else if list == nil {
		http.Error(w, "List not found", 404)
	} else {
		util.WriteJSONResponse(w, &ListGetResponse{
			TodoList: list,
		})
	}
}

func (l *ListsAPI) put(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	list, err := l.Lists.Find(id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if list == nil {
		http.Error(w, "List not found", 404)
		return
	}

	var payload = &ListGetResponse{}
	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding JSON (%s)", err), 500)
		return
	}

	list.Title = payload.TodoList.Title
	list.Description = payload.TodoList.Description

	if err := l.Lists.Update(list); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	util.WriteJSONResponse(w, &ListGetResponse{
		TodoList: list,
	})
}

func (l *ListsAPI) post(w http.ResponseWriter, r *http.Request) {
	var payload = &ListGetResponse{}
	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding JSON (%s)", err), 500)
		return
	}

	postList := payload.TodoList
	list := &model.TodoList{
		Title:       postList.Title,
		Description: postList.Description,
	}

	if err := l.Lists.Insert(list); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	util.WriteJSONResponse(w, &ListGetResponse{
		TodoList: list,
	})
}

func (l *ListsAPI) delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if err := l.Lists.Delete(id); err == model.ListNotFound {
		http.Error(w, "List not found", 404)
		return
	} else if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

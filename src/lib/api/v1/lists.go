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

type ListListResponse struct {
	TodoLists []*model.TodoList `json:"todoLists"`
}

type ListGetResponse struct {
	TodoList *model.TodoList `json:"todoList"`
}

func listsList(w http.ResponseWriter, r *http.Request) {
	if lists, err := model.GetAllTodoLists(lib.GetDB()); err != nil {
		http.Error(w, err.Error(), 500)
	} else {
		util.WriteJSONResponse(w, &ListListResponse{
			TodoLists: lists,
		})
	}
}

func listGet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	if list, err := model.FindTodoList(lib.GetDB(), id); err != nil {
		http.Error(w, err.Error(), 500)
	} else if list == nil {
		http.Error(w, "List not found", 404)
	} else {
		util.WriteJSONResponse(w, &ListGetResponse{
			TodoList: list,
		})
	}
}

func listPut(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	list, err := model.FindTodoList(lib.GetDB(), id)
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

	if err := model.UpdateTodoList(lib.GetDB(), list); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func listPost(w http.ResponseWriter, r *http.Request) {
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

	if err := model.InsertTodoList(lib.GetDB(), list); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func init() {
	lib.Routes.HandleFunc("/api/v1/todoLists", listsList).Methods("GET")
	lib.Routes.HandleFunc("/api/v1/todoLists/{id}", listGet).Methods("GET")
	lib.Routes.HandleFunc("/api/v1/todoLists/{id}", listPut).Methods("PUT")
	lib.Routes.HandleFunc("/api/v1/todoLists", listPost).Methods("POST")
}

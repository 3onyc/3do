package api1

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/3onyc/threedo-backend"
	"github.com/3onyc/threedo-backend/model"
	"github.com/3onyc/threedo-backend/util"
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
	if lists, err := model.GetAllTodoLists(threedo.GetDB()); err != nil {
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

	if list, err := model.FindTodoList(threedo.GetDB(), id); err != nil {
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

	list, err := model.FindTodoList(threedo.GetDB(), id)
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

	if err := model.UpdateTodoList(threedo.GetDB(), list); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	util.WriteJSONResponse(w, &ListGetResponse{
		TodoList: list,
	})
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

	if err := model.InsertTodoList(threedo.GetDB(), list); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	util.WriteJSONResponse(w, &ListGetResponse{
		TodoList: list,
	})
}

func listDelete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if err := model.DeleteTodoList(threedo.GetDB(), id); err == model.ListNotFound {
		http.Error(w, "List not found", 404)
		return
	} else if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func init() {
	threedo.Routes.HandleFunc("/api/v1/todoLists", listsList).Methods("GET")
	threedo.Routes.HandleFunc("/api/v1/todoLists/{id}", listGet).Methods("GET")
	threedo.Routes.HandleFunc("/api/v1/todoLists/{id}", listPut).Methods("PUT")
	threedo.Routes.HandleFunc("/api/v1/todoLists/{id}", listDelete).Methods("DELETE")
	threedo.Routes.HandleFunc("/api/v1/todoLists", listPost).Methods("POST")
}

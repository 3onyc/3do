package v1

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

type ListListResponse struct {
	TodoLists []*model.TodoList `json:"todoLists"`
}

type ListGetResponse struct {
	TodoList *model.TodoList `json:"todoList"`
}

func listsList(w http.ResponseWriter, r *http.Request) {
	if lists, err := model.GetAllTodoLists(model.GetDB()); err != nil {
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

	if list, err := model.FindTodoList(model.GetDB(), id); err != nil {
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

	list, err := model.FindTodoList(model.GetDB(), id)
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

	if err := model.UpdateTodoList(model.GetDB(), list); err != nil {
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

	if err := model.InsertTodoList(model.GetDB(), list); err != nil {
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

	if err := model.DeleteTodoList(model.GetDB(), id); err == model.ListNotFound {
		http.Error(w, "List not found", 404)
		return
	} else if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func init() {
	api.Routes.HandleFunc("/api/v1/todoLists", listsList).Methods("GET")
	api.Routes.HandleFunc("/api/v1/todoLists/{id}", listGet).Methods("GET")
	api.Routes.HandleFunc("/api/v1/todoLists/{id}", listPut).Methods("PUT")
	api.Routes.HandleFunc("/api/v1/todoLists/{id}", listDelete).Methods("DELETE")
	api.Routes.HandleFunc("/api/v1/todoLists", listPost).Methods("POST")
}

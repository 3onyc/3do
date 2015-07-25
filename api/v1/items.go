// +build nobuild
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

type ItemListResponse struct {
	TodoItems []*model.TodoItem `json:"todoItems"`
}

type ItemGetResponse struct {
	TodoItem *model.TodoItem `json:"todoItem"`
}

func itemsList(w http.ResponseWriter, r *http.Request) {
	if items, err := model.GetAllTodoItems(threedo.GetDB()); err != nil {
		http.Error(w, err.Error(), 500)
	} else {
		util.WriteJSONResponse(w, &ItemListResponse{
			TodoItems: items,
		})
	}
}

func itemGet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	if item, err := model.FindTodoItem(threedo.GetDB(), id); err != nil {
		http.Error(w, err.Error(), 500)
	} else if item == nil {
		http.Error(w, "Item not found", 404)
	} else {
		util.WriteJSONResponse(w, &ItemGetResponse{
			TodoItem: item,
		})
	}
}

func itemPut(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	item, err := model.FindTodoItem(threedo.GetDB(), id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if item == nil {
		http.Error(w, "Item not found", 404)
		return
	}

	var payload = &ItemGetResponse{}
	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding JSON (%s)", err), 500)
		return
	}

	item.Title = payload.TodoItem.Title
	item.Description = payload.TodoItem.Description
	item.Done = payload.TodoItem.Done
	item.DoneAt = payload.TodoItem.DoneAt
	item.Group = payload.TodoItem.Group

	if err := model.UpdateTodoItem(threedo.GetDB(), item); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	util.WriteJSONResponse(w, &ItemGetResponse{
		TodoItem: item,
	})
}

func itemPost(w http.ResponseWriter, r *http.Request) {
	var payload = &ItemGetResponse{}
	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding JSON (%s)", err), 500)
		return
	}

	postItem := payload.TodoItem
	item := &model.TodoItem{
		Title:       postItem.Title,
		Description: postItem.Description,
		Done:        postItem.Done,
		DoneAt:      postItem.DoneAt,
		Group:       postItem.Group,
	}

	if err := model.InsertTodoItem(threedo.GetDB(), item); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	util.WriteJSONResponse(w, &ItemGetResponse{
		TodoItem: item,
	})
}

func itemDelete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if err := model.DeleteTodoItem(threedo.GetDB(), id); err == model.ItemNotFound {
		http.Error(w, err.Error(), 404)
		return
	} else if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func init() {
	threedo.Routes.HandleFunc("/api/v1/todoItems", itemsList).Methods("GET")
	threedo.Routes.HandleFunc("/api/v1/todoItems/{id}", itemGet).Methods("GET")
	threedo.Routes.HandleFunc("/api/v1/todoItems/{id}", itemPut).Methods("PUT")
	threedo.Routes.HandleFunc("/api/v1/todoItems/{id}", itemDelete).Methods("DELETE")
	threedo.Routes.HandleFunc("/api/v1/todoItems", itemPost).Methods("POST")
}

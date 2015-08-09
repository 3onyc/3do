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

type ItemsAPI struct {
	*util.Context
}

func NewItemsAPI(ctx *util.Context) *ItemsAPI {
	return &ItemsAPI{ctx}
}

func (i *ItemsAPI) Register() {
	i.Router.HandleFunc("/api/v1/todoItems", i.list).Methods("GET")
	i.Router.HandleFunc("/api/v1/todoItems/{id}", i.get).Methods("GET")
	i.Router.HandleFunc("/api/v1/todoItems/{id}", i.put).Methods("PUT")
	i.Router.HandleFunc("/api/v1/todoItems/{id}", i.delete).Methods("DELETE")
	i.Router.HandleFunc("/api/v1/todoItems/{id}/done", i.done).Methods("PUT")
	i.Router.HandleFunc("/api/v1/todoItems/{id}/todo", i.todo).Methods("PUT")
	i.Router.HandleFunc("/api/v1/todoItems", i.post).Methods("POST")
}

type ItemListResponse struct {
	TodoItems []*model.TodoItem `json:"todoItems"`
}

type ItemGetResponse struct {
	TodoItem *model.TodoItem `json:"todoItem"`
}

func (i *ItemsAPI) list(w http.ResponseWriter, r *http.Request) {
	if items, err := i.Items.GetAll(); err != nil {
		http.Error(w, err.Error(), 500)
	} else {
		util.WriteJSONResponse(w, &ItemListResponse{
			TodoItems: items,
		})
	}
}

func (i *ItemsAPI) get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if item, err := i.Items.Find(id); err != nil {
		http.Error(w, err.Error(), 500)
	} else if item == nil {
		http.Error(w, "Item not found", 404)
	} else {
		util.WriteJSONResponse(w, &ItemGetResponse{
			TodoItem: item,
		})
	}
}

func (i *ItemsAPI) put(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	item, err := i.Items.Find(id)
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

	if err := i.Items.Update(item); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	util.WriteJSONResponse(w, &ItemGetResponse{
		TodoItem: item,
	})
}

func (i *ItemsAPI) post(w http.ResponseWriter, r *http.Request) {
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
		GroupID:     postItem.GroupID,
	}

	if err := i.Items.Insert(item); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	util.WriteJSONResponse(w, &ItemGetResponse{
		TodoItem: item,
	})
}

func (i *ItemsAPI) delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if err := i.Items.Delete(id); err == model.ItemNotFound {
		http.Error(w, err.Error(), 404)
		return
	} else if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func (i *ItemsAPI) done(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if item, err := i.Items.Find(id); err != nil {
		http.Error(w, err.Error(), 500)
	} else if item == nil {
		http.Error(w, "Item not found", 404)
	} else {
		util.WriteJSONResponse(w, true)
		i.Bus.Publish("item:done", item)
	}
}

func (i *ItemsAPI) todo(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if item, err := i.Items.Find(id); err != nil {
		http.Error(w, err.Error(), 500)
	} else if item == nil {
		http.Error(w, "Item not found", 404)
	} else {
		i.Bus.Publish("item:todo", item)
		util.WriteJSONResponse(w, true)
	}
}

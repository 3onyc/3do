package api1

import (
	"github.com/gorilla/mux"
	"lib"
	"lib/model"
	"lib/util"
	"net/http"
	"strconv"
)

type ListListResponse struct {
	TodoLists []model.TodoList `json:"todoLists"`
}

type ListGetResponse struct {
	TodoList model.TodoList `json:"todoList"`
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
			TodoList: *list,
		})
	}
}

func init() {
	lib.Routes.HandleFunc("/api/v1/todoLists", listsList).Methods("GET")
	lib.Routes.HandleFunc("/api/v1/todoLists/{id}", listGet).Methods("GET")
}

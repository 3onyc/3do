package api1

import (
	"fmt"
	"github.com/gorilla/mux"
	"lib"
	"lib/model"
	"lib/util"
	"net/http"
)

type ListListResponse struct {
	TodoLists []model.TodoList `json:"todoLists"`
}

type ListGetResponse struct {
	TodoList model.TodoList `json:"todoList"`
}

func listsList(w http.ResponseWriter, r *http.Request) {
	db, err := lib.GetDB()
	if err != nil {
		http.Error(w, fmt.Sprintf("DB error (%s)", err.Error()), 500)
	}

	var lists []model.TodoList
	db.Preload("Groups").Find(&lists)

	util.WriteJSONResponse(w, &ListListResponse{
		TodoLists: lists,
	})
}

func listGet(w http.ResponseWriter, r *http.Request) {
	var list = &model.TodoList{}
	id := mux.Vars(r)["id"]

	db, err := lib.GetDB()
	if err != nil {
		http.Error(w, fmt.Sprintf("DB error (%s)", err.Error()), 500)
	}

	if db.Preload("Groups").First(list, id); list == nil {
		http.Error(w, "List Not Found", 404)
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

// +build nobuild
package api1

import (
	"github.com/gorilla/mux"
	"lib"
	"lib/model"
	"lib/util"
	"net/http"
	"strconv"
)

type ItemListResponse struct {
	TodoItems []model.TodoItem `json:"todoItems"`
}

type ItemGetResponse struct {
	TodoItem model.TodoItem `json:"todoItem"`
}

func itemsList(w http.ResponseWriter, r *http.Request) {
	if items, err := model.GetAllTodoItems(lib.GetDB()); err != nil {
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

	if item, err := model.FindTodoItem(lib.GetDB(), id); err != nil {
		http.Error(w, err.Error(), 500)
	} else if item == nil {
		http.Error(w, "Item not found", 404)
	} else {
		util.WriteJSONResponse(w, &ItemGetResponse{
			TodoItem: *item,
		})
	}
}

func itemPut(w http.ResponseWriter, r *http.Request) {
	// var idx int
	// id := mux.Vars(r)["id"]

	// for i, v := range ItemsFixture {
	// 	if strconv.FormatUint(uint64(v.ID), 10) == id {
	// 		idx = i
	// 		break
	// 	}
	// }

	// if idx == 0 {
	// 	http.Error(w, "List Not Found", 404)
	// 	return
	// }

	// var payload = &ItemGetResponse{}
	// if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
	// 	http.Error(w, fmt.Sprintf("Error encoding JSON (%s)", err), 500)
	// 	return
	// }

	// ItemsFixture[idx] = payload.TodoItem
}

func itemPost(w http.ResponseWriter, r *http.Request) {
	// var maxID uint
	// for _, v := range ItemsFixture {
	// 	if v.ID > maxID {
	// 		maxID = v.ID
	// 	}
	// }

	// var payload = &ItemGetResponse{}
	// if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
	// 	http.Error(w, fmt.Sprintf("Error encoding JSON (%s)", err), 500)
	// 	return
	// }

	// payload.TodoItem.ID = maxID + 1
	// ItemsFixture = append(ItemsFixture, payload.TodoItem)
}

func init() {
	lib.Routes.HandleFunc("/api/v1/todoItems", itemsList).Methods("GET")
	lib.Routes.HandleFunc("/api/v1/todoItems/{id}", itemGet).Methods("GET")
	// lib.Routes.HandleFunc("/api/v1/todoItems/{id}", itemPut).Methods("PUT")
	// lib.Routes.HandleFunc("/api/v1/todoItems", itemPost).Methods("POST")
}

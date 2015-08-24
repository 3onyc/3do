package v1

import (
	"encoding/json"
	"fmt"
	"github.com/3onyc/3do/model"
	"github.com/3onyc/3do/parser"
	"github.com/3onyc/3do/util"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type ListsAPI struct {
	*util.Context
}

func NewListsAPI(ctx *util.Context) *ListsAPI {
	return &ListsAPI{ctx}
}

func (l *ListsAPI) Register() {
	l.Router.HandleFunc("/api/v1/todoLists", l.list).Methods("GET")
	l.Router.HandleFunc("/api/v1/todoLists/import", l.importList).Methods("POST")
	l.Router.HandleFunc("/api/v1/todoLists/{id}/export", l.export).Methods("GET")
	l.Router.HandleFunc("/api/v1/todoLists/{id}", l.get).Methods("GET")
	l.Router.HandleFunc("/api/v1/todoLists/{id}", l.put).Methods("PUT")
	l.Router.HandleFunc("/api/v1/todoLists/{id}", l.delete).Methods("DELETE")
	l.Router.HandleFunc("/api/v1/todoLists", l.post).Methods("POST")
}

type ListListResponse struct {
	TodoLists []*model.TodoList `json:"todoLists"`
}

type ListGetResponse struct {
	TodoList   *model.TodoList    `json:"todoList"`
	TodoGroups []*model.TodoGroup `json:"todoGroups,omitempty"`
	TodoItems  []*model.TodoItem  `json:"todoItems,omitempty"`
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

func (l *ListsAPI) importList(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	importText := r.FormValue("import_text") == "true"

	var file io.Reader
	if importText {
		file = strings.NewReader(r.FormValue("text"))
	} else {
		var err error
		file, _, err = r.FormFile("file")
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}

	list, err := parser.Parse(title, file)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if err := l.Lists.Insert(list); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	for _, g := range list.Groups {
		g.ListID = list.ID.Int64
		if err := l.Groups.Insert(g); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		for _, i := range g.Items {
			i.GroupID = g.ID.Int64
			if err := l.Items.Insert(i); err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
		}
	}

	util.WriteJSONResponse(w, &ListGetResponse{
		TodoList: list,
	})
}

func (l *ListsAPI) get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	list, err := l.Lists.Find(id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	} else if list == nil {
		http.Error(w, "List not found", 404)
		return
	}

	groups, err := l.Groups.GetAllForList(id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	items := []*model.TodoItem{}
	for _, g := range groups {
		is, err := l.Items.GetAllForGroup(g.ID.Int64)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		items = append(items, is...)
	}

	util.WriteJSONResponse(w, &ListGetResponse{
		TodoList:   list,
		TodoGroups: groups,
		TodoItems:  items,
	})
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

func (l *ListsAPI) export(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if list, err := l.Lists.Find(id); err != nil {
		if err.Error() == "sql: no rows in result set" {
			http.Error(w, "File Not Found", 404)
		} else {
			http.Error(w, err.Error(), 500)
		}
	} else if list == nil {
		http.Error(w, "List not found", 404)
	} else {
		// TODO: refactor into the repositories
		if err := l.Groups.AddGroupsToList(list); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		for _, g := range list.Groups {
			if err := l.Items.AddItemsToGroup(g); err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
		}

		w.Header().Set("Content-Type", "text/markdown")
		w.Header().Set("Content-Disposition", fmt.Sprintf(
			"attachment; filename=%s.md",
			list.Title,
		))

		parser.NewWriter().Write(list, w)
	}
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

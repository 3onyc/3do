package model

import (
	"github.com/jmoiron/sqlx"
	"lib/util"
	"time"
)

type TodoList struct {
	ID          util.NullInt64 `json:"id,string"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	CreatedAt   *time.Time     `json:"createdAt" db:"created_at"`
	UpdatedAt   *time.Time     `json:"updatedAt" db:"updated_at"`
	Groups      []int64        `json:"groups,string"`
}

func InsertTodoList(db *sqlx.DB, list *TodoList) error {
	now := time.Now()
	list.CreatedAt = &now
	list.UpdatedAt = &now

	r, err := db.NamedExec(TODO_LIST_INSERT_QUERY, list)
	if err != nil {
		return err
	}

	lastID, err := r.LastInsertId()
	if err != nil {
		return err
	}

	list.ID = util.NewNullInt64(lastID)
	return nil
}

func UpdateTodoList(db *sqlx.DB, list *TodoList) error {
	now := time.Now()
	list.UpdatedAt = &now

	if _, err := db.NamedExec(TODO_LIST_UPDATE_QUERY, list); err != nil {
		return err
	}

	return nil
}

func SetListGroupIDs(db *sqlx.DB, l *TodoList) error {
	return db.Select(&l.Groups, TODO_GROUP_SELECT_IDS_WITH_LIST_QUERY, l.ID)
}

func GetAllTodoLists(db *sqlx.DB) ([]*TodoList, error) {
	var lists = []*TodoList{}
	if err := db.Select(&lists, TODO_LIST_SELECT_QUERY); err != nil {
		return nil, err
	}

	for _, list := range lists {
		if err := SetListGroupIDs(db, list); err != nil {
			return nil, err
		}
	}

	return lists, nil
}

func FindTodoList(db *sqlx.DB, id int64) (*TodoList, error) {
	var list TodoList
	if err := db.Get(&list, TODO_LIST_SELECT_ID_QUERY, id); err != nil {
		return nil, err
	}

	if err := SetListGroupIDs(db, &list); err != nil {
		return nil, err
	}

	return &list, nil
}

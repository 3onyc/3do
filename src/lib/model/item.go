package model

import (
	"github.com/jmoiron/sqlx"
	"lib/util"
	"time"
)

type TodoItem struct {
	ID          util.NullInt64 `json:"id,string"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Done        bool           `json:"done"`
	DoneAt      *time.Time     `json:"doneAt" db:"done_at"`
	CreatedAt   *time.Time     `json:"createdAt" db:"created_at"`
	UpdatedAt   *time.Time     `json:"updatedAt" db:"updated_at"`
	Group       int64          `json:"group,string" db:"group_id"`
}

func InsertTodoItem(db *sqlx.DB, item *TodoItem) error {
	now := time.Now()
	item.CreatedAt = &now
	item.UpdatedAt = &now

	r, err := db.NamedExec(TODO_ITEM_INSERT_QUERY, item)
	if err != nil {
		return err
	}

	lastID, err := r.LastInsertId()
	if err != nil {
		return err
	}

	item.ID = util.NewNullInt64(lastID)
	return nil
}

func UpdateTodoItem(db *sqlx.DB, item *TodoItem) error {
	now := time.Now()
	item.UpdatedAt = &now

	if _, err := db.NamedExec(TODO_ITEM_UPDATE_QUERY, item); err != nil {
		return err
	}

	return nil
}

func GetAllTodoItems(db *sqlx.DB) ([]*TodoItem, error) {
	var items = []*TodoItem{}
	if err := db.Select(&items, TODO_ITEM_SELECT_QUERY); err != nil {
		return nil, err
	}

	return items, nil
}

func GetAllTodoItemsForGroup(db *sqlx.DB, groupID int64) ([]*TodoItem, error) {
	var items = []*TodoItem{}
	if err := db.Select(&items, TODO_ITEM_SELECT_WITH_GROUP_QUERY, groupID); err != nil {
		return nil, err
	}

	return items, nil
}

func FindTodoItem(db *sqlx.DB, id int64) (*TodoItem, error) {
	var item TodoItem
	if err := db.Get(&item, TODO_ITEM_SELECT_ID_QUERY, id); err != nil {
		return nil, err
	}

	return &item, nil
}

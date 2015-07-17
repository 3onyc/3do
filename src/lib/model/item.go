package model

import (
	"github.com/jmoiron/sqlx"
	"lib/util"
	"time"
)

var (
	TODO_ITEM_INSERT_QUERY = `
		INSERT INTO
			"todo_items"
		(
			"id",
			"title",
			"description",
			"done",
			"done_at",
			"created_at",
			"updated_at",
			"group_id"
		) VALUES (
			:id,
			:title,
			:description,
			:done,
			:done_at,
			:created_at,
			:updated_at,
			:group_id
		)
	`
	TODO_ITEM_SELECT_QUERY = `
		SELECT
			"id",
			"title",
			"description",
			"done",
			"done_at",
			"created_at",
			"updated_at",
			"group_id"
		FROM
			"todo_items"
	`
	TODO_ITEM_SELECT_ID_QUERY = TODO_ITEM_SELECT_QUERY + `
		WHERE
			id = ?
	`
	TODO_ITEM_SELECT_WITH_GROUP_QUERY = TODO_ITEM_SELECT_QUERY + `
		WHERE
			group_id = ?
	`
	TODO_ITEM_SELECT_IDS_WITH_GROUP_QUERY = `
		SELECT 
			"id" 
		FROM 
			"todo_items"
		WHERE 
			"group_id" = ?
	`
)

type TodoItem struct {
	ID          util.NullInt64 `json:"id,string"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Done        bool           `json:"done"`
	DoneAt      time.Time      `json:"doneAt" db:"done_at"`
	CreatedAt   time.Time      `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time      `json:"updatedAt" db:"updated_at"`
	Group       int64          `json:"group,string" db:"group_id"`
}

func InsertTodoItem(db *sqlx.DB, item TodoItem) (int64, error) {
	r, err := db.NamedExec(TODO_ITEM_INSERT_QUERY, item)
	if err != nil {
		return 0, err
	}

	lastID, err := r.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int64(lastID), nil
}

func GetAllTodoItems(db *sqlx.DB) ([]TodoItem, error) {
	var items = []TodoItem{}
	if err := db.Select(&items, TODO_ITEM_SELECT_QUERY); err != nil {
		return nil, err
	}

	return items, nil
}

func GetAllTodoItemsForGroup(db *sqlx.DB, groupID int64) ([]TodoItem, error) {
	var items = []TodoItem{}
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

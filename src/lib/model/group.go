package model

import (
	"github.com/jmoiron/sqlx"
	"lib/util"
	"time"
)

var (
	TODO_GROUP_INSERT_QUERY = `
		INSERT INTO
			"todo_groups"
		(
			"id",
			"title",
			"created_at",
			"updated_at",
			"list_id"
		) VALUES (
			:id,
			:title,
			:created_at,
			:updated_at,
			:list_id
		)
	`
	TODO_GROUP_SELECT_QUERY = `
		SELECT
			"id",
			"title",
			"created_at",
			"updated_at",
			"list_id"
		FROM
			"todo_groups"
	`
	TODO_GROUP_SELECT_ID_QUERY = TODO_GROUP_SELECT_QUERY + `
		WHERE
			id = ?
	`
	TODO_GROUP_SELECT_WITH_LIST_QUERY = TODO_GROUP_SELECT_QUERY + `
		WHERE
			list_id = ?
	`
	TODO_GROUP_SELECT_IDS_WITH_LIST_QUERY = `
		SELECT 
			"id" 
		FROM 
			"todo_groups"
		WHERE 
			"list_id" = ?
	`
)

type TodoGroup struct {
	ID        util.NullInt64 `json:"id,string"`
	Title     string         `json:"title"`
	CreatedAt time.Time      `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time      `json:"updatedAt" db:"updated_at"`
	List      int64          `json:"list,string" db:"list_id"`
	Items     []int64        `json:"items,string"`
}

func InsertTodoGroup(db *sqlx.DB, group TodoGroup) (int64, error) {
	r, err := db.NamedExec(TODO_GROUP_INSERT_QUERY, group)
	if err != nil {
		return 0, err
	}

	lastID, err := r.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int64(lastID), nil
}

func SetGroupItemIDs(db *sqlx.DB, g *TodoGroup) error {
	return db.Select(&g.Items, TODO_ITEM_SELECT_IDS_WITH_GROUP_QUERY, g.ID)
}

func GetAllTodoGroups(db *sqlx.DB) ([]TodoGroup, error) {
	var groups = []TodoGroup{}
	if err := db.Select(&groups, TODO_GROUP_SELECT_QUERY); err != nil {
		return nil, err
	}

	for _, group := range groups {
		if err := SetGroupItemIDs(db, &group); err != nil {
			return nil, err
		}
	}

	return groups, nil
}

func GetAllTodoGroupsForList(db *sqlx.DB, listID int64) ([]TodoGroup, error) {
	var groups = []TodoGroup{}
	if err := db.Select(&groups, TODO_GROUP_SELECT_WITH_LIST_QUERY, listID); err != nil {
		return nil, err
	}

	for _, group := range groups {
		if err := SetGroupItemIDs(db, &group); err != nil {
			return nil, err
		}
	}

	return groups, nil
}

func FindTodoGroup(db *sqlx.DB, id int64) (*TodoGroup, error) {
	var group TodoGroup
	if err := db.Get(&group, TODO_GROUP_SELECT_ID_QUERY, id); err != nil {
		return nil, err
	}

	if err := SetGroupItemIDs(db, &group); err != nil {
		return nil, err
	}

	return &group, nil
}

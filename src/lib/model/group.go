package model

import (
	"github.com/jmoiron/sqlx"
	"lib/util"
	"time"
)

type TodoGroup struct {
	ID        util.NullInt64 `json:"id,string"`
	Title     string         `json:"title"`
	CreatedAt *time.Time     `json:"createdAt" db:"created_at"`
	UpdatedAt *time.Time     `json:"updatedAt" db:"updated_at"`
	List      int64          `json:"list,string" db:"list_id"`
	Items     []int64        `json:"items,string"`
}

func InsertTodoGroup(db *sqlx.DB, group *TodoGroup) (int64, error) {
	now := time.Now()
	group.CreatedAt = &now
	group.UpdatedAt = &now

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

func UpdateTodoGroup(db *sqlx.DB, group *TodoGroup) error {
	now := time.Now()
	group.UpdatedAt = &now

	if _, err := db.NamedExec(TODO_GROUP_UPDATE_QUERY, group); err != nil {
		return err
	}

	return nil
}

func SetGroupItemIDs(db *sqlx.DB, g *TodoGroup) error {
	return db.Select(&g.Items, TODO_ITEM_SELECT_IDS_WITH_GROUP_QUERY, g.ID)
}

func GetAllTodoGroups(db *sqlx.DB) ([]*TodoGroup, error) {
	var groups = []*TodoGroup{}
	if err := db.Select(&groups, TODO_GROUP_SELECT_QUERY); err != nil {
		return nil, err
	}

	for _, group := range groups {
		if err := SetGroupItemIDs(db, group); err != nil {
			return nil, err
		}
	}

	return groups, nil
}

func GetAllTodoGroupsForList(db *sqlx.DB, listID int64) ([]*TodoGroup, error) {
	var groups = []*TodoGroup{}
	if err := db.Select(&groups, TODO_GROUP_SELECT_WITH_LIST_QUERY, listID); err != nil {
		return nil, err
	}

	for _, group := range groups {
		if err := SetGroupItemIDs(db, group); err != nil {
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

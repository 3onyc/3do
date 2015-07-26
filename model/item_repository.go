package model

import (
	"errors"
	"github.com/3onyc/3do/dbtype"
	"github.com/jmoiron/sqlx"
	"time"
)

var (
	ItemNotFound = errors.New("Item Not Found")
)

type TodoItemRepository struct {
	DB *sqlx.DB
}

func NewTodoItemRepository(db *sqlx.DB) *TodoItemRepository {
	return &TodoItemRepository{
		DB: db,
	}
}

func (repo *TodoItemRepository) Insert(item *TodoItem) error {
	now := time.Now()
	item.CreatedAt = &now
	item.UpdatedAt = &now

	r, err := repo.DB.NamedExec(TODO_ITEM_INSERT_QUERY, item)
	if err != nil {
		return err
	}

	lastID, err := r.LastInsertId()
	if err != nil {
		return err
	}

	item.ID = dbtype.NewNullInt64(lastID)
	return nil
}

func (repo *TodoItemRepository) Update(item *TodoItem) error {
	now := time.Now()
	item.UpdatedAt = &now

	if _, err := repo.DB.NamedExec(TODO_ITEM_UPDATE_QUERY, item); err != nil {
		return err
	}

	return nil
}

func (repo *TodoItemRepository) Delete(id int64) error {
	if result, err := repo.DB.Exec(TODO_ITEM_DELETE_QUERY, id); err != nil {
		return err
	} else if affected, err := result.RowsAffected(); err != nil {
		return err
	} else if affected == 0 {
		return ItemNotFound
	} else {
		return nil
	}
}

func (repo *TodoItemRepository) GetAll() ([]*TodoItem, error) {
	var items = []*TodoItem{}
	if err := repo.DB.Select(&items, TODO_ITEM_SELECT_QUERY); err != nil {
		return nil, err
	}

	return items, nil
}

func (repo *TodoItemRepository) GetAllForGroup(groupID int64) ([]*TodoItem, error) {
	var items = []*TodoItem{}
	if err := repo.DB.Select(&items, TODO_ITEM_SELECT_WITH_GROUP_QUERY, groupID); err != nil {
		return nil, err
	}

	return items, nil
}

func (repo *TodoItemRepository) Find(id int64) (*TodoItem, error) {
	var item TodoItem
	if err := repo.DB.Get(&item, TODO_ITEM_SELECT_ID_QUERY, id); err != nil {
		return nil, err
	}

	return &item, nil
}

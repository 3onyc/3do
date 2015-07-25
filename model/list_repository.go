package model

import (
	"errors"
	"github.com/3onyc/threedo-backend/dbtype"
	"github.com/jmoiron/sqlx"
	"time"
)

var (
	ListNotFound = errors.New("List Not Found")
)

type TodoListRepository struct {
	DB *sqlx.DB
}

func NewTodoListRepository(db *sqlx.DB) *TodoListRepository {
	return &TodoListRepository{
		DB: db,
	}
}

func (repo *TodoListRepository) Insert(list *TodoList) error {
	now := time.Now()
	list.CreatedAt = &now
	list.UpdatedAt = &now

	r, err := repo.DB.NamedExec(TODO_LIST_INSERT_QUERY, list)
	if err != nil {
		return err
	}

	lastID, err := r.LastInsertId()
	if err != nil {
		return err
	}

	list.ID = dbtype.NewNullInt64(lastID)
	return nil
}

func (repo *TodoListRepository) Update(list *TodoList) error {
	now := time.Now()
	list.UpdatedAt = &now

	if _, err := repo.DB.NamedExec(TODO_LIST_UPDATE_QUERY, list); err != nil {
		return err
	}

	return nil
}

func (repo *TodoListRepository) Delete(id int64) error {
	if result, err := repo.DB.Exec(TODO_LIST_DELETE_QUERY, id); err != nil {
		return err
	} else if affected, err := result.RowsAffected(); err != nil {
		return err
	} else if affected == 0 {
		return ListNotFound
	} else {
		return nil
	}
}

func (repo *TodoListRepository) SetGroupIDs(l *TodoList) error {
	return repo.DB.Select(&l.Groups, TODO_GROUP_SELECT_IDS_WITH_LIST_QUERY, l.ID)
}

func (repo *TodoListRepository) GetAll() ([]*TodoList, error) {
	var lists = []*TodoList{}
	if err := repo.DB.Select(&lists, TODO_LIST_SELECT_QUERY); err != nil {
		return nil, err
	}

	for _, list := range lists {
		if err := repo.SetGroupIDs(list); err != nil {
			return nil, err
		}
	}

	return lists, nil
}

func (repo *TodoListRepository) Find(id int64) (*TodoList, error) {
	var list TodoList
	if err := repo.DB.Get(&list, TODO_LIST_SELECT_ID_QUERY, id); err != nil {
		return nil, err
	}

	if err := repo.SetGroupIDs(&list); err != nil {
		return nil, err
	}

	return &list, nil
}

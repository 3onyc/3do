package model

import (
	"errors"
	"github.com/3onyc/3do/dbtype"
	"github.com/jmoiron/sqlx"
	"time"
)

var (
	GroupNotFound = errors.New("List Not Found")
)

type TodoGroupRepository struct {
	DB *sqlx.DB
}

func NewTodoGroupRepository(db *sqlx.DB) *TodoGroupRepository {
	return &TodoGroupRepository{
		DB: db,
	}
}

func (repo *TodoGroupRepository) Insert(group *TodoGroup) error {
	now := time.Now()
	group.CreatedAt = &now
	group.UpdatedAt = &now

	r, err := repo.DB.NamedExec(TODO_GROUP_INSERT_QUERY, group)
	if err != nil {
		return err
	}

	lastID, err := r.LastInsertId()
	if err != nil {
		return err
	}

	group.ID = dbtype.NewNullInt64(lastID)
	return nil
}

func (repo *TodoGroupRepository) Update(group *TodoGroup) error {
	now := time.Now()
	group.UpdatedAt = &now

	if _, err := repo.DB.NamedExec(TODO_GROUP_UPDATE_QUERY, group); err != nil {
		return err
	}

	return nil
}

func (repo *TodoGroupRepository) Delete(id int64) error {
	if result, err := repo.DB.Exec(TODO_GROUP_DELETE_QUERY, id); err != nil {
		return err
	} else if affected, err := result.RowsAffected(); err != nil {
		return err
	} else if affected == 0 {
		return GroupNotFound
	} else {
		return nil
	}
}

func (repo *TodoGroupRepository) SetItemIDs(g *TodoGroup) error {
	return repo.DB.Select(&g.ItemIDs, TODO_ITEM_SELECT_IDS_WITH_GROUP_QUERY, g.ID)
}

func (repo *TodoGroupRepository) GetAll() ([]*TodoGroup, error) {
	var groups = []*TodoGroup{}
	if err := repo.DB.Select(&groups, TODO_GROUP_SELECT_QUERY); err != nil {
		return nil, err
	}

	for _, group := range groups {
		if err := repo.SetItemIDs(group); err != nil {
			return nil, err
		}
	}

	return groups, nil
}

func (repo *TodoGroupRepository) GetAllForList(listID int64) ([]*TodoGroup, error) {
	var groups = []*TodoGroup{}
	if err := repo.DB.Select(&groups, TODO_GROUP_SELECT_WITH_LIST_QUERY, listID); err != nil {
		return nil, err
	}

	for _, group := range groups {
		if err := repo.SetItemIDs(group); err != nil {
			return nil, err
		}
	}

	return groups, nil
}

func (repo *TodoGroupRepository) Find(id int64) (*TodoGroup, error) {
	var group TodoGroup
	if err := repo.DB.Get(&group, TODO_GROUP_SELECT_ID_QUERY, id); err != nil {
		return nil, err
	}

	if err := repo.SetItemIDs(&group); err != nil {
		return nil, err
	}

	return &group, nil
}

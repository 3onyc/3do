package lib

import (
	"github.com/jinzhu/gorm"
	"lib/model"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var (
	DB *gorm.DB
)

func GetDB() (*gorm.DB, error) {
	if DB != nil {
		return DB, nil
	}

	db, err := gorm.Open("sqlite3", "/tmp/3do.sqlite3")
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(
		&model.TodoItem{},
		&model.TodoGroup{},
		&model.TodoList{},
	)

	DB = &db
	return &db, nil
}

func SeedDB(db *gorm.DB) error {
	i1 := model.TodoItem{
		ID:          1,
		Title:       "Item 1",
		Description: "# Foo",
		Done:        false,
		DoneAt:      time.Now(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		GroupID:     1,
	}
	i2 := model.TodoItem{
		ID:          2,
		Title:       "Item 2",
		Description: "Bar",
		Done:        false,
		DoneAt:      time.Now(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		GroupID:     1,
	}

	g1 := model.TodoGroup{
		ID:        1,
		Title:     "Group 1",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Items:     []model.TodoItem{i1, i2},
		ListID:    1,
	}

	l1 := model.TodoList{
		ID:          1,
		Title:       "Hello",
		Description: "Foo",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Groups:      []model.TodoGroup{g1},
	}
	l2 := model.TodoList{
		ID:          2,
		Title:       "Bye",
		Description: "Bar",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	db.Create(&l1)
	db.Create(&l2)
	db.Create(&g1)
	db.Create(&i1)
	db.Create(&i2)

	return nil
}

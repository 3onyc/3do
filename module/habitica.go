package module

import (
	"errors"
	"github.com/3onyc/3do/model"
	"github.com/3onyc/3do/util"
	"github.com/3onyc/go-habitica"
	"github.com/namsral/flag"
)

var (
	HABIT_USER = flag.String("habit-user", "", "HabitRPG user id")
	HABIT_KEY  = flag.String("habit-key", "", "HabitRPG API key")
	HABIT_TASK = flag.String("habit-task", "", "HabitRPG task ID")
)

type Habitica struct {
	*util.Context
	client *habitica.HabiticaClient
}

func NewHabitica(ctx *util.Context) *Habitica {
	return &Habitica{ctx, nil}
}

func (m *Habitica) Init() error {
	m.client = habitica.New(*HABIT_USER, *HABIT_KEY)
	up, err := m.client.Status()
	if err != nil {
		return err
	}

	if !up {
		return errors.New("Habitica is down")
	}

	m.Bus.SubscribeAsync("item:done", m.handleDone, false)
	m.Bus.SubscribeAsync("item:todo", m.handleTodo, false)

	return nil
}

func (m *Habitica) handleDone(item *model.TodoItem) {
	if err := m.client.TaskUp(*HABIT_TASK); err != nil {
		m.Log(
			"module", "habitica",
			"action", "task-up",
			"task-id", *HABIT_TASK,
			"result", false,
			"message", err,
		)
	}
}

func (m *Habitica) handleTodo(item *model.TodoItem) {
	if err := m.client.TaskDown(*HABIT_TASK); err != nil {
		m.Log(
			"module", "habitica",
			"action", "task-down",
			"task-id", *HABIT_TASK,
			"result", false,
			"message", err,
		)
	}
}

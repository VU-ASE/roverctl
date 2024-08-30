package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// This is an action that can be performed by a page
type Action[T interface{}] struct {
	Name     string
	Result   bool
	Error    error
	Started  bool
	Finished bool
	Attempt  uint // allows you to ignore results from previous attempts
	Data     *T
}

// This is a message returned after an action is performed
// it describes the action and the data that was returned
type ActionResult[T interface{}] struct {
	Name       string
	Result     bool
	Error      error
	ForAttempt uint // allows you to ignore results from previous attempts
	Data       *T
}

// A collection of Actions that can be used by a model
type Actions []*Action[any]

// Check if a given result is for a specific action and attempt
func (a ActionResult[T]) IsFor(action *Action[T]) bool {
	return a.Name == action.Name && a.ForAttempt-1 == action.Attempt
}

// ProcessResult takes an ActionResult and updates the Actions where the name matches
func (a Actions) ProcessResult(res ActionResult[any]) {
	for _, action := range a {
		if res.IsFor(action) {
			action.Result = res.Result
			action.Error = res.Error
			action.Finished = true
			action.Data = res.Data
		}
	}
}

func (a Action[T]) IsLoading() bool {
	return a.Started && !a.Finished
}

func (a Action[T]) IsSuccess() bool {
	return a.Started && a.Finished && a.Result
}

func (a Action[T]) IsError() bool {
	return a.Started && a.Finished && !a.Result
}

func (a Action[T]) IsDone() bool {
	return a.Started && a.Finished
}

func (a *Action[T]) Reset() {
	a.Started = false
	a.Finished = false
	a.Result = false
	a.Error = nil
	a.Data = nil
}

func (a *Action[T]) Start() {
	a.Started = true
}

// Generate a new ActionResult from an Action
func NewResult[T interface{}](a Action[T], success bool, err error, data *T, attempt uint) ActionResult[T] {
	return ActionResult[T]{
		Name:       a.Name,
		Result:     success,
		Error:      err,
		ForAttempt: attempt,
		Data:       data,
	}
}

func NewAction[T interface{}](name string) Action[T] {
	return Action[T]{
		Name:     name,
		Result:   false,
		Error:    nil,
		Started:  false,
		Finished: false,
		Attempt:  0,
		Data:     nil,
	}
}

type ActionFunction[T interface{}] func() (*T, error)

// Wrapper that makes your life easier when performing an action
// You can also use the oldschool method of creating a function that returns a tea.Cmd, and use tui.NewResult() with that
func PerformAction[T interface{}](action *Action[T], f ActionFunction[T]) tea.Cmd {
	action.Attempt++
	action.Started = true
	action.Finished = false
	attempt := action.Attempt

	return func() tea.Msg {
		data, err := f()
		return NewResult(*action, err == nil, err, data, attempt)
	}
}

package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// This is an action that can be performed by a page
type Action struct {
	Name     string
	Result   bool
	Error    error
	Started  bool
	Finished bool
	Attempt  uint // allows you to ignore results from previous attempts
}

// This is a message returned after an action is performed
// it describes the action and the data that was returned
type ActionResult struct {
	Name       string
	Result     bool
	Error      error
	ForAttempt uint // allows you to ignore results from previous attempts
}

// A collection of Actions that can be used by a model
type Actions []*Action

// ProcessResult takes an ActionResult and updates the Actions where the name matches
func (a Actions) ProcessResult(res ActionResult) {
	for _, action := range a {
		if action.Name == res.Name && action.Attempt == res.ForAttempt {
			action.Result = res.Result
			action.Error = res.Error
			action.Finished = true
		}
	}
}

func (a Action) IsLoading() bool {
	return a.Started && !a.Finished
}

func (a Action) IsSuccess() bool {
	return a.Started && a.Finished && a.Result
}

func (a Action) IsError() bool {
	return a.Started && a.Finished && !a.Result
}

func (a Action) IsDone() bool {
	return a.Started && a.Finished
}

// Generate a new ActionResult from an Action
func NewResult(a Action, success bool, err error, attempt uint) ActionResult {
	return ActionResult{
		Name:       a.Name,
		Result:     success,
		Error:      err,
		ForAttempt: attempt,
	}
}

func NewAction(name string) Action {
	return Action{
		Name:     name,
		Result:   false,
		Error:    nil,
		Started:  false,
		Finished: false,
		Attempt:  0,
	}
}

type ActionFunction func() error

// Wrapper that makes your life easier when performing an action
// You can also use the oldschool method of creating a function that returns a tea.Cmd, and use tui.NewResult() with that
func PerformAction(action *Action, f ActionFunction) tea.Cmd {
	action.Attempt++
	action.Started = true
	action.Finished = false
	attempt := action.Attempt

	return func() tea.Msg {
		err := f()
		action.Finished = true

		return NewResult(*action, err == nil, err, attempt)
	}
}

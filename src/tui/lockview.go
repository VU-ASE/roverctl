package tui

import (
	"fmt"

	roverlock "github.com/VU-ASE/rover/src/lock"
	"github.com/VU-ASE/rover/src/state"
	"github.com/VU-ASE/rover/src/style"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// A reusable view that can show lock information
// if not locked yet, it shows a loading state or a message, if successfully locked, it shows the string passed in
// NB: the spinner needs to be initialized already
func LockView(lockAction Action, unlockAction Action, spinner spinner.Model, view string) string {
	if lockAction.IsLoading() {
		return spinner.View() + " locking your Rover"
	} else if unlockAction.IsLoading() {
		return spinner.View() + " unlocking your Rover"
	} else if lockAction.IsError() {
		s := lipgloss.NewStyle().Foreground(style.ErrorPrimary).Render("Failed to lock your Rover")
		s += "\n\nThe operation you were about to perform required exclusive access to your Rover\nbut this could not be obtained."
		if lockAction.Error != nil {
			s += " The Rover reported the following error:"
			s += "\n > " + lipgloss.NewStyle().Foreground(style.WarningPrimary).Render(lockAction.Error.Error())
		}

		s += "\n\nRebooting the Rover by holding the power button for 10 seconds will release the lock. "
		s += lipgloss.NewStyle().Foreground(style.GrayPrimary).Render("\n\nPress 'r' to retry, or 'q' to go back.")
		return s
	} else if unlockAction.IsError() {
		s := lipgloss.NewStyle().Foreground(style.ErrorPrimary).Render("Failed to unlock your Rover")
		s += "\n\nTo grant others permission to use your Rover, it needs to be unlocked by you\nbut the unlocking failed."
		if unlockAction.Error != nil {
			s += " The Rover reported the following error:"
			s += "\n > " + lipgloss.NewStyle().Foreground(style.WarningPrimary).Render(unlockAction.Error.Error())
		}

		s += "\n\nRebooting the Rover by holding the power button for 10 seconds will release the lock. "
		s += lipgloss.NewStyle().Foreground(style.GrayPrimary).Render("\n\nPress 'r' to retry, or 'q' to go back.")
		return s
	} else {
		return view
	}
}

// A reusable update function that can be used to retry lock and unlock actions
func LockUpdate(m tea.Model, msg tea.Msg, lockAction *Action, unlockAction *Action) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "r":
			// Retry the unlock
			if !lockAction.IsSuccess() {
				return m, unlock(unlockAction)
			} else {
				// Retry the lock
				return m, lock(lockAction)
			}
		}
	case ActionResult:
		actions := Actions{lockAction, unlockAction}
		actions.ProcessResult(msg)
	}

	return m, nil
}

func LockInit(lockAction *Action, unlockAction *Action) tea.Cmd {
	return tea.Batch(
		lock(lockAction),
		unlock(unlockAction),
	)
}

func lock(lockAction *Action) tea.Cmd {
	return PerformAction(lockAction, func() error {
		// Get the active rover
		active := state.Get().RoverConnections.GetActive()
		if active == nil {
			return fmt.Errorf("Not connected to an active Rover")
		}

		// Unlock the rover
		return roverlock.Unlock(*active)
	})
}

func unlock(unlockAction *Action) tea.Cmd {
	return PerformAction(unlockAction, func() error {
		// Get the active rover
		active := state.Get().RoverConnections.GetActive()
		if active == nil {
			return fmt.Errorf("Not connected to an active Rover")
		}

		// Unlock the rover
		return roverlock.Lock(*active)
	})
}

package main

import (
	"fmt"
	"winser-cli/winser"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/samber/lo"
)

type menuFind struct {
	focusIndex  int
	searchName  InputText
	name        string
	displayName string
	description string
	executable  string
	argument    string
	startMode   winser.StartMode
	status      winser.Status
	statusBar   StatusBar
}

func newMenuFind() menuFind {
	m := menuFind{}

	tname := textinput.New()
	tname.Placeholder = "Name"
	tname.Focus()
	tname.PromptStyle = focusedStyle
	tname.TextStyle = focusedStyle
	tname.CharLimit = 256
	m.searchName = InputText{"name", tname}

	m.statusBar = StatusBar{Type: StatusInfo, Msg: "Ready"}
	return m
}

func (m menuFind) view() string {
	// title
	v := appName + "\n\n"
	v += info(" -- FIND SERVICE -- ") + "\n\n"

	// input
	v += m.searchName.Render() + "\n\n"

	// button manage
	v += lo.Ternary(m.focusIndex == 1, bginfo("> FIND "), bgsecondary("  FIND ")) + "    "
	v += lo.Ternary(m.focusIndex == 2, bgwarning("> EDIT "), bgsecondary("  EDIT ")) + "    "
	v += lo.Ternary(m.focusIndex == 3, bgdanger("> REMOVE "), bgsecondary("  REMOVE ")) + "\n\n"

	// textview
	v += secondary("Name         : ") + m.name + "\n"
	v += secondary("Display Name : ") + m.displayName + "\n"
	v += secondary("Description  : ") + m.description + "\n"
	v += secondary("Executable   : ") + m.executable + "\n"
	v += secondary("Argument     : ") + m.argument + "\n"

	startMode := warning("UNKNOWN")
	switch m.startMode {
	case winser.Auto:
		startMode = info("AUTO")
	case winser.Manual:
		startMode = info("MANUAL")
	case winser.Disabled:
		startMode = info("DISABLED")
	}
	v += secondary("Start Mode   : ") + startMode + "\n"

	status := warning("UNKNOWN")
	switch m.status {
	case winser.Running:
		status = success("RUNNING")
	case winser.Paused:
		status = warning("PAUSED")
	case winser.Stopped:
		status = danger("STOPPED")
	case winser.StartPending:
		status = warning("START PENDING")
	case winser.PausePending:
		status = warning("PAUSE PENDING")
	case winser.StopPending:
		status = warning("STOP PENDING")
	case winser.ContinuePending:
		status = warning("CONTINUE PENDING")
	}
	v += secondary("Status       : ") + status + "\n\n"

	// button control
	v += lo.Ternary(m.focusIndex == 4, bgsuccess("> START "), bgsecondary("  START ")) + "    "
	v += lo.Ternary(m.focusIndex == 5, bgdanger("> STOP "), bgsecondary("  STOP ")) + "    "
	v += lo.Ternary(m.focusIndex == 6, bgwarning("> PAUSE "), bgsecondary("  PAUSE ")) + "\n\n\n"

	// statusbar
	v += m.statusBar.Render() + "\n"

	// help
	switch m.focusIndex {
	case 1:
		v += secondary("Press ENTER to Find • ")
	case 2:
		v += secondary("Press ENTER to Edit • ")
	case 3:
		v += warning("Press ENTER to REMOVE") + secondary(" • ")
	case 4:
		v += secondary("Press ENTER to Start • ")
	case 5:
		v += secondary("Press ENTER to Stop • ")
	case 6:
		v += secondary("Press ENTER to Pause • ")
	}
	v += secondary("↑/j: up • ↓/k: down • esc: back")
	return v
}

func (m *menuFind) event(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// button action
			if s == "enter" && m.focusIndex == 1 {
				name := m.searchName.Value()
				m.statusBar.Info(fmt.Sprintf("Find (%s)", name))
				return findService(name)
			}
			if s == "enter" && m.focusIndex == 2 {
				m.statusBar.Warning("FEATURE IS UNDER CONSTRUCTION!")
				return nil
			}
			if s == "enter" && m.focusIndex == 3 {
				if m.name == "" {
					m.statusBar.Warning("Find the service first")
					return nil
				}
				m.statusBar.Warning(fmt.Sprintf("Remove (%s)", m.name))
				return removeService(m.name)
			}
			if s == "enter" && m.focusIndex == 4 {
				if m.name == "" {
					m.statusBar.Warning("Find the service first")
					return nil
				}
				m.statusBar.Info(fmt.Sprintf("Start (%s)", m.name))
				return startCtrlService(m.name)
			}
			if s == "enter" && m.focusIndex == 5 {
				if m.name == "" {
					m.statusBar.Warning("Find the service first")
					return nil
				}
				m.statusBar.Info(fmt.Sprintf("Stop (%s)", m.name))
				return stopCtrlService(m.name)
			}
			if s == "enter" && m.focusIndex == 6 {
				if m.name == "" {
					m.statusBar.Warning("Find the service first")
					return nil
				}
				m.statusBar.Info(fmt.Sprintf("Pause (%s)", m.name))
				return pauseCtrlService(m.name)
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > 6 {
				m.focusIndex = 6
			} else if m.focusIndex < 0 {
				m.focusIndex = 0
			}

			cmd := m.searchName.Focus(m.focusIndex == 0)
			return cmd
		}

	// service message
	case installServiceMsg:
		if msg.err != nil {
			m.statusBar.Error(msg.err.Error())
		} else {
			m.statusBar.Success(fmt.Sprintf("Installed (%s)", msg.name))
		}
	case findServiceMsg:
		if msg.err != nil {
			m.statusBar.Error(msg.err.Error())
			m.name = ""
			m.displayName = ""
			m.description = ""
			m.executable = ""
			m.argument = ""
			m.startMode = winser.StartModeUnknown
			m.status = winser.StatusUnknown
		} else {
			m.name = msg.service.Name
			m.displayName = msg.service.DisplayName
			m.description = msg.service.Description
			m.executable = msg.service.Executable
			m.argument = msg.service.Argument
			m.startMode = msg.service.StartMode
			m.status = msg.service.Status
			m.statusBar.Success(fmt.Sprintf("Found (%s)", m.name))
		}
	case removeServiceMsg:
		if msg.err != nil {
			m.statusBar.Error(msg.err.Error())
		} else {
			m.statusBar.Success(fmt.Sprintf("Removed (%s)", msg.name))
		}
	case startCtrlServiceMsg:
		if msg.err != nil {
			m.statusBar.Error(msg.err.Error())
		} else {
			m.statusBar.Success(fmt.Sprintf("Started (%s)", msg.name))
		}
	case stopCtrlServiceMsg:
		if msg.err != nil {
			m.statusBar.Error(msg.err.Error())
		} else {
			m.statusBar.Success(fmt.Sprintf("Stopped (%s)", msg.name))
		}
	case pauseCtrlServiceMsg:
		if msg.err != nil {
			m.statusBar.Error(msg.err.Error())
		} else {
			m.statusBar.Success(fmt.Sprintf("Paused (%s)", msg.name))
		}
	}

	// input name
	var cmd tea.Cmd
	m.searchName.Model, cmd = m.searchName.Update(msg)
	return cmd
}

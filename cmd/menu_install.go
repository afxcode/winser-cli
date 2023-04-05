package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/samber/lo"
)

type menuInstall struct {
	focusIndex int
	inputs     []Input
	statusBar  StatusBar
}

func newMenuInstall() menuInstall {
	m := menuInstall{
		inputs: make([]Input, 0),
	}

	tname := textinput.New()
	tname.Placeholder = "Name"
	tname.Focus()
	tname.PromptStyle = focusedStyle
	tname.TextStyle = focusedStyle
	tname.CharLimit = 256
	m.inputs = append(m.inputs, &InputText{"name", tname})

	tdname := textinput.New()
	tdname.Placeholder = "Display Name"
	tdname.CharLimit = 256
	m.inputs = append(m.inputs, &InputText{"dname", tdname})

	tdesc := textinput.New()
	tdesc.Placeholder = "Description"
	m.inputs = append(m.inputs, &InputText{"desc", tdesc})

	texe := textinput.New()
	texe.Placeholder = "Executable"
	m.inputs = append(m.inputs, &InputText{"exe", texe})

	targ := textinput.New()
	targ.Placeholder = "Argument"
	m.inputs = append(m.inputs, &InputText{"arg", targ})

	m.inputs = append(m.inputs, &InputSelect{
		id:      "startmode",
		label:   "Start Mode",
		options: []option{{"Auto", "auto"}, {"Manual", "manual"}, {"Disabled", "disabled"}},
	})

	m.statusBar = StatusBar{Type: StatusInfo, Msg: "Ready"}
	return m
}

func (m menuInstall) view() string {
	// title
	v := appName + "\n\n"
	v += info(" -- INSTALL SERVICE -- ") + "\n\n"

	// inputs
	selectsInputIndex := make([]int, 0)
	for i, input := range m.inputs {
		switch input := input.(type) {
		case *InputText:
			v += input.Render()
		case *InputSelect:
			v += input.Render()
			selectsInputIndex = append(selectsInputIndex, i)
		}
		v += "\n"
	}
	v += "\n"

	// button
	buttonIndex := len(m.inputs)
	v += lo.Ternary(m.focusIndex == buttonIndex, bgsuccess("> INSTALL "), bgsecondary("  INSTALL ")) + "\n\n\n"

	// statusbar
	v += m.statusBar.Render() + "\n"

	// help
	if lo.Contains(selectsInputIndex, m.focusIndex) {
		v += secondary("←/h: left • →/l: right • ")
	}
	if m.focusIndex == buttonIndex {
		v += secondary("Press ENTER to Install • ")
	}
	v += secondary("esc: back")
	return v
}

func (m *menuInstall) event(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			actionButtonIndex := len(m.inputs)
			if s == "enter" && m.focusIndex == actionButtonIndex {
				var name, displayName, description, executable, argument, startMode string
				for _, input := range m.inputs {
					switch inp := input.(type) {
					case *InputText:
						switch inp.id {
						case "name":
							name = inp.Value()
						case "dname":
							displayName = inp.Value()
						case "desc":
							description = inp.Value()
						case "exe":
							executable = inp.Value()
						case "arg":
							argument = inp.Value()
						}
					case *InputSelect:
						if inp.id == "startmode" {
							startMode = inp.selectedOption.value
						}
					}
				}

				m.statusBar.Info(fmt.Sprintf("Installing (%s)", name))
				return installService(name, displayName, description, executable, argument, startMode)
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > actionButtonIndex {
				m.focusIndex = actionButtonIndex
			} else if m.focusIndex < 0 {
				m.focusIndex = 0
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i, input := range m.inputs {
				cmds[i] = input.Focus(i == m.focusIndex)
			}

			return tea.Batch(cmds...)
		}

		// input select
		if m.focusIndex >= len(m.inputs) {
			return nil
		}
		switch input := m.inputs[m.focusIndex].(type) {
		case *InputSelect:
			switch msg.String() {
			case "right", "l":
				input.NextCursor()
			case "left", "h":
				input.PrevCursor()
			}
			return nil
		}

	// service message
	case installServiceMsg:
		if msg.err != nil {
			m.statusBar.Error(msg.err.Error())
		} else {
			m.statusBar.Success(fmt.Sprintf("Installed (%s)", msg.name))
		}
		return nil
	}

	// input text
	if m.focusIndex >= len(m.inputs) {
		return nil
	}
	switch input := m.inputs[m.focusIndex].(type) {
	case *InputText:
		var cmd tea.Cmd
		input.Model, cmd = input.Update(msg)
		return cmd
	}
	return nil
}

func (m *menuInstall) reset() {
	for i, input := range m.inputs {
		input.Focus(i == 0)
		input.SetValue("")
	}
	m.statusBar.Msg = "Ready"
	m.statusBar.Type = StatusInfo
	m.focusIndex = 0
}

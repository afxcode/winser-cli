package main

import tea "github.com/charmbracelet/bubbletea"

type menuUpdate struct {
	statusBar StatusBar
}

func newMenuUpdate() menuUpdate {
	m := menuUpdate{}
	m.statusBar = StatusBar{Type: StatusInfo, Msg: "Ready"}
	return m
}

func (m menuUpdate) view() string {
	// title
	v := appName + "\n\n"
	v += bginfo(" -- UPDATE -- ") + "\n\n"

	// statusbar
	v += m.statusBar.Render() + "\n"

	// help
	v += secondary("esc: back")
	return v
}

func (m *menuUpdate) event(msg tea.Msg) tea.Cmd {
	return nil
}

package main

import tea "github.com/charmbracelet/bubbletea"

type menuList struct {
	statusBar StatusBar
}

func newMenuList() menuList {
	m := menuList{}
	m.statusBar = StatusBar{Type: StatusWarning, Msg: "FEATURE IS UNDER CONSTRUCTION!"}
	return m
}

func (m menuList) view() string {
	// title
	v := appName + "\n\n"
	v += bginfo(" -- LIST -- ") + "\n\n"

	// statusbar
	v += m.statusBar.Render() + "\n"

	// help
	v += secondary("esc: back")
	return v
}

func (m *menuList) event(msg tea.Msg) tea.Cmd {
	return nil
}

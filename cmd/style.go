package main

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	primary   = makeFgStyle("207")
	secondary = makeFgStyle("244")
	success   = makeFgStyle("82")
	info      = makeFgStyle("123")
	warning   = makeFgStyle("214")
	danger    = makeFgStyle("196")

	bgprimary   = makeFgBgStyle("236", "207")
	bgsecondary = makeFgBgStyle("236", "244")
	bgsuccess   = makeFgBgStyle("236", "82")
	bginfo      = makeFgBgStyle("236", "123")
	bgwarning   = makeFgBgStyle("236", "214")
	bgdanger    = makeFgBgStyle("252", "196")

	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("207"))
	noStyle      = lipgloss.NewStyle()
)

func makeFgStyle(color string) func(...string) string {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Render
}

func makeFgBgStyle(fg, bg string) func(...string) string {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(fg)).Background(lipgloss.Color(bg)).Render
}

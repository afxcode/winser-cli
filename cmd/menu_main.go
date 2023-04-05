package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/samber/lo"
)

type menuMain struct {
	focusIndex int
	menus      []Menu
}

func newMenuMain() menuMain {
	m := menuMain{
		menus: []Menu{KeyMenuInstall, KeyMenuFind, KeyMenuList},
	}
	return m
}

func (m menuMain) view() string {
	v := appName + "\n\n"

	k := m.menus[m.focusIndex]
	v += lo.Ternary(k == KeyMenuInstall, bgsuccess("> INSTALL SERVICE "), bgsecondary("  INSTALL SERVICE ")) + "\n"
	v += lo.Ternary(k == KeyMenuFind, bgsuccess("> FIND SERVICE    "), bgsecondary("  FIND SERVICE    ")) + "\n"
	v += lo.Ternary(k == KeyMenuList, bgsuccess("> LIST SERVICE    "), bgsecondary("  LIST SERVICE    ")) + "\n\n"
	v += secondary("↑/j: up • ↓/k: down • ctrl+c: quit")
	return v
}

func (m *menuMain) event(msg tea.Msg) (menu Menu, cmd tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "j":
			m.focusIndex--
			if m.focusIndex < 0 {
				m.focusIndex = 0
			}
		case "down", "k":
			m.focusIndex++
			if m.focusIndex >= len(m.menus) {
				m.focusIndex = len(m.menus) - 1
			}
		case "enter":
			return m.menus[m.focusIndex], nil
		}
	}
	return "", nil
}

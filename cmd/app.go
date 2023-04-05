package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/reflow/indent"
)

var appName = bgprimary(" WINSER-CLI V0.1 ")

type appCLI struct {
	Quitting    bool
	CurrentMenu Menu
	menuMain
	menuInstall
	menuFind
	menuUpdate
	menuList
}

func (a appCLI) Init() tea.Cmd {
	return nil
}

type Menu string

const (
	KeyMenuMain    = "main"
	KeyMenuInstall = "install"
	KeyMenuFind    = "find"
	KeyMenuUpdate  = "update"
	KeyMenuList    = "list"
)

func (a appCLI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var k string
	if msg, ok := msg.(tea.KeyMsg); ok {
		k = msg.String()
	}

	// Make sure these keys always quit
	if k == "ctrl+c" {
		a.Quitting = true
		return a, tea.Quit
	}

	switch a.CurrentMenu {
	case KeyMenuMain:
		menu, cmd := a.menuMain.event(msg)
		if menu != "" {
			a.CurrentMenu = menu
		}
		return a, cmd
	case KeyMenuInstall:
		if k == "esc" {
			a.CurrentMenu = KeyMenuMain
			a.menuInstall.reset()
			return a, nil
		}
		return a, a.menuInstall.event(msg)
	case KeyMenuFind:
		if k == "esc" {
			a.CurrentMenu = KeyMenuMain
			return a, nil
		}
		return a, a.menuFind.event(msg)
	case KeyMenuUpdate:
		if k == "esc" {
			a.CurrentMenu = KeyMenuFind
			return a, nil
		}
		return a, a.menuUpdate.event(msg)
	case KeyMenuList:
		if k == "esc" {
			a.CurrentMenu = KeyMenuMain
			return a, nil
		}
		return a, a.menuList.event(msg)
	}

	return a, nil
}

func (a appCLI) View() string {
	var s string
	if a.Quitting {
		return "\n  Finished!\n\n"
	}

	switch a.CurrentMenu {
	case KeyMenuMain:
		s = a.menuMain.view()
	case KeyMenuInstall:
		s = a.menuInstall.view()
	case KeyMenuFind:
		s = a.menuFind.view()
	case KeyMenuUpdate:
		s = a.menuUpdate.view()
	case KeyMenuList:
		s = a.menuList.view()
	}
	return indent.String("\n"+s+"\n\n", 2)
}

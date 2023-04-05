package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	app := appCLI{
		Quitting:    false,
		CurrentMenu: KeyMenuMain,
		menuMain:    newMenuMain(),
		menuInstall: newMenuInstall(),
		menuFind:    newMenuFind(),
		menuUpdate:  newMenuUpdate(),
		menuList:    newMenuList(),
	}

	p := tea.NewProgram(app)
	if _, err := p.Run(); err != nil {
		fmt.Println("could not start program:", err)
	}
}

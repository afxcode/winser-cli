package main

import (
	"winser-cli/winser"

	tea "github.com/charmbracelet/bubbletea"
)

type installServiceMsg struct {
	name string
	err  error
}

func installService(name, displayName, description, executable, argument, startMode string) func() tea.Msg {
	return func() tea.Msg {
		err := winser.Install(winser.Service{
			Name:        name,
			DisplayName: displayName,
			Description: description,
			Executable:  executable,
			Argument:    argument,
			StartMode:   winser.StartMode(startMode),
		})
		return installServiceMsg{name, err}
	}
}

type updateServiceMsg struct {
	name string
	err  error
}

func updateService(name, displayName, description, executable, argument, startMode string) func() tea.Msg {
	return func() tea.Msg {
		return updateServiceMsg{}
	}
}

type removeServiceMsg struct {
	name string
	err  error
}

func removeService(name string) func() tea.Msg {
	return func() tea.Msg {
		err := winser.Remove(name)
		return removeServiceMsg{name, err}
	}
}

type findServiceMsg struct {
	service winser.Service
	err     error
}

func findService(name string) func() tea.Msg {
	return func() tea.Msg {
		service, err := winser.Find(name)
		return findServiceMsg{
			service: service,
			err:     err,
		}
	}
}

type getAllServiceMsg struct {
	services []winser.Service
	err      error
}

func getAllService() func() tea.Msg {
	return func() tea.Msg {
		return getAllServiceMsg{}
	}
}

type startCtrlServiceMsg struct {
	name string
	err  error
}

func startCtrlService(name string) func() tea.Msg {
	return func() tea.Msg {
		err := winser.Start(name)
		return startCtrlServiceMsg{name, err}
	}
}

type stopCtrlServiceMsg struct {
	name string
	err  error
}

func stopCtrlService(name string) func() tea.Msg {
	return func() tea.Msg {
		err := winser.Stop(name)
		return stopCtrlServiceMsg{name, err}
	}
}

type pauseCtrlServiceMsg struct {
	name string
	err  error
}

func pauseCtrlService(name string) func() tea.Msg {
	return func() tea.Msg {
		err := winser.Pause(name)
		return pauseCtrlServiceMsg{name, err}
	}
}

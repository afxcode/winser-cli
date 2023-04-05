package main

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/samber/lo"
)

type Input interface {
	Focus(bool) tea.Cmd
	Render() string
	SetValue(string)
}

// InputText
type InputText struct {
	id string
	textinput.Model
}

func (t *InputText) Focus(focus bool) tea.Cmd {
	if focus {
		// Set focused state
		cmd := t.Model.Focus()
		t.Model.PromptStyle = focusedStyle
		t.Model.TextStyle = focusedStyle
		return cmd
	}
	// Remove focused state
	t.Model.Blur()
	t.Model.PromptStyle = noStyle
	t.Model.TextStyle = noStyle
	return nil
}
func (t *InputText) Render() string { return t.View() }

// InputSelect
type InputSelect struct {
	id             string
	label          string
	options        []option
	selectedOption option
	cursor         int
	focused        bool
}

type option struct {
	label string
	value string
}

func (s *InputSelect) Render() string {
	var b strings.Builder
	label := secondary(" " + s.label + ": ")
	if s.focused {
		b.WriteString(primary(">") + label)
	} else {
		b.WriteString(">" + label)
	}

	for i := 0; i < len(s.options); i++ {
		if s.cursor == i {
			b.WriteString(info(" [" + s.options[i].label + "] "))
		} else {
			b.WriteString("  " + s.options[i].label + "  ")
		}
		b.WriteString("  ")
	}
	return b.String()
}

func (s *InputSelect) Focus(focus bool) tea.Cmd {
	s.focused = focus
	return nil
}

func (s *InputSelect) NextCursor() {
	s.cursor++
	if s.cursor >= len(s.options) {
		s.cursor = len(s.options) - 1
	}
	s.selectedOption = s.options[s.cursor]
}

func (s *InputSelect) PrevCursor() {
	s.cursor--
	if s.cursor <= 0 {
		s.cursor = 0
	}
	s.selectedOption = s.options[s.cursor]
}

func (s *InputSelect) SetValue(v string) {
	opt, i, _ := lo.FindIndexOf(s.options, func(el option) bool { return el.value == v })
	s.cursor = i
	if i < 0 {
		i = 0
	}
	s.cursor = i
	s.selectedOption = opt
	return
}

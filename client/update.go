package client

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/cursor"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewport.SetWidth(msg.Width)
		m.prompt.SetWidth(msg.Width)
		m.viewport.SetHeight(msg.Height - m.prompt.Height())

		if len(m.messages) > 0 {
			// Wrap content before setting it.
			m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width()).Render(strings.Join(m.messages, "\n")))
		}
		m.viewport.GotoBottom()

	case incomingMsg:
		m.messages = append(m.messages, string(msg))
		m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width()).Render(strings.Join(m.messages, "\n")))
		m.viewport.GotoBottom()
		return m, listen(m.connection.conn)

	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			fmt.Println(m.prompt.Value())
			return m, tea.Quit

		case "enter":
			m.connection.send(m.prompt.Value() + "\n")
			m.messages = append(m.messages, m.senderStyle.Render("You: ")+m.prompt.Value())
			m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width()).Render(strings.Join(m.messages, "\n")))
			m.prompt.Reset()
			m.viewport.GotoBottom()
			return m, nil

		default:
			// Send all other keypresses to the textarea.
			var cmd tea.Cmd
			m.prompt, cmd = m.prompt.Update(msg)
			return m, cmd
		}

	case cursor.BlinkMsg:
		// Textarea should also process cursor blinks.
		var cmd tea.Cmd
		m.prompt, cmd = m.prompt.Update(msg)
		return m, cmd
	}

	return m, nil
}

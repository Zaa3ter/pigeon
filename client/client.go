package client

import (
	"fmt"
	"net"
	"os"

	"charm.land/bubbles/v2/textarea"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

const MSGLENGTH = 280

func Connect(addres string) {
	conn, err := net.Dial("tcp", addres)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	c := connection{conn}
	c.login()
	app := tea.NewProgram(initialModel(c))
	if _, err = app.Run(); err != nil {
		fmt.Println(err)
		os.Exit(3)
	}
}

type connection struct {
	conn net.Conn
}

func (c connection) send(msg string) {
	fmt.Fprint(c.conn, msg)
}

func (c connection) login() {
	input := make([]byte, 10)
	ok := make([]byte, 1)
	for {

		fmt.Print("Username: ")
		n, err := os.Stdin.Read(input)
		if err != nil {
			fmt.Println(("Sorry...Something Wrong Happen"))
			fmt.Println(err)
			continue
		}

		c.conn.Write(input[:n])
		n, err = c.conn.Read(ok)
		if ok[0] == 1 {
			user = string(input[:n])
			return
		}
		fmt.Println("Username exist tye something else")
	}
}

type incomingMsg string

func listen(conn net.Conn) tea.Cmd {
	return func() tea.Msg {
		buf := make([]byte, MSGLENGTH)

		n, err := conn.Read(buf)
		if err != nil {
			return err
		}

		return incomingMsg(string(buf[:n]))
	}
}

var user string

type model struct {
	connection  connection
	viewport    viewport.Model
	messages    []string
	prompt      textarea.Model
	senderStyle lipgloss.Style
	err         error
}

func initialModel(conn connection) model {
	prompt := textarea.New()
	prompt.Placeholder = "Send a message..."
	prompt.SetVirtualCursor(false)
	prompt.Focus()

	prompt.Prompt = "┃ "
	prompt.CharLimit = MSGLENGTH

	prompt.SetWidth(30)
	prompt.SetHeight(3)

	// Remove cursor line styling
	s := prompt.Styles()
	s.Focused.CursorLine = lipgloss.NewStyle()
	prompt.SetStyles(s)

	prompt.ShowLineNumbers = false

	vp := viewport.New(viewport.WithWidth(30), viewport.WithHeight(5))
	vp.SetContent("Welcome to the chat room! Type a message and press Enter to send.")
	vp.KeyMap.Left.SetEnabled(false)
	vp.KeyMap.Right.SetEnabled(false)

	prompt.KeyMap.InsertNewline.SetEnabled(false)

	return model{
		connection:  conn,
		prompt:      prompt,
		messages:    []string{},
		viewport:    vp,
		senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		err:         nil,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		textarea.Blink,
		listen(m.connection.conn),
	)
}

func (m model) View() tea.View {
	viewportView := m.viewport.View()
	v := tea.NewView(viewportView + "\n" + m.prompt.View())
	c := m.prompt.Cursor()
	if c != nil {
		c.Y += lipgloss.Height(viewportView)
	}
	v.Cursor = c
	v.AltScreen = true
	return v
}

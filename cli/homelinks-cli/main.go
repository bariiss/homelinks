package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type HomeLink struct {
	Name      string `json:"Name"`
	Text      string `json:"Text"`
	URL       string `json:"URL"`
	Color     string `json:"Color"`
	TextColor string `json:"TextColor"`
	AltText   string `json:"AltText"`
}

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type model struct {
	table      table.Model
	selectedURL string // This field will store the URL to display
}
	
func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			if m.table.Focused() {
				m.selectedURL = m.table.SelectedRow()[2] // Update the selectedURL field
			}
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	view := baseStyle.Render(m.table.View()) // Render the table
	if m.selectedURL != "" {
		// If a URL has been selected, add it to the view
		view += "\nSelected link: " + m.selectedURL + "\n"
	}
	return view + "\n" // Return the full view
}


func loadHomeLinks() ([]HomeLink, error) {
	var homeLinks []HomeLink

	// read file ~/.config/homelinks.json
	data, err := os.ReadFile(os.Getenv("HOME") + "/.config/homelinks.json")
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &homeLinks)
	if err != nil {
		return nil, err
	}
	return homeLinks, nil
}

func main() {
	homeLinks, err := loadHomeLinks()
	if err != nil {
		fmt.Printf("Failed to load home links: %v\n", err)
		os.Exit(1)
	}

	columns := []table.Column{
		{Title: "Name", Width: 10},
		{Title: "Text", Width: 15},
		{Title: "URL", Width: 25},
		{Title: "Color", Width: 10},
		{Title: "TextColor", Width: 10},
		{Title: "AltText", Width: 25},
	}

	var rows []table.Row
	for _, link := range homeLinks {
		row := table.Row{link.Name, link.Text, link.URL, link.Color, link.TextColor, link.AltText}
		rows = append(rows, row)
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(len(homeLinks) + 2),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)
	
	p := tea.NewProgram(model{table: t})
	if err := p.Start(); err != nil {
		fmt.Printf("Failed to start program: %v\n", err)
		os.Exit(1)
	}
}

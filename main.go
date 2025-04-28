package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

type step int

const (
	stepSelectEntryType step = iota
	stepInputCurrentGoals
	stepInputCompanyGoals
	stepInputHighlights
	stepInputLearnings
	stepInputImprovements
	stepPreview
)

type model struct {
	step           step
	entryType      string
	currentGoals   []string
	companyGoals   []string
	highlights     []string
	learnings      []string
	improvements   []string
	textarea       textarea.Model
	choices        list.Model
	progress       progress.Model
	err            error
	previewContent string
}

func initialModel() model {
	items := []list.Item{
		listItem("daily"),
		listItem("weekly"),
	}

	choiceList := list.New(items, list.NewDefaultDelegate(), 20, 10)
	choiceList.Title = "Choose journal entry type"

	ta := textarea.New()
	ta.Placeholder = "Write here..."
	ta.Focus()
	ta.Prompt = "" // REMOVE LINE NUMBER / BULLETS IN INPUT AREA

	return model{
		step:    stepSelectEntryType,
		choices: choiceList,
		textarea: ta,
	}
}

type listItem string

func (i listItem) Title() string       { return string(i) }
func (i listItem) Description() string { return "" }
func (i listItem) FilterValue() string { return string(i) }

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.step {
		case stepSelectEntryType:
			switch msg.String() {
			case "enter":
				i := m.choices.Index()
				if i >= 0 && i < len(m.choices.Items()) {
					m.entryType = m.choices.Items()[i].FilterValue()
					m.step = stepInputCurrentGoals
					m.updatePlaceholder()
					m.textarea.SetValue("")
					if m.entryType == "daily" {
						m.progress = progress.New(progress.WithDefaultGradient())
					} else {
						m.progress = progress.New(progress.WithScaledGradient("#8e44ad", "#3498db"))
					}
				}
			default:
				m.choices, cmd = m.choices.Update(msg)
			}

		case stepInputCurrentGoals, stepInputCompanyGoals, stepInputHighlights, stepInputLearnings, stepInputImprovements:
			switch msg.String() {
			case "enter":
				text := strings.TrimSpace(m.textarea.Value())
				if text == "" {
					switch m.step {
					case stepInputCurrentGoals:
						m.step = stepInputCompanyGoals
					case stepInputCompanyGoals:
						m.step = stepInputHighlights
					case stepInputHighlights:
						m.step = stepInputLearnings
					case stepInputLearnings:
						m.step = stepInputImprovements
					case stepInputImprovements:
						m.step = stepPreview
					}
					m.textarea.SetValue("")
					m.updatePlaceholder()
				} else {
					switch m.step {
					case stepInputCurrentGoals:
						m.currentGoals = append(m.currentGoals, text)
					case stepInputCompanyGoals:
						m.companyGoals = append(m.companyGoals, text)
					case stepInputHighlights:
						m.highlights = append(m.highlights, text)
					case stepInputLearnings:
						m.learnings = append(m.learnings, text)
					case stepInputImprovements:
						m.improvements = append(m.improvements, text)
					}
					m.updatePlaceholder()
					m.textarea.SetValue("")
				}
			default:
				m.textarea, cmd = m.textarea.Update(msg)
			}

		case stepPreview:
			switch msg.String() {
			case "s":
				err := saveMarkdown(m)
				if err != nil {
					m.err = err
				} else {
					return m, tea.Quit
				}
			case "e":
				m.step = stepInputCurrentGoals
				m.textarea.SetValue("")
			case "q":
				return m, tea.Quit
			}
		}

	case tea.WindowSizeMsg:
		m.choices.SetSize(msg.Width, msg.Height)
		m.textarea.SetWidth(msg.Width - 10)
	}

	return m, cmd
}

func (m *model) updatePlaceholder() {
	var count int
	switch m.step {
	case stepInputCurrentGoals:
		count = len(m.currentGoals) + 1
	case stepInputCompanyGoals:
		count = len(m.companyGoals) + 1
	case stepInputHighlights:
		count = len(m.highlights) + 1
	case stepInputLearnings:
		count = len(m.learnings) + 1
	case stepInputImprovements:
		count = len(m.improvements) + 1
	}
	m.textarea.Placeholder = fmt.Sprintf("%d Add next item", count)
}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v", m.err)
	}

	header := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00CED1")).Render(`
#   __    __   ______  ______  ______         __   ______   __  __   ______   __   __   ______   __        
#  /\ "-./  \ /\  ___\/\__  _\/\  __ \       /\ \ /\  __ \ /\ \/\ \ /\  == \ /\ "-.\ \ /\  __ \ /\ \       
#  \ \ \-./\ \\ \  __\\/_/\ \/\ \  __ \     _\_\ \\ \ \/\ \\ \ \_\ \\ \  __< \ \ \-.  \\ \  __ \\ \ \____  
#   \ \_\ \ \_\\ \_____\ \ \_\ \ \_\ \_\   /\_____\\ \_____\\ \_____\\ \_\ \_\\ \_\\"\_\\ \_\ \_\\ \_____\ 
#    \/_/  \/_/ \/_____/  \/_/  \/_/\/_/   \/_____/ \/_____/ \/_____/ \/_/ /_/ \/_/ \/_/ \/_/\/_/ \/_____/ 
#                                                                                                          
`)

	progressBar := m.progress.ViewAs(float64(m.step) / float64(stepPreview))

	switch m.step {
	case stepSelectEntryType:
		return header + "\n" + m.choices.View()

	case stepInputCurrentGoals, stepInputCompanyGoals, stepInputHighlights, stepInputLearnings, stepInputImprovements:
		var entries []string
		var title string
		switch m.step {
		case stepInputCurrentGoals:
			entries = m.currentGoals
			title = "Current Goals"
		case stepInputCompanyGoals:
			entries = m.companyGoals
			title = "Company Goals"
		case stepInputHighlights:
			entries = m.highlights
			title = "Highlights"
		case stepInputLearnings:
			entries = m.learnings
			title = "Learnings / Challenges"
		case stepInputImprovements:
			entries = m.improvements
			title = "Improvements for Next Week"
		}

		var listContent string
		for _, item := range entries {
			listContent += fmt.Sprintf("- %s\n", item)
		}

		return lipgloss.JoinVertical(lipgloss.Top,
			header,
			progressBar,
			lipgloss.NewStyle().Bold(true).Render(title),
			listContent,
			m.textarea.View(),
		) + "\n(Enter to add. Empty input to move to next step)"

	case stepPreview:
		md := buildMarkdown(m)
		out, err := glamour.Render(md, "dark")
		if err != nil {
			return fmt.Sprintf("Error rendering markdown preview: %v", err)
		}
		return header + "\n" + out + "\n\n[s] Save   [e] Edit   [q] Quit"
	}

	return "Loading..."
}

func buildMarkdown(m model) string {
	now := time.Now()
	dateStr := now.Format("2006-01-02")
	weekStart := now.AddDate(0, 0, -int(now.Weekday()-1)).Format("2006-01-02")

	currentGoalsList := formatList(m.currentGoals)
	companyGoalsList := formatList(m.companyGoals)

	currentGoalsBullets := formatBulletPoints(m.currentGoals)
	companyGoalsBullets := formatBulletPoints(m.companyGoals)

	return strings.TrimSpace(fmt.Sprintf(`---
date: %s
entryType: %s
weekStart: %s
currentGoals: %s
companyGoals: %s
---

# Weekly Journal

## Current Goals
%s

## Company Goals
%s

## Highlights
%s

## Learnings / Challenges
%s

## Improvements for Next Week
%s
`,
		dateStr,
		m.entryType,
		weekStart,
		currentGoalsList,
		companyGoalsList,
		currentGoalsBullets,
		companyGoalsBullets,
		formatBulletPoints(m.highlights),
		formatBulletPoints(m.learnings),
		formatBulletPoints(m.improvements),
	))
}

func saveMarkdown(m model) error {
	now := time.Now()
	dateStr := now.Format("2006-01-02")
	weekStart := now.AddDate(0, 0, -int(now.Weekday()-1)).Format("2006-01-02")

	dir := filepath.Join(getSaveBasePath(), "week-beginning", weekStart)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	filename := fmt.Sprintf("%s_%s_entry.md", dateStr, m.entryType)
	filepath := filepath.Join(dir, filename)

	content := buildMarkdown(m)
	return os.WriteFile(filepath, []byte(content), 0644)
}

func formatList(list []string) string {
	if len(list) == 0 {
		return "[]"
	}
	return fmt.Sprintf("[%s]", strings.Join(list, ", "))
}

func formatBulletPoints(points []string) string {
	if len(points) == 0 {
		return "* No items recorded\n"
	}
	var formatted string
	for _, p := range points {
		formatted += fmt.Sprintf("- %s\n", p)
	}
	return formatted
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func getSaveBasePath() string {
    path := os.Getenv("METAJOURNAL_SAVE_PATH")
    if path == "" {
        fmt.Println("⚠️ Warning: METAJOURNAL_SAVE_PATH not set, defaulting to ./data/journal")
        path = "./data/journal"
    }
    return path
}
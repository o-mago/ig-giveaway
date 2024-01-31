package main

import (
	"fmt"
	"os"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/exp/maps"
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle  = focusedStyle.Copy()
	noStyle      = lipgloss.NewStyle()

	winnerStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	winnerTextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("220"))

	helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render

	selectedContenderStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	notSelectedContenderStyle = lipgloss.NewStyle()
)

const (
	padding  = 2
	maxWidth = 80
)

type model struct {
	focusIndex             int
	inputs                 []textinput.Model
	progress               progress.Model
	multipleEntries        bool
	focusedMultipleEntries bool
	allContenders          bool
	focusedAllContenders   bool
	selectedContenders     []int
	contenders             []string
	submitted              bool
	percent                float64
	winners                map[string][]string
	finish                 bool
}

func initialModel() *model {
	m := &model{
		inputs:   make([]textinput.Model, 6),
		progress: progress.New(progress.WithDefaultGradient()),
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.Cursor.Style = cursorStyle

		switch i {
		case 0:
			t.Placeholder = "Instagram user name (e.g. techhubjf)"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Placeholder = "Instagram post code (e.g. C2unenNseJB)"
		case 2:
			t.Placeholder = "Graph API Token"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
		case 3:
			t.Placeholder = "Number of mentions (default 3)"
			t.Validate = func(text string) error {
				for _, s := range text {
					if s < '0' || s > '9' {
						return fmt.Errorf("not a number")
					}
				}

				return nil
			}
		case 4:
			t.Placeholder = "Number of winners (default 1)"
			t.Validate = func(text string) error {
				for _, s := range text {
					if s < '0' || s > '9' {
						return fmt.Errorf("not a number")
					}
				}

				return nil
			}
		case 5:
			t.Placeholder = "Blocklist (comma separated)"
		}

		m.inputs[i] = t
	}

	return m
}

func (m *model) Init() tea.Cmd {
	return textinput.Blink
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			if s == "enter" && m.focusIndex >= len(m.inputs)+2 {
				m.submitted = true

				m.percent = 0
				m.progress.SetPercent(m.percent)
				m.winners = giveaway{}

				if m.inputs[3].Value() == "" {
					m.inputs[3].SetValue("3")
				}

				if m.inputs[4].Value() == "" {
					m.inputs[4].SetValue("1")
				}

				totalMentions, err := strconv.Atoi(m.inputs[3].Value())
				if err != nil {
					panic(err)
				}

				totalWinners, err := strconv.Atoi(m.inputs[4].Value())
				if err != nil {
					panic(err)
				}

				if m.focusIndex == len(m.inputs)+3 {
					updatedBlockList := m.inputs[5].Value()

					if m.inputs[5].Value() != "" {
						updatedBlockList += ","
					}

					for username := range m.winners {
						updatedBlockList += username + ","
					}

					updatedBlockList = strings.TrimSuffix(updatedBlockList, ",")

					m.inputs[5].SetValue(updatedBlockList)
				}

				blockList := strings.Split(m.inputs[5].Value(), ",")

				m.selectedContenders = make([]int, totalWinners)
				for i := range m.selectedContenders {
					m.selectedContenders[i] = -1
				}

				input := startGiveawayInput{
					userName:        m.inputs[0].Value(),
					postCode:        m.inputs[1].Value(),
					token:           m.inputs[2].Value(),
					totalMentions:   totalMentions,
					totalWinners:    totalWinners,
					blockList:       blockList,
					multipleEntries: m.multipleEntries,
					allContenders:   m.allContenders,
				}

				go m.startGiveaway(input)

				return m, tickCmd()
			}

			if s == "enter" && m.focusIndex == len(m.inputs) {
				m.multipleEntries = !m.multipleEntries
				return m, nil
			}

			if s == "enter" && m.focusIndex == len(m.inputs)+1 {
				m.allContenders = !m.allContenders
				return m, nil
			}

			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs)+2 {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					// Set focused state
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = noStyle
				m.inputs[i].TextStyle = noStyle
			}

			if m.focusIndex == len(m.inputs) {
				m.focusedMultipleEntries = true
			} else {
				m.focusedMultipleEntries = false
			}

			if m.focusIndex == len(m.inputs)+1 {
				m.focusedAllContenders = true
			} else {
				m.focusedAllContenders = false
			}

			return m, tea.Batch(cmds...)
		}
	case tickMsg:
		cmd := m.progress.SetPercent(m.percent)
		return m, tea.Batch(tickCmd(), cmd)

	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width - padding*2 - 4
		if m.progress.Width > maxWidth {
			m.progress.Width = maxWidth
		}
		return m, nil

	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd

	default:
		return m, nil
	}

	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m *model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m *model) View() string {
	var b strings.Builder

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	multipleCheck := checkbox("Multiple entries per user (every x mentions, 1 entry)", m.multipleEntries, m.focusedMultipleEntries)
	b.WriteString(multipleCheck)

	allContendersCheck := checkbox("Show all contenders animation", m.allContenders, m.focusedAllContenders)
	b.WriteString(allContendersCheck)

	buttonText := "Submit"

	if m.finish {
		buttonText = "Repeat"
	}

	submitButton := fmt.Sprintf("[ %s ]", blurredStyle.Render(buttonText))
	if m.focusIndex == len(m.inputs)+2 {
		submitButton = focusedStyle.Render(fmt.Sprintf("[ %s ]", buttonText))
	}
	fmt.Fprintf(&b, "\n\n%s", submitButton)

	if m.finish {
		repeatButtonText := "Repeat without last winners"

		repeatButton := fmt.Sprintf("[ %s ]", blurredStyle.Render(repeatButtonText))
		if m.focusIndex == len(m.inputs)+3 {
			repeatButton = focusedStyle.Render(fmt.Sprintf("[ %s ]", repeatButtonText))
		}
		fmt.Fprintf(&b, "	%s", repeatButton)
	}

	b.WriteString("\n\n")

	if m.submitted {
		b.WriteString(m.progress.View() + "\n\n")
	}

	if m.allContenders {
		for index, contender := range m.contenders {
			if index%6 == 0 {
				b.WriteString("\n")
			}

			if slices.Contains(m.selectedContenders, index) {
				b.WriteString(selectedContenderStyle.Render("@"+contender) + " ")

				continue
			}

			b.WriteString(notSelectedContenderStyle.Render("@"+contender) + " ")
		}

		b.WriteString("\n\n")
	}

	if len(m.winners) > 0 {
		userNames := maps.Keys(m.winners)

		// We must sort a slice to show, because maps don't have a fixed order of their keys
		// This way we avoid the screen to keep change winners position
		sort.Strings(userNames)

		for _, userName := range userNames {
			winner := winnerStyle.Render("@" + userName)
			winnerText := winnerTextStyle.Render(m.winners[userName]...)

			finishMessage := fmt.Sprintf("\nThe winner was: %s\nThe mentions were: %s", winner, winnerText)

			b.WriteString(finishMessage + "\n\n")
		}
	}

	b.WriteString(helpStyle("Press esc to quit") + "\n")

	return b.String()
}

func checkbox(text string, conditional, focused bool) string {
	check := " "
	if conditional {
		check = "x"
	}

	check = fmt.Sprintf("\n[%s] ", check)

	checkPlaceholder := text

	if focused {
		check = focusedStyle.Render(check)
		checkPlaceholder = focusedStyle.Render(checkPlaceholder)
	} else {
		checkPlaceholder = blurredStyle.Render(checkPlaceholder)
	}

	return check + checkPlaceholder
}

type tickMsg time.Time

func tickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*200, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func main() {
	if _, err := tea.NewProgram(initialModel()).Run(); err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}
}

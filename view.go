package main

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// ---------------- STYLES ----------------
var letterStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#e5e9f0"))

var selectedLetterStyle = lipgloss.NewStyle().
	Inherit(letterStyle).
	Underline(true)

var correctLetterStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color(
		"#4c566a",
	))

var incorrectLetterStyle = lipgloss.NewStyle().
	Inherit(letterStyle).
	Foreground(lipgloss.Color(
		"#bf616a",
	))

var paragraphStyle = lipgloss.NewStyle().
	Align(lipgloss.Center)

// ---------------- VIEW ----------------

// Render words with styles
func (m model) renderWords() string {
	var result string
	for i, word := range m.words {
		for j, char := range word {
			if char.state == Incomplete {
				if m.wordPos == i && m.charPos == j {
					result += selectedLetterStyle.Render(string(char.showChar))
				} else {
					result += letterStyle.Render(string(char.showChar))
				}
			} else if char.state == Correct {
				result += correctLetterStyle.Render(string(char.showChar))
			} else if char.state == Incorrect {
				result += incorrectLetterStyle.Render(string(char.showChar))
			}
		}
		result += " "
	}
	return result
}

func (m model) View() tea.View {
	doneness := ""
	if m.done {
		doneness = "DONE LMAO"
	}

	// timer := ""

	// Final layouting
	wordField := m.renderWords()
	return tea.NewView(lipgloss.JoinVertical(
		lipgloss.Center,
		wordField,
		doneness,
	))
}

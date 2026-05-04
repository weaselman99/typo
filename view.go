package main

import (
	"fmt"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// ---------------- STYLES ----------------
// Colours
var (
	defaultTextColour     = lipgloss.Color("#5e81ac")
	correctLetterColour   = lipgloss.Color("#3b4252")
	defaultLetterColour   = lipgloss.Color("#e5e9f0")
	incorrectLetterColour = lipgloss.Color("#bf616a")
	borderForeColour      = lipgloss.Color("#5e81ac")
	backgroundColour      = lipgloss.Color("#222222")
)

var letterStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(defaultLetterColour)

var selectedLetterStyle = lipgloss.NewStyle().
	Inherit(letterStyle).
	Underline(true)

var correctLetterStyle = lipgloss.NewStyle().
	Inherit(letterStyle).
	Foreground(correctLetterColour)

var incorrectLetterStyle = lipgloss.NewStyle().
	Inherit(letterStyle).
	Foreground(incorrectLetterColour)

var defaultTextStyle = lipgloss.NewStyle().
	Foreground(defaultTextColour).Background(backgroundColour)

var paragraphStyle = lipgloss.NewStyle().
	Align(lipgloss.Center).
	Margin(1, 2, 1, 2).
	Border(lipgloss.ThickBorder()).
	BorderForeground(borderForeColour)

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

// Render help
func (m model) renderHelp() string {
	var result string
	suffix := "\t"
	helpStrings := []string{
		"[esc] Quit" + suffix,
		"[tab] Reset" + suffix,
		"[1] 10" + suffix,
		"[2] 25" + suffix,
		"[3] 50" + suffix,
		"[4] 100",
	}
	result += defaultTextStyle.Render(lipgloss.JoinHorizontal(lipgloss.Top, helpStrings...))

	return result
}

// Renders the timer and wpm at the bottom of the screen
func (m model) renderStats() string {
	var result string

	if time.Time.IsZero(m.startTime) {
		// Blank defaults
		result += fmt.Sprint("0wpm\t")
		result += fmt.Sprint("0.0s\t")
		result += fmt.Sprint("0%")
	} else {
		wpm, acc := stats(m.words, m.elapsedTimeSeconds)
		// WPM
		result += fmt.Sprintf("%.0fwpm\t", wpm)
		// Timer
		result += fmt.Sprintf("%.1fs\t", m.elapsedTimeSeconds)
		// Accuracy
		result += fmt.Sprintf("%.0f%%", acc*100)
	}

	return defaultTextStyle.Render(result)
}

// Render text field
func (m model) renderParagraph() string {
	return paragraphStyle.Render(lipgloss.JoinVertical(
		lipgloss.Center,
		m.renderHelp(),
		m.renderWords(),
		m.renderStats(),
	))
}

func (m model) View() tea.View {
	doneness := ""
	if m.done {
		doneness = "DONE LMAO"
	}

	// timer := ""

	// Final layouting
	return tea.NewView(lipgloss.JoinVertical(
		lipgloss.Center,
		m.renderParagraph(),
		doneness,
	))
}

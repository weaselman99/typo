package main

import (
	"fmt"
	"time"

	"charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// ---------------- STYLES ----------------
// Colours
var (
	defaultTextColour     = lipgloss.Color("#5e81ac")
	correctLetterColour   = lipgloss.Color("#3b4252")
	defaultLetterColour   = lipgloss.Color("#e5e9f0")
	mutedTextColour       = lipgloss.Color("#c0c0c0")
	incorrectLetterColour = lipgloss.Color("#bf616a")
	borderForeColour      = lipgloss.Color("#5e81ac")
	highlightColour       = lipgloss.Color("#ebcb8b")
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
	Foreground(defaultTextColour)

var areaStyle = lipgloss.NewStyle().
	Margin(1, 1, 1, 1).Align(lipgloss.Left)

var titleStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(defaultLetterColour)

	// Border(lipgloss.RoundedBorder(), true, true, true).
	// BorderForeground(borderForeColour).
	// Padding(0, 1).
	// Margin(0, 16, 0, 0)

var borderStyle = lipgloss.NewStyle().
	Align(lipgloss.Center).
	Border(lipgloss.DoubleBorder()).
	BorderForeground(borderForeColour).
	Width(80)

var paragraphStyle = lipgloss.NewStyle().
	// Width(80).
	// Align(lipgloss.Center).
	Margin(1, 0)

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
	style := defaultTextStyle.MarginRight(2)

	helpItems := []string{
		"[esc] quit",
		"[tab] reset",
	}

	var rendered []string
	rendered = append(rendered, titleStyle.MarginRight(2).Render("WeaselTypo"))
	for _, s := range helpItems {
		rendered = append(rendered, style.Render(s))
	}

	// Word count item with highlighted value
	label := defaultTextStyle.Render("[1-4] word count:")
	value := style.Bold(true).Foreground(mutedTextColour).Render(fmt.Sprintf(" %d", m.wordAmt))
	rendered = append(rendered, label+value)

	return lipgloss.JoinHorizontal(lipgloss.Top, rendered...)
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

	if m.done {
		return defaultTextStyle.Bold(true).Foreground(highlightColour).Render(result)
	} else {
		return defaultTextStyle.Render(result)
	}
}

// Render text field
func (m model) renderParagraph() string {
	innerWidth := 64 // borderStyle width(60) - 2 padding sides - 2 border chars

	help := defaultTextStyle.Width(innerWidth).Align(lipgloss.Center).Render(m.renderHelp())
	words := paragraphStyle.Width(innerWidth).Align(lipgloss.Left).Render(m.renderWords())
	stats := defaultTextStyle.Width(innerWidth).Align(lipgloss.Center).Render(m.renderStats())

	return lipgloss.JoinVertical(lipgloss.Left, help, words, stats)
}

// Render app
func (m model) renderAll() string {
	var status string
	statusStyle := letterStyle.Border(lipgloss.RoundedBorder(), true, true, false).BorderForeground(borderForeColour).Padding(0, 1)
	if time.Time.IsZero(m.startTime) {
		status = statusStyle.Foreground(highlightColour).Render("Start Typing!")
	} else if !m.done {
		status = statusStyle.Foreground(correctLetterColour).BorderForeground(correctLetterColour).Render("...")
	} else {
		status = statusStyle.Foreground(highlightColour).Render("Done!")
	}

	letterStyle.Foreground(lipgloss.Color("#ebcb8b")).Render(status)
	return areaStyle.Render(lipgloss.JoinVertical(
		lipgloss.Center,
		status,
		borderStyle.Render(m.renderParagraph()),
	))
}

func (m model) View() tea.View {
	v := tea.NewView(m.renderAll())
	v.AltScreen = true
	return v
}

package main

import (
	"slices"
	"time"

	tea "charm.land/bubbletea/v2"
)

// Main update function
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKey(msg)
	case TickMsg:
		if !m.done && !m.startTime.IsZero() {
			m.elapsedTimeSeconds = float64(time.Since(m.startTime).Milliseconds()) / 1000
		}
		return m, doTick()
	}

	return m, nil
}

// Get current letter for mutation
func (m *model) currentLetter() *letter {
	return &m.words[m.wordPos][m.charPos]
}

// Handles any keypress event
func (m model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.(tea.KeyMsg).Key()
	switch {
	case key.Code == tea.KeyEsc: // ESC - Quit app
		return m, tea.Quit
	case key.Code == tea.KeyTab: // TAB - Reset typing field
		m.resetModel()
	case key.Code == tea.KeySpace: // SPACE - Next word
		return m.handleSpace()
	case key.Code == tea.KeyBackspace: // BACKSPACE - delete word
		return m.handleBackspace()
	case len(key.Text) == 1: // Character key press
		if !m.done {
			return m.handleRune(msg)
		}
	}
	return m, nil
}

// Handles backspace keypress
func (m model) handleBackspace() (tea.Model, tea.Cmd) {
	// Delete previous word if user hasn't stated the next one,
	// only if the previous word had errorrs
	if m.charPos == 0 && m.wordPos != 0 && !allCorrect(m.words[m.wordPos-1]) {
		m.wordPos--
		m.words[m.wordPos] = slices.Clone(m.originalWords[m.wordPos])
	} else {
		m.words[m.wordPos] = slices.Clone(m.originalWords[m.wordPos])
		m.charPos = 0
	}
	return m, nil
}

// Handles space keypress
func (m model) handleSpace() (tea.Model, tea.Cmd) {
	if m.charPos > 0 {
		// check if at the last word
		if m.wordPos >= len(m.words)-1 {
			m.done = true
			return m, nil
		} else {
			// else move to next word
			m.wordPos++
			m.charPos = 0
		}
	}
	return m, nil
}

// Handles the default key press
func (m model) handleRune(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.(tea.KeyMsg).Key()
	// Check if we have started yet
	if m.startTime.IsZero() {
		m.startTime = time.Now()
	}
	// Check if its a single char
	input := rune(key.Text[0])

	// Numbers to change word amounts
	wordAmtValues := map[rune]int{
		'1': 10,
		'2': 25,
		'3': 50,
		'4': 100,
	}
	if val, ok := wordAmtValues[input]; ok {
		if val != m.wordAmt {
			m.wordAmt = val
			m.resetModel()
		}
		return m, nil
	}

	// Check if inputs are currently inbounds
	if m.charPos >= len(m.originalWords[m.wordPos]) {
		// out of bounds - add incorrect char
		m.words[m.wordPos] = append(m.words[m.wordPos], letter{
			state:    Incorrect,
			realChar: 0,
			showChar: input,
		})
	} else {
		// within bounds - check if correct
		curr := m.currentLetter()
		curr.showChar = input
		if input != curr.realChar {
			curr.state = Incorrect
		} else {
			curr.state = Correct
		}
		m.charPos++
	}

	// set done if lsat word is correct
	if m.wordPos == len(m.words)-1 && m.charPos == len(m.words[m.wordPos]) && allCorrect(m.words[m.wordPos]) {
		m.done = true
	}
	return m, nil
}

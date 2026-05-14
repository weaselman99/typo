package main

import (
	"slices"
	"time"

	tea "charm.land/bubbletea/v2"
)

// ---------------- Model ----------------
type model struct {
	wordAmt            int        // How many words should be typed per session, default 10
	originalWords      [][]letter // Used to reset words easily
	words              [][]letter // Tracks current progress of user
	wordPos            int        // Index of word in words
	charPos            int        // Index of letter in word
	done               bool       // Used to stop timers and wpm
	startTime          time.Time  // Default to time.Time{}
	elapsedTimeSeconds float64    // Time since startTime
}

// Used to deep duplicate the [][]letter slice
func copyWords(words [][]letter) [][]letter {
	copiedWords := make([][]letter, len(words))
	for i := range words {
		copiedWords[i] = make([]letter, len(words[i]))
		copiedWords[i] = slices.Clone(words[i])
	}
	return copiedWords
}

// Default state of the program, we call this when tab is pressed and in initModel
func (m *model) resetModel() {
	// Get random words and copy
	randomisedWords := getWords(m.wordAmt)
	copiedWords := copyWords(randomisedWords)

	// Allocate model states
	m.originalWords = randomisedWords
	m.words = copiedWords
	m.wordPos = 0
	m.charPos = 0
	m.startTime = time.Time{}
	m.elapsedTimeSeconds = 0
	m.done = false
}

// Initialises the model with default values
func initModel() model {
	defaultAmount := 10

	m := model{
		wordAmt: defaultAmount,
	}

	// Add other defaults - we only want to short persist the
	// amount if we change it later and other states if we add more
	m.resetModel()

	return m
}

// For the timer
type TickMsg time.Time

// Tick command .........
func doTick() tea.Cmd {
	return tea.Tick(time.Millisecond, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

// Run when program starts
func (m model) Init() tea.Cmd {
	// Start ticker to keep track of time
	return doTick()
}

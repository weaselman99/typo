package main

import (
	"math/rand"
	"strings"
)

type letterState int

const (
	Incomplete letterState = iota
	Correct
	Incorrect
)

type letter struct {
	realChar rune
	showChar rune
	state    letterState
}

// Resets the letters of a word to incomplete state
func resetWord(word []letter) {
	for _, char := range word {
		char.state = Incomplete
	}
}

// stats() returns the calculated wpm and accuracy from completed session
func stats(words [][]letter, timePassed float64) (float64, float64) {
	if timePassed <= 0 {
		return 0, 0
	}

	// Increment the correct values for wpm and acc calculation
	var numCorrect, numTotal int
	for _, word := range words {
		for _, curr := range word {
			numTotal++
			if curr.state == Correct {
				numCorrect++
			}
		}
	}

	wpm := (float64(numCorrect) / 5) / (timePassed / 60)
	acc := float64(numCorrect) / float64(numTotal)
	return wpm, acc
}

// Returns an array of []letter words from a []string
func castToLetters(words []string) [][]letter {
	// Allocate words
	result := make([][]letter, len(words))
	for i, word := range words {
		// Allocate letters
		result[i] = make([]letter, len(word))

		// cast to letter
		for j, char := range word {
			result[i][j] = letter{
				realChar: char,
				showChar: char,
				state:    Incomplete,
			}
		}
	}
	return result
}

// checks if the word is all correct
func allCorrect(word []letter) bool {
	for _, char := range word {
		if char.state == Incorrect || char.state == Incomplete {
			return false
		}
	}
	return true
}

// getWords rrturns a slice of the wordlist based on the word amount chosen
func getWords(amount int) [][]letter {
	shuffleWords(wordlist)
	return castToLetters(wordlist[0:amount])
}

// shuffleWords randomises the positions of each word in wordlist
func shuffleWords(words []string) {
	rand.Shuffle(len(words), func(i, j int) {
		words[i], words[j] = words[j], words[i]
	})
}

var wordlist = strings.Fields(`
	the be of and a to in he have it that for they with as not on she at by
	you do but from or which one would all will there say who make when can
	more if no man out other so what time up go about than into could state
	only new year some take come these know see use get like then first any
	work now may such give over think most even find day also after way eye 
	many must look before great back through long where much should well
	people down own just because good each those feel seem how high too end
	place little world very still nation hand old life tell write become
	house both between need mean call develop under last right move thing
	general school never same another begin while number part turn real leave
	might want point form off child few small since against ask late home
	interest large person open public follow during present without again
	hold govern around possible head consider word program problem however
	lead system set order plan run keep face fact group play stand increase
	early course change help line this we here show skibidi sigma timo
`)

package random

import (
	"math/rand"
	"time"
)

const (
	digits     = "0123456789"
	consonants = "bcdfghjklmnpqrstvwxz"
	vowels     = "aeiouy"
)

func init() {
	rand.Seed(time.Now().UnixMicro())
}

//Generate random string with specified lenght
func String(lenght int) string {

	var word = make([]rune, lenght)
	var vovelFlag = rand.Intn(2)
	for i := 0; i < lenght; i++ {
		if vovelFlag > 0 {
			word[i] = randomConsonant()
			vovelFlag = 0
		} else {
			word[i] = randomVowel()
			vovelFlag = 1
		}
	}

	word = randomUpperCase(word, 2)

	return string(word)
}

func randomUpperCase(word []rune, count int) []rune {
	for i := 0; i < len(word); i++ {
		if count == 0 {
			break
		}
		if rand.Intn(2) > 0 {
			word[i] = word[i] - 32
			count = count - 1
		}
	}

	return word
}

const lenD, lenC, lenV = len(digits), len(consonants), len(vowels)

//Return random digit rune
func randomDigit() rune {
	var i = rand.Intn(lenD)

	return rune(digits[i])
}

//Return random consonant rune
func randomConsonant() rune {
	var i = rand.Intn(lenC)

	return rune(consonants[i])
}

//Return random vowel rune
func randomVowel() rune {
	var i = rand.Intn(lenV)

	return rune(vowels[i])
}

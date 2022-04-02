package random

import (
	"math/rand"
	"time"
)

const (
	digits     = "0123456789"
	consonants = "bcdfghjklmnpqrstvwxz"
	vowel      = "aeiouy"
)

func init() {
	rand.Seed(time.Now().UnixMicro())
}

func String(lenght int) string {

	var word = make([]byte, lenght)
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

	var uppercaseCount = 1
	for i := 0; i < lenght; i++ {
		if uppercaseCount == 0 {
			break
		}
		if rand.Intn(2) > 0 {
			word[i] = word[i] - 32
			uppercaseCount = uppercaseCount - 1
		}
	}

	return string(word)
}

func randomDigit() byte {
	return digits[rand.Intn(len(digits))]
}

func randomConsonant() byte {
	return consonants[rand.Intn(len(consonants))]
}

func randomVowel() byte {
	return vowel[rand.Intn(len(vowel))]
}

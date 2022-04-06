package random

import (
	"crypto/rand"
	"log"
	"math/big"
)

func String(lenght int) string {

	var word = make([]rune, lenght)

	var cFlag = randomBool()

	for i := 0; i < lenght; i++ {
		if cFlag {
			word[i] = randomConsonant()
			cFlag = false
		} else {
			word[i] = randomVowel()
			cFlag = true
		}
	}

	return string(word)
}

func randomDigit() rune {
	return rune(randomInt(9))
}

func randomBool() bool {
	return randomInt(100) > 50
}

func randomConsonant() rune {
	var consonants = []rune("bcdfghjklmnpqrstvwxz")
	var lenght = len(consonants)
	var n = randomInt(lenght)

	return consonants[n]
}

func randomVowel() rune {
	var vowel = []rune("aeiouy")
	var lenght = len(vowel)
	var n = randomInt(lenght)

	return vowel[n]
}

func randomInt(n int) int {
	var max = int64(n)
	nBig, err := rand.Int(rand.Reader, big.NewInt(max))
	if err != nil {
		log.Panic(err)
	}

	return int(nBig.Int64())
}

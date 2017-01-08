package utils

import (
	"math/rand"
	"time"

	"github.com/swanwish/go-common/logs"
)

const (
	RandomTypeCapitalString = 1 << iota
	RandomTypeLowercaseChar
	RandomTypeDigital
	RandomTypeSymbol

	RandomTypeDefault = RandomTypeCapitalString | RandomTypeLowercaseChar | RandomTypeDigital | RandomTypeSymbol

	DefaultSource = "4he+91Hl3C^Fslhixe-iVJ40!%BI^v6r"
)

var (
	character = [][]string{
		{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"},
		{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"},
		{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"},
		{"!", "@", "#", "$", "%", "^", "&", "*", "(", ")", "-", "+", "_"}}
)

func GenerateRandomDigital(length int64) string {
	return GenerateRandomStringEx(RandomTypeDigital, length)
}

func GenerateRandomMonogram(length int64) string {
	return GenerateRandomStringEx(RandomTypeCapitalString|RandomTypeLowercaseChar, length)
}

func GenerateRandomString(length int64) string {
	return GenerateRandomStringEx(RandomTypeDefault, length)
}

func GenerateRandomStringEx(randomType, length int64) string {
	if randomType == 0 {
		randomType = RandomTypeDefault
	}
	randomCharacters := [][]string{}
	if randomType&RandomTypeCapitalString != 0 {
		randomCharacters = append(randomCharacters, character[0])
	}
	if randomType&RandomTypeLowercaseChar != 0 {
		randomCharacters = append(randomCharacters, character[1])
	}
	if randomType&RandomTypeDigital != 0 {
		randomCharacters = append(randomCharacters, character[2])
	}
	if randomType&RandomTypeSymbol != 0 {
		randomCharacters = append(randomCharacters, character[3])
	}
	typeLength := len(randomCharacters)
	if typeLength == 0 {
		logs.Errorf("Invalid random type %d", randomType)
		return ""
	}
	date := time.Now().Unix()
	source := DefaultSource
	seed := date + getSourceOffset(source)
	if seed == 0 {
		seed = time.Now().Unix()
	}
	if length == 0 {
		length = 8
	}
	rand.Seed(seed)
	randomString := ""
	for i := int64(0); i < length; i++ {
		randomType := rand.Intn(typeLength)
		randomIndex := rand.Intn(len(randomCharacters[randomType]))
		randomString += randomCharacters[randomType][randomIndex]
	}
	return randomString
}

func getSourceOffset(source string) int64 {
	var offset int64
	for i := 0; i < len(source); i++ {
		offset += int64(source[i])
	}
	return offset
}

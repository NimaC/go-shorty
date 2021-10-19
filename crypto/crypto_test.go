package crypto

import (
	"regexp"
	"testing"
)

var stringInput = "www.google.de/testing/plfpelfp?name=shorty&color=purple"

func TestValidCharacters(t *testing.T) {
	IsLettersAndNumbers := regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString
	encodedString := Encode([]byte(stringInput))
	if !IsLettersAndNumbers(encodedString) {
		t.Errorf("Invalid Characters Found In String %q", encodedString)
	}
}

func TestCollision(t *testing.T) {
	encodedString := Encode([]byte(stringInput))
	sameEncodedString := Encode([]byte(stringInput))
	if encodedString != sameEncodedString {
		t.Errorf("Identical String Inputs do not Collide %q", stringInput)
	}
}

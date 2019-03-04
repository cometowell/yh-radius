package main

import "strings"


func LeftPadChar(source string, padChar byte, size int) string {
	sourceLength := len(source)
	if sourceLength >= size {
		return source
	}
	return strings.Repeat(string(padChar), size - sourceLength) + source
}

func RightPadChar(source string, padChar byte, size int) string {
	sourceLength := len(source)
	if sourceLength >= size {
		return source
	}
	return source + strings.Repeat(string(padChar), size - sourceLength)
}
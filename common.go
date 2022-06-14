package main

import (
	"os"
	"strconv"
	"strings"

	"github.com/nleeper/goment"
)

var alphabet = map[byte]int{
	'a': 0,
	'b': 1,
	'c': 2,
	'd': 3,
	'e': 4,
	'f': 5,
	'g': 6,
	'h': 7,
	'i': 8,
	'j': 9,
	'k': 10,
	'l': 11,
	'm': 12,
	'n': 13,
	'o': 14,
	'p': 15,
	'q': 16,
	'r': 17,
	's': 18,
	't': 19,
	'u': 20,
	'v': 21,
	'w': 22,
	'x': 23,
	'y': 24,
	'z': 25,
}

func columnToIndex(column string) int {
	index := -1
	noColIdentifiers := len(column)
	for i := 0; i < noColIdentifiers; i++ {
		letter := column[i]
		index = alphabet[(strings.ToLower(string(letter)))[0]]
	}
	return index
}

func cellReturnValue(valueStr string, returnType string) interface{} {
	var returnable interface{}
	if returnType == "n" {
		returnable, _ = strconv.Atoi(valueStr)
	} else if returnType == "b" {
		returnable = (strings.ToLower(valueStr) == "true")
	} else if returnType == "t" {
		returnable = 0
		if len(valueStr) > 0 {
			format := os.Getenv("SHEETS_DATE_FORMAT")
			moment, _ := goment.New(valueStr, format)
			returnable = moment.ToUnix()
		}
	} else {
		returnable = valueStr
	}
	return returnable
}

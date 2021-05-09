package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

// Helper functions
// -------------------------------------------
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func generateCode(n int) string {
	// Regenerate seed so code is unique
	rand.Seed(time.Now().UnixNano())

	// Available characters
	var letters = []rune("02345689") // 1 and 7 removed as can look similar

	// Build string
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}

func isCrossover(p1 string, p2 string, c string) bool {
	return strings.Contains(p1, c) || strings.Contains(p2, c)
}

func checkWin(m string) bool {
	winConditions := []string{
		"1,2,3", "4,5,6", "7,8,9", // Rows
		"1,4,7", "2,5,8", "3,6,9", // Cols
		"1,5,9", "3,5,7", // Diags
	}
	isWin := false
	for _, cond := range winConditions {
		c := strings.Split(cond, ",")
		if strings.Contains(m, c[0]) && strings.Contains(m, c[1]) && strings.Contains(m, c[2]) {
			isWin = true
		}
	}
	return isWin
}

// https://www.golangprograms.com/remove-duplicate-values-from-slice.html
func unique(intSlice []int) []int {
	keys := make(map[int]bool)
	list := []int{}
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// https://stackoverflow.com/questions/37532255/one-liner-to-transform-int-into-string
func arrayToString(a []int, delim string) string {
	return strings.Trim(strings.Replace(fmt.Sprint(a), " ", delim, -1), "[]")
}

package gotils

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

// Perror is Syntax Sugga for panicing on error
func Perror(err error) {
	if err != nil {
		Lerror(err)
		panic(err)
	}
}

// Lerror is Syntax Sugga for logging an error if its there
func Lerror(err error) {
	if err != nil {
		fmt.Printf("[%v] - %+v\n", strings.ToUpper(os.Getenv("ENV")), err)
	}
}

// UnixInMilliseconds returns the current time in milliseconds
func UnixInMilliseconds() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

/**
 * Sudorandom String Generator
 */
// Sudo Random
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"

//PseudorandomStringWithLength creates a random string w/ given length
func PseudorandomStringWithLength(n int) string {
	return randomString(n)
}

func randomString(l int) string {
	// First seed the rand
	rand.Seed(time.Now().UTC().UnixNano())

	// Product the string
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(48, 122))
	}
	return string(bytes)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

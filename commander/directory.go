package cmdr

import (
	"strings"
)

// dots count consecutive dots
func dots(s string) (count int) {
	l := len(s)
	for count = 0; count < l; count++ {
		if s[count] != '.' {
			break
		}
	}
	return
}

// directory copies
func directory(directoryIn []string, args ...string) (directoryOut []string) {
	// work on a copy, leave source alone
	directoryOut = make([]string, len(directoryIn))
	copy(directoryOut, directoryIn)

	if len(args) < 1 {
		return
	}
	first := args[0]
	lenLeadingDots := dots(first)
	if lenLeadingDots > 0 {
		lenArg := len(first)
		lenDir := len(directoryOut)
		if lenArg > lenLeadingDots {
			// remove leading dots from first argument
			args[0] = first[lenLeadingDots:]
		} else {
			// all dots so discard first argument
			args = args[1:]
		}
		if lenLeadingDots > lenDir {
			// discard excess dots
			lenLeadingDots = lenDir
		}
		// remove trailing directory path
		directoryOut = directoryOut[:lenDir-lenLeadingDots]
	}

	// separate remaining arguments from their dotsdc
	var tokens []string
	for _, arg := range args {
		tokens = append(tokens, strings.Split(arg, ".")...)
	}

	for _, token := range tokens {
		// ignore zero length tokens resulting from
		// multiple dot separators
		if len(token) < 1 {
			continue
		}
		directoryOut = append(directoryOut, token)
	}
	return
}

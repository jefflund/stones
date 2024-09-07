package tui

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

func Log(s string, args ...any) string {
	var subject any
	var subjectNP string

	replaceFn := func(match string) string {
		if match[0] == '%' {
			if len(args) == 0 {
				return fmt.Sprintf("%%!%s(MISSING)", match[1:])
			}

			arg := args[0]
			args = args[1:]
			switch match {
			case "%s":
				subject = arg
				subjectNP = getNounPhrase(subject)
				return subjectNP
			case "%v":
				return conjugate(arg, subjectNP)
			case "%o":
				if arg == subject {
					return getReflexivePronoun(subjectNP)
				}
				return getNounPhrase(arg)
			case "%x":
				return fmt.Sprint(arg)
			}
		}

		verb := match[1 : len(match)-1]
		return conjugate(verb, subjectNP)
	}

	s = strings.ReplaceAll(s, "%%", "%")
	s = formatRE.ReplaceAllStringFunc(s, replaceFn)
	s = strings.ToUpper(s[:1]) + s[1:]
	s = ensurePunctuated(s)

	if len(args) > 0 {
		argstrs := make([]string, len(args))
		for i, arg := range args {
			argstrs[i] = fmt.Sprint(arg)
		}
		extra := strings.Join(argstrs, ", ")
		s = fmt.Sprintf("%s%%!(EXTRA %s)", s, extra)
	}

	return s
}

var (
	formatRE            = regexp.MustCompile("%s|%v|%o|%x|<.+?>")
	articles            = []string{"the", "a", "an"}
	esEndings           = []string{"ch", "sh", "ss", "x", "o"}
	endPunctuation      = []string{".", "!", "?"}
	irregularVerbsFirst = map[string]string{
		"be": "am",
	}
	irregularVerbsSecond = map[string]string{
		"be": "are",
	}
	irregularVerbsThird = map[string]string{
		"do":     "does",
		"be":     "is",
		"have":   "has",
		"can":    "can",
		"cannot": "cannot",
		"could":  "could",
		"may":    "may",
		"must":   "must",
		"shall":  "shall",
		"should": "should",
		"will":   "will",
		"would":  "would",
	}
)

func getNounPhrase(arg any) string {
	noun := fmt.Sprint(arg)
	if noun == "I" || noun == "you" || isProper(noun) || includesArticle(noun) {
		return noun
	}
	return fmt.Sprintf("the %s", noun)
}

func isProper(noun string) bool {
	return unicode.IsUpper([]rune(noun)[0])
}

func includesArticle(noun string) bool {
	for _, article := range articles {
		if strings.HasPrefix(noun, fmt.Sprintf("%s ", article)) {
			return true
		}
	}
	return false
}

func conjugate(arg any, subject string) string {
	verb := fmt.Sprint(arg)
	switch subject {
	case "I":
		if conjugated, irregular := irregularVerbsFirst[verb]; irregular {
			return conjugated
		}
		return verb
	case "you":
		if conjugated, irregular := irregularVerbsSecond[verb]; irregular {
			return conjugated
		}
		return verb
	default:
		if conjugated, irregular := irregularVerbsThird[verb]; irregular {
			return conjugated
		}
		for _, ending := range esEndings {
			if strings.HasSuffix(verb, ending) {
				return verb + "es"
			}
		}
		if strings.HasSuffix(verb, "y") {
			return verb[:len(verb)-1] + "ies"
		}
		return verb + "s"
	}
}

func getReflexivePronoun(noun string) string {
	if noun == "I" {
		return "myself"
	} else if noun == "you" {
		return "yourself"
	} else if isProper(noun) {
		return "themself"
	}
	return "itself"
}

func ensurePunctuated(s string) string {
	for _, punctuation := range endPunctuation {
		if strings.HasSuffix(s, punctuation) {
			return s
		}
	}
	return s + "."
}

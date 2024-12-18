package hjkl

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

// NameQuery is an Event querying an Entity for its string name.
type NameQuery struct {
	Field[string]
}

// Log applies a formatting language to create a log message.
//
// The format specifiers include the following:
//
//	%s - subject
//	%o - object
//	%v - verb
//	%x - literal
//
// Additionally, verb literals may be included using the form <verb>.
//
// Each format specifier can be mapped to any arbitrary value, and is converted
// to a string by fmt package. Consequently, format values should implement the
// fmt.Stringer interface to ensure that the values are correctly represented
// in the formatted string.
//
// Example usage:
//
//	Log("%s <hit> %o", hero, bear) yields "You hit the bear."
//	Log("%s %v %o", tiger, verb, hero) yields "The saber-tooth slashes you."
//	Log("%s <hit> %o!", tiger, rabbit) yields "The saber-tooth hits the rabbit!"
//	Log("%s %v %o?", bear, verb, bear) yields "The bear hits itself?"
//	Log("%s <laugh>", unique) yields "Gorp laughs."
//
// Note that if the fmt.String conversion for a value is "you" so that the
// formatter knows which grammatical-person to use. Named monsters should have
// string representations which are capitalized so the formatter knows not to
// add articles to the names.
//
// Also note that if no ending punctuation is given, then a period is added
// automatically. The sentence is also capitalized if was not already.
func Log(s string, args ...any) string {
	var subject any
	var subjectNP string

	replaceFn := func(match string) string {
		if match[0] == '%' {
			if len(args) == 0 {
				return fmt.Sprintf("%%!%c(MISSING)", match[1])
			}

			arg := args[0]
			argstr := getArgstr(arg)
			args = args[1:]
			if argstr == "" {
				return fmt.Sprintf("%%!%c(EMPTY)", match[1])
			}
			switch match {
			case "%s":
				subject = arg
				subjectNP = getNounPhrase(argstr)
				return subjectNP
			case "%v":
				return conjugate(argstr, subjectNP)
			case "%o":
				if arg == subject {
					return getReflexivePronoun(subjectNP)
				}
				return getNounPhrase(argstr)
			case "%x":
				return argstr
			}
		}

		// If match doesn't start with '%', it must be a verb literal.
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
			argstrs[i] = getArgstr(arg)
		}
		extra := strings.Join(argstrs, ", ")
		s = fmt.Sprintf("%s%%!(EXTRA %s)", s, extra)
	}

	return s
}

// Data nneded by Log helper functions. These should be regarded as constants.
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

// getArgstr converts an arbitrary argument to a string.
func getArgstr(arg any) string {
	if m, ok := arg.(Entity); ok {
		if name := Get(m, &NameQuery{}); name != "" {
			return name
		}
	}
	return fmt.Sprint(arg)
}

// getNounPhrase prepends the article 'the' to the noun unless it already
// begins with an article or is a unique name.
func getNounPhrase(noun string) string {
	if noun == "I" || noun == "you" || isProper(noun) || includesArticle(noun) {
		return noun
	}
	return fmt.Sprintf("the %s", noun)
}

// isProper returns true if the noun string is capitalized.
func isProper(noun string) bool {
	return unicode.IsUpper([]rune(noun)[0])
}

// includesArticle returns true if the noun string begins with an article.
func includesArticle(noun string) bool {
	for _, article := range articles {
		if strings.HasPrefix(noun, fmt.Sprintf("%s ", article)) {
			return true
		}
	}
	return false
}

// conjugate returns the conjugated version of a verb for a given subject.
func conjugate(verb, subject string) string {
	switch subject {
	case "I": // First person.
		if conjugated, irregular := irregularVerbsFirst[verb]; irregular {
			return conjugated
		}
		return verb
	case "you": // Second person.
		if conjugated, irregular := irregularVerbsSecond[verb]; irregular {
			return conjugated
		}
		return verb
	default: // Third person.
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

// getReflexivePronoun gets the reflexive pronoun for a noun.
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

// ensurePunctuated appends a '.' to a string unless it is already punctuated.
func ensurePunctuated(s string) string {
	for _, punctuation := range endPunctuation {
		if strings.HasSuffix(s, punctuation) {
			return s
		}
	}
	return s + "."
}

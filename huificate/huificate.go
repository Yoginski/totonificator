package huificate

import (
	"regexp"
	"strings"
)

var PREFIX_TO_SKIP_RE = regexp.MustCompile("^[бвгджзйклмнпрстфхцчшщьъ]+")
var NON_LETTERS_RE = regexp.MustCompile("[^а-яё-]+")
var ONLY_DASHES_RE = regexp.MustCompile("^-*$")
var RULES = map[string]string{"о": "е", "а": "я", "у": "ю", "ы": "и"}

const RULES_VALUES string = "еяюи"
const PREFIX string = "ху"
const VOWELS string = "оеаяуюы"

func Huify(text string, amount int) string {
	huified := TryHuify(text, amount)
	if huified == "" {
		return ""
	}
	return huified
}

func TryHuify(text string, amount int) string {
	words := strings.Fields(text)
	if len(words) < 1 || len(words)-amount > 4 {
		return ""
	}
	var answer []string
	candidate_words := words
	if len(words) > amount {
		candidate_words = words[len(words)-amount:]
	}
	isHuified := false
	for _, word := range candidate_words {
		output, ok := TryHuifyWord(word)
		if len(output) > 0 {
			answer = append(answer, output)
		}
		isHuified = isHuified || ok
	}
	if isHuified {
		return strings.Join(answer, " ")
	} else {
		return ""
	}
}

func IsHuifyApplicable(word string) (*string, bool) {
	// Пропускаем слова с дефисами, у которых после преобразования ничего, кроме дефисов не осталось
	if ONLY_DASHES_RE.MatchString(word) {
		return nil, false
	}
	postfix := PREFIX_TO_SKIP_RE.ReplaceAllString(word, "")
	// Пропускаем уже хуифицированные слова
	if len(postfix) < 6 || word[:4] == PREFIX && strings.Index(RULES_VALUES, string(postfix[2:4])) >= 0 {
		return nil, false
	}

	return &postfix, true
}

func HuifyWord(postfix string) string {
	if _, ok := RULES[postfix[0:2]]; ok {
		if strings.Index(VOWELS, postfix[2:4]) < 0 {
			return PREFIX + RULES[postfix[0:2]] + postfix[2:]
		}
		if huified, ok := RULES[postfix[2:4]]; ok {
			return PREFIX + huified + postfix[4:]
		}
		return PREFIX + postfix[2:]
	}
	return PREFIX + postfix
}

func TryHuifyWord(text string) (string, bool) {
	word := NON_LETTERS_RE.ReplaceAllString(strings.ToLower(text), "")

	// Отдельная обработка слова бот
	if word == "бот" {
		return "хуебот", true
	}

	if strings.Contains(strings.ToLower(word), "егор") {
		return "охуел?!", true
	}

	if postfix, ok := IsHuifyApplicable(word); ok {
		return HuifyWord(*postfix), true
	}

	return word, false
}

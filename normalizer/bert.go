package normalizer

import (
	"unicode"
)

type BertNormalizer struct {
	CleanText          bool // Whether to remove Control characters and all sorts of whitespaces replaced with single ` ` space
	Lowercase          bool // Whether to do lowercase
	HandleChineseChars bool // Whether to put spaces around chinese characters so they get split
	StripAccents       bool // whether to remove accents
}

func NewBertNormalizer(cleanText, lowercase, handleChineseChars, stripAccents bool) BertNormalizer {
	return BertNormalizer{
		CleanText:          cleanText,
		Lowercase:          lowercase,
		HandleChineseChars: handleChineseChars,
		StripAccents:       stripAccents,
	}
}

// IsWhitespace checks whether rune c is a BERT whitespace character
func isWhitespace(c rune) bool {
	switch c {
	case ' ':
		return true
	case '\t':
		return true
	case '\n':
		return true
	case '\r':
		return true
	}
	return unicode.Is(unicode.Zs, c)
}

// IsControl checks whether rune c is a BERT control character
func isControl(c rune) bool {
	switch c {
	case '\t':
		return false
	case '\n':
		return false
	case '\r':
		return false
	}
	return unicode.In(c, unicode.Cc, unicode.Cf)
}

// bpunc is the BERT extension of the Punctuation character range
var bpunc = &unicode.RangeTable{
	R16: []unicode.Range16{
		{0x0021, 0x002f, 1}, // 33-47
		{0x003a, 0x0040, 1}, // 58-64
		{0x005b, 0x0060, 1}, // 91-96
		{0x007b, 0x007e, 1}, // 123-126
	},
	LatinOffset: 4, // All less than 0x00FF
}

// IsPunctuation checks whether rune c is a BERT punctuation character
func isPunctuation(c rune) bool {
	return unicode.In(c, bpunc, unicode.P)
}

// This defines a "chinese character" as anything in the CJK Unicode block:
//   https://en.wikipedia.org/wiki/CJK_Unified_Ideographs_(Unicode_block)
// Note that the CJK Unicode block is NOT all Japanese and Korean characters,
// despite its name. The modern Korean Hangul alphabet is a different block,
// as is Japanese Hiragana and Katakana. Those alphabets are used to write
// space-separated words, so they are not treated specially and handled
// like for all of the other languages.
var cjk = &unicode.RangeTable{
	R16: []unicode.Range16{
		{0x4e00, 0x9fff, 1},
		{0x3400, 0x4dbf, 1},
		{0xf900, 0xfaff, 1},
	},
	R32: []unicode.Range32{
		{Lo: 0x20000, Hi: 0x2a6df, Stride: 1},
		{Lo: 0x2a700, Hi: 0x2b73f, Stride: 1},
		{Lo: 0x2b740, Hi: 0x2b81f, Stride: 1},
		{Lo: 0x2b820, Hi: 0x2ceaf, Stride: 1},
		{Lo: 0x2f800, Hi: 0x2fa1f, Stride: 1},
	},
}

// isChinese validates that rune c is in the CJK range according to BERT spec
func isChinese(c rune) bool {
	return unicode.In(c, cjk, unicode.P)
}

func doCleanText(n Normalized) Normalized {

	s := n.normalizedString.Normalized
	var changeMap []ChangeMap

	// Fisrt, reverse the string
	var oRunes []rune = []rune(s)

	revRunes := make([]rune, 0)
	for i := len(oRunes) - 1; i >= 0; i-- {
		revRunes = append(revRunes, oRunes[i])
	}

	// Then, clean up
	var removed int = 0
	for _, r := range revRunes {
		if r == 0 || r == 0xfffd || isControl(r) {
			removed += 1
		} else {
			if removed > 0 {
				if isWhitespace(r) {
					r = ' '
				}
				changeMap = append(changeMap, ChangeMap{
					RuneVal: string(r),
					Changes: -removed,
				})
				removed = 0
			} else if removed == 0 {
				changeMap = append(changeMap, ChangeMap{
					RuneVal: string(r),
					Changes: 0,
				})
			}
		}
	}

	// Flip back changeMap
	var unrevMap []ChangeMap
	for i := len(changeMap) - 1; i >= 0; i-- {
		unrevMap = append(unrevMap, changeMap[i])
	}

	n.Transform(unrevMap, removed)

	return n
}

func doHandleChineseChars(n Normalized) Normalized {
	var changeMap []ChangeMap
	runes := []rune(n.normalizedString.Normalized)
	for _, r := range runes {
		// padding around chinese char
		if isChinese(r) {
			changeMap = append(changeMap, []ChangeMap{
				{
					RuneVal: string(' '),
					Changes: 1,
				},
				{
					RuneVal: string(r),
					Changes: 0,
				},
				{
					RuneVal: string(' '),
					Changes: 1,
				},
			}...)
		}

		// No changes
		changeMap = append(changeMap, ChangeMap{
			RuneVal: string(r),
			Changes: 0,
		})
	}

	n.Transform(changeMap, 0)

	return n
}

func doLowercase(n Normalized) Normalized {
	n.Lowercase()

	return n
}

func stripAccents(n Normalized) Normalized {
	n.RemoveAccents()

	return n
}

// Normalize implements Normalizer interface for BertNormalizer
func (bn BertNormalizer) Normalize(n Normalized) (Normalized, error) {
	if bn.CleanText {
		n = doCleanText(n)
	}

	if bn.HandleChineseChars {
		n = doHandleChineseChars(n)
	}

	if bn.Lowercase {
		n = doLowercase(n)
	}

	if bn.StripAccents {
		n = stripAccents(n)
	}

	return n, nil
}

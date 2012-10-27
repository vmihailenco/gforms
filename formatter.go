package gforms

import (
	"unicode"
)

type splitter struct {
	s    string
	slen int
	i    int
}

func newSplitter(s string) *splitter {
	return &splitter{
		s:    s,
		slen: len(s),
	}
}

func (s *splitter) next() rune {
	r := rune(s.s[s.i])
	s.i++
	return r
}

func (s *splitter) skipUpper() {
	for s.i < s.slen {
		if unicode.IsLower(s.next()) {
			s.i--
			break
		}
	}
}

func (s *splitter) skipLower() {
	for s.i < s.slen {
		if unicode.IsUpper(s.next()) {
			s.i--
			break
		}
	}
}

func (s *splitter) split() []string {
	words := make([]string, 0)
	if s.s == "" {
		return words
	}

	for s.i < s.slen-1 {
		start := s.i
		r1 := s.next()
		r2 := s.next()

		if unicode.IsUpper(r1) && unicode.IsUpper(r2) {
			s.skipUpper()
			if s.i != s.slen {
				s.i--
			}
		} else {
			s.skipLower()
		}
		words = append(words, s.s[start:s.i])
	}

	if s.i < s.slen {
		words = append(words, s.s[s.i:])
	}

	return words
}

func splitWords(s string) []string {
	splitter := newSplitter(s)
	return splitter.split()
}

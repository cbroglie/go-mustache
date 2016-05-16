package mustache

import (
	"regexp"
)

// stringScanner is modeled off of Ruby's StringScanner, which is the
// foundation of the parser in the canonical implementation of mustache.
type stringScanner struct {
	input string
	pos   int
}

// scan tries to match the pattern at the current position. If there’s a match,
// the scan pointer is advanced to that location.
// It returns a slice containing the match and any other matched subgroups. If
// there is no match, nil is returned.
func (s *stringScanner) scan(re *regexp.Regexp) []string {
	loc := re.FindStringSubmatchIndex(s.input[s.pos:])
	if loc == nil {
		return nil
	}
	if loc[0] != 0 {
		return nil
	}
	numMatches := len(loc) / 2
	matches := make([]string, numMatches)
	for i := 0; i < numMatches; i++ {
		start, stop := loc[i*2], loc[i*2+1]
		if start >= 0 && stop >= 0 {
			matches[i] = s.input[s.pos+start : s.pos+stop]
		}
	}
	s.pos += loc[1]
	return matches
}

// scanUntil scans the string until the pattern is matched, updating the scan
// pointer to that location.
// It returns a slice containing the substring up to and including the end of
// the match, the match, and any other matched subgroups. If there is no match,
// nil is returned.
func (s *stringScanner) scanUntil(re *regexp.Regexp) []string {
	loc := re.FindStringSubmatchIndex(s.input[s.pos:])
	if loc == nil {
		return nil
	}
	numMatches := (len(loc) / 2)
	matches := make([]string, numMatches+1)
	matches[0] = s.input[s.pos : s.pos+loc[1]] // include everything up to and including the match
	for i := 0; i < numMatches; i++ {
		start, stop := loc[i*2], loc[i*2+1]
		if start >= 0 && stop >= 0 {
			matches[i+1] = s.input[s.pos+loc[i*2] : s.pos+loc[i*2+1]]
		}
	}
	s.pos += loc[1]
	return matches
}

// done returns true when the scan pointer has reached the end of the string.
func (s *stringScanner) done() bool {
	return s.pos == len(s.input)
}

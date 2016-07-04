package mustache

import (
	"fmt"
	"regexp"
)

// stringScanner is modeled off of Ruby's StringScanner, which is the
// foundation of the parser in the canonical implementation of mustache.
type stringScanner struct {
	input string
	pos   int
}

// Scan tries to match the pattern at the current position. If there’s a match,
// the scan pointer is advanced to that location.
// It returns a slice containing the match and any other matched subgroups. If
// there is no match, nil is returned.
func (s *stringScanner) Scan(re *regexp.Regexp) []string {
	matches := s.Check(re)
	if matches == nil {
		return nil
	}
	s.pos += len(matches[0])
	return matches
}

// Check tries to match the pattern at the current position without advancing
// the scan pointer.
// It returns a slice containing the match and any other matched subgroups. If
// there is no match, nil is returned.
func (s *stringScanner) Check(re *regexp.Regexp) []string {
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
	return matches
}

// ScanUntil scans the string until the pattern is matched, updating the scan
// pointer to that location.
// It returns a slice containing the substring up to and including the end of
// the match, the match, and any other matched subgroups. If there is no match,
// nil is returned.
func (s *stringScanner) ScanUntil(re *regexp.Regexp) []string {
	matches := s.CheckUntil(re)
	if matches == nil {
		return nil
	}
	s.pos += len(matches[0])
	return matches
}

// CheckUntil scans the string until the pattern is matched, without updating
// the scan pointer.
// It returns a slice containing the substring up to and including the end of
// the match, the match, and any other matched subgroups. If there is no match,
// nil is returned.
func (s *stringScanner) CheckUntil(re *regexp.Regexp) []string {
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
			matches[i+1] = s.input[s.pos+start : s.pos+stop]
		}
	}
	return matches
}

// Done returns true when the scan pointer has reached the end of the string.
func (s *stringScanner) Done() bool {
	return s.pos == len(s.input)
}

// Pos returns the current location of the scan pointer.
func (s *stringScanner) Pos() int {
	return s.pos
}

// SetPos sets the current location of the scan pointer.
func (s *stringScanner) SetPos(pos int) error {
	if pos < 0 || pos > len(s.input) {
		return fmt.Errorf("pos %d is outside the allowed range of [0, %d]", pos, len(s.input))
	}
	s.pos = pos
	return nil
}

// Len returns the length of the input string.
func (s *stringScanner) Len() int {
	return len(s.input)
}

// Substring returns the substring of the input string for the given range.
func (s *stringScanner) Substring(start, end int) (string, error) {
	if start < 0 || start >= len(s.input) {
		return "", fmt.Errorf("start index %d is outside the allowed range of [0, %d)", start, len(s.input))
	}
	if end < 0 || end > len(s.input) {
		return "", fmt.Errorf("end index %d is outside the allowed range of [0, %d]", end, len(s.input))
	}
	if start >= end {
		return "", fmt.Errorf("start index %d is >= end index %d", start, end)
	}
	return s.input[start:end], nil
}

// StartOfLine returns true if the pointer is at the beginning of a line.
func (s *stringScanner) StartOfLine() bool {
	if s.pos == 0 {
		return true
	}
	if s.input[s.pos-1] == '\n' {
		return true
	}
	return false
}

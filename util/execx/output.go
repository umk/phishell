package execx

import "strings"

// ProcessOutput represents lines printed by a process.
type ProcessOutput struct {
	lines []string
	tail  bool
	err   error
}

func (s *ProcessOutput) Get() (output string, tail bool, err error) {
	if s.err != nil {
		return "", false, s.err
	}

	output = strings.Join(s.lines, "\n")

	return output, s.tail, nil
}

func (s *ProcessOutput) Empty() bool {
	if s.err != nil {
		return true
	}

	return len(s.lines) == 0 && !s.tail
}

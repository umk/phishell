package termx

import (
	"sync"
	"time"

	"github.com/briandowns/spinner"
)

var Spinner = spinnerCore{
	spinner: spinner.New(spinner.CharSets[14], 100*time.Millisecond),
}

type spinnerCore struct {
	spinnerMu    sync.Mutex
	spinner      *spinner.Spinner
	spinnerCount int
}

func (s *spinnerCore) Start() {
	s.change(1)
}

func (s *spinnerCore) Stop() {
	s.change(-1)
}

func (s *spinnerCore) change(d int) {
	s.spinnerMu.Lock()
	defer s.spinnerMu.Unlock()

	s.spinnerCount += d

	switch {
	case s.spinnerCount == 0:
		s.spinner.Stop()
	case s.spinnerCount > 0:
		s.spinner.Start()
	case s.spinnerCount < 0:
		panic("spinner count is out of range")
	}
}

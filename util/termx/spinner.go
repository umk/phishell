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
	spinnerMu sync.Mutex
	spinner   *spinner.Spinner

	spinnerCount     int
	suppressionCount int
}

func (s *spinnerCore) Start() {
	s.spinnerMu.Lock()
	defer s.spinnerMu.Unlock()

	s.spinnerCount++
	s.updateSpinner()
}

func (s *spinnerCore) Stop() {
	s.spinnerMu.Lock()
	defer s.spinnerMu.Unlock()

	s.spinnerCount--
	s.updateSpinner()
}

func (s *spinnerCore) Suppress() {
	s.spinnerMu.Lock()
	defer s.spinnerMu.Unlock()

	s.suppressionCount++
	s.updateSpinner()
}

func (s *spinnerCore) Unsuppress() {
	s.spinnerMu.Lock()
	defer s.spinnerMu.Unlock()

	s.suppressionCount--
	s.updateSpinner()
}

func (s *spinnerCore) updateSpinner() {
	if s.spinnerCount < 0 || s.suppressionCount < 0 {
		panic("spinner or suppression count is out of range")
	} else if s.spinnerCount > 0 && s.suppressionCount == 0 {
		s.spinner.Start()
	} else {
		s.spinner.Stop()
	}
}

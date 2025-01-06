package termx

import (
	"sync"
	"time"

	"github.com/briandowns/spinner"
)

var spinnerMu sync.Mutex

var spinner_ = spinner.New(spinner.CharSets[14], 100*time.Millisecond)

var spinnerCount int

func SpinnerStart() {
	changeSpinner(1)
}

func SpinnerStop() {
	changeSpinner(-1)
}

func changeSpinner(d int) {
	spinnerMu.Lock()
	defer spinnerMu.Unlock()

	spinnerCount += d

	switch {
	case spinnerCount == 0:
		spinner_.Stop()
	case spinnerCount > 0:
		spinner_.Start()
	case spinnerCount < 0:
		panic("spinner count is out of range")
	}
}

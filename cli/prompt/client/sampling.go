package client

// Samples implements a cheesy way to estimate the number of bytes per
// token registered through several calls to the LLM service.
type Samples struct {
	bytesPerTok []float32

	n int
}

func newSamples(n int, def float32) *Samples {
	if n <= 0 {
		panic("number of items must be a positive value")
	}

	bytesPerTok := make([]float32, n)
	for i := range bytesPerTok {
		bytesPerTok[i] = def
	}

	return &Samples{
		bytesPerTok: bytesPerTok,
	}
}

func (s *Samples) BytesPerTok() float32 {
	var r float32
	for _, v := range s.bytesPerTok {
		r += v
	}

	return r / float32(len(s.bytesPerTok))
}

func (s *Samples) put(bytesPerTok float32) {
	n := (s.n + 1) % len(s.bytesPerTok)
	s.n = n

	s.bytesPerTok[s.n] = bytesPerTok
}

package splitter

type Document []byte

type Chunk struct {
	Document Document
	Metadata Metadata
	Range    Range
}

type Pos struct {
	Offset int // offset to the character in bytes
}

type Range struct {
	Start Pos
	End   Pos // exclusive
}

func (r *Range) Get(document Document) string {
	return string(document[r.Start.Offset:r.End.Offset])
}

type Metadata = map[string]string

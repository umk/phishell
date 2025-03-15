package vector

import (
	"math"

	"github.com/pehringer/simd"
	"github.com/umk/phishell/util/slicesx"
)

const pooledVectorSize = 20_000

var vectorsPool = slicesx.NewSlicesPool[float32](pooledVectorSize)

func cosineSimilarity(a, b []float32, normA, normB float64, tmp []float32) float64 {
	simd.MulFloat32(a, b, tmp)

	var sum float32
	for _, v := range tmp {
		sum += v
	}

	return float64(sum) / (normA * normB)
}

func vectorNorm(vector []float32, tmp []float32) float64 {
	simd.MulFloat32(vector, vector, tmp)

	var sum float32
	for _, v := range tmp {
		sum += v
	}

	return math.Sqrt(float64(sum))
}

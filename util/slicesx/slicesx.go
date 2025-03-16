package slicesx

// Unique returns a new slice containing only the unique values from the input slice.
// The function preserves the original order of elements.
func Unique[T comparable](slice []T) []T {
	if len(slice) == 0 {
		return []T{}
	}

	seen := make(map[T]bool)
	result := make([]T, 0, len(slice))

	for _, item := range slice {
		if _, exists := seen[item]; !exists {
			seen[item] = true
			result = append(result, item)
		}
	}

	return result
}

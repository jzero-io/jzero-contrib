package slicex

func Paginate[T any](slice []T, page, size int) []T {
	start := (page - 1) * size
	if start >= len(slice) {
		return []T{}
	}

	end := start + size
	if end > len(slice) {
		end = len(slice)
	}

	return slice[start:end]
}

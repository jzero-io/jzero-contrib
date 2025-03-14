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

// ToMap 是一个泛型函数，用于将切片转换为映射
func ToMap[K comparable, T any](rows []T, keyFunc func(row T) K) map[K]T {
	// 初始化一个空的映射，键的类型为 K，值的类型为 T
	res := make(map[K]T)
	// 遍历切片中的每个元素
	for _, row := range rows {
		// 调用 keyFunc 函数获取当前元素对应的键
		key := keyFunc(row)
		// 将键和对应的元素存入映射中
		res[key] = row
	}
	return res
}

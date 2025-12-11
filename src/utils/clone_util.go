package utils

// 一些clone方法，非递归

// 克隆切片的指针
// 切片的元素非引用类型或者指针类型
func CloneSlicePtr[T any](s *[]T) *[]T {
	if s == nil {
		return nil
	}
	res := CloneSlice(*s)
	return &res
}

// 克隆指针的切片
// 指针指向的元素不能是引用类型或者指针类型，并且如果是结构体，属性不能是引用类型或者指针类型
func ClonePtrSlice[T any](s []*T) []*T {
	res := make([]*T, len(s))
	for i, ptr := range s {
		res[i] = ClonePtr(ptr)
	}
	return res
}

// 克隆切片
func CloneSlice[T any](s []T) []T {
	res := make([]T, len(s))
	copy(res, s)
	return res
}

// 克隆指针
func ClonePtr[T any](p *T) *T {
	if p == nil {
		return nil
	}
	v := *p
	return &v
}

// 克隆map
func CloneMap[K comparable, V any](m map[K]V) map[K]V {
	res := make(map[K]V)
	for k, v := range m {
		res[k] = v
	}
	return res
}

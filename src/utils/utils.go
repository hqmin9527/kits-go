package utils

import "reflect"

func If[T any](judge bool, valueA, valueB T) T {
	if judge {
		return valueA
	} else {
		return valueB
	}
}

func Ref[T any](v T) *T {
	return &v
}

func GetOrDefault[T any](v *T, def T) T {
	if v != nil {
		return *v
	}
	return def
}

func GetOrZero[T any](v *T) T {
	if v != nil {
		return *v
	}
	var zero T
	return zero
}

// 通过反射，判断接口的值是否为nil
func IsValueNil(a any) bool {
	if a == nil {
		return true
	}
	v := reflect.ValueOf(a)
	return v.IsNil()
}

package collection

type Number interface {
	~int | ~int32 | ~int64 | ~float64 | ~float32
}

type Direct interface {
	~bool | ~string | ~map[string]any | ~[]any | any
}

func ParseNumber[N Number](val any) N {
	if val == nil {
		return N(0.0)
	}
	return N(val.(float64))
}

func ParseDirect[D Direct](val any) D {
	if val == nil {
		var zero D
		return zero
	}
	return val.(D)
}

func ParseNumberPtr[N Number](val any) *N {
	if val == nil {
		return nil
	}
	v := N(val.(float64))
	return &v
}

func ParseDirectPtr[D Direct](val any) *D {
	if val == nil {
		return nil
	}
	v := val.(D)
	return &v
}

func ParseSliceNumber[N Number](val any) []N {
	if val == nil {
		return nil
	}
	vs := val.([]any)
	ns := make([]N, len(vs))
	for i, v := range vs {
		ns[i] = ParseNumber[N](v)
	}
	return ns
}

func ParseSliceDirect[D Direct](val any) []D {
	if val == nil {
		return nil
	}
	vs := val.([]any)
	ns := make([]D, len(vs))
	for i, v := range vs {
		ns[i] = ParseDirect[D](v)
	}
	return ns
}

type Equal interface {
	~int | ~int32 | ~int64 | ~float64 | ~float32 |
		~bool | ~string
}

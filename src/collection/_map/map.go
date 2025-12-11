package _map

// Golang不支持结构方法使用泛型

type Map map[string]any

func (m Map) GetIntPtr(key string) *int {
	return GetIntPtr(m, key)
}

func (m Map) GetInt64Ptr(key string) *int64 {
	return GetInt64Ptr(m, key)
}

func (m Map) GetStrPtr(key string) *string {
	return GetStrPtr(m, key)
}

func (m Map) GetSlice(key string) []any {
	return GetSlice(m, key)
}

func (m Map) GetSliceStr(key string) []string {
	return GetSliceStr(m, key)
}

func (m Map) GetSlicePtr(key string) *[]any {
	return GetSlicePtr(m, key)
}

func (m Map) GetSliceStrPtr(key string) *[]string {
	return GetSliceStrPtr(m, key)
}

func (m Map) GetSliceMap(key string) []map[string]any {
	return GetSliceMap(m, key)
}

func (m Map) GetFloat32Ptr(key string) *float32 {
	return GetFloat32Ptr(m, key)
}

func (m Map) GetFloat64Ptr(key string) *float64 {
	return GetFloat64Ptr(m, key)
}

func (m Map) GetBoolPtr(key string) *bool {
	return GetBoolPtr(m, key)
}

func (m Map) GetInt(key string) int {
	return GetInt(m, key)
}

func (m Map) GetInt64(key string) int64 {
	return GetInt64(m, key)
}

func (m Map) GetString(key string) string {
	return GetString(m, key)
}

func (m Map) GetBool(key string) bool {
	return GetBool(m, key)
}

func (m Map) GetMap(key string) map[string]any {
	return GetMap(m, key)
}

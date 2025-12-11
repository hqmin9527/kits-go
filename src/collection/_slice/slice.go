package _slice

type Slice []any

func (s Slice) GetIntPtr(idx int) *int {
	return GetIntPtr(s, idx)
}

func (s Slice) GetSliceStr(idx int) []string {
	return GetSliceStr(s, idx)
}

func (s Slice) GetSliceMap(idx int) []map[string]any {
	return GetSliceMap(s, idx)
}

func (s Slice) GetInt(idx int) int {
	return GetInt(s, idx)
}

func (s Slice) GetInt64(idx int) int64 {
	return GetInt64(s, idx)
}

func (s Slice) GetString(idx int) string {
	return GetString(s, idx)
}

func (s Slice) GetBool(idx int) bool {
	return GetBool(s, idx)
}

func (s Slice) GetMap(idx int) map[string]any {
	return GetMap(s, idx)
}

func (s Slice) GetSlice(idx int) []any {
	return GetSlice(s, idx)
}

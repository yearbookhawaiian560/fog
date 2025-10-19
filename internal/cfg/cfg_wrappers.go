package cfg

func Set(key string, val interface{}) error {
	return K.Set(key, val)
}

func Bool(path string) bool {
	return K.Bool(path)
}

func Str(path string) string {
	return K.String(path)
}

func MustStr(path string) string {
	return K.MustString(path)
}

func Strings(path string) []string {
	return K.Strings(path)
}

func Int(path string) int {
	return K.Int(path)
}

func Int64(path string) int64 {
	return K.Int64(path)
}

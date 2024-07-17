package utility

import "unsafe"

func StringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&struct {
		string
		Cap int
	}{s, len(s)}))
}

func ByteToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

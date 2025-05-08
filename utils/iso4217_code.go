package utils

func Iso4217Number2Code(number string) (string, bool) {
	code, ok := iso4217Data[number]
	return code, ok
}

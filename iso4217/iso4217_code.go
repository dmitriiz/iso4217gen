package iso4217

func Iso4217Numeric2Code(numeric string) (string, bool) {
	code, ok := iso4217N2C[numeric]
	return code, ok
}

func Iso4217Code2Numeric(code string) (string, bool) {
	code, ok := iso4217C2N[code]
	return code, ok
}

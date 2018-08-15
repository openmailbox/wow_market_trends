package internal

// CheckError is a global helper that panics on all errors
func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

package helper

// ErrorHandler is a temparory error solution
func ErrorHandler(err error) {
	if err != nil {
		panic(err)
	}
}

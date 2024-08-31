package cmd

type testMemExit struct {
	code int
}

func (e *testMemExit) Exit(i int) {
	e.code = i
}

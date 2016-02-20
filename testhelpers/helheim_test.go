package testhelpers_test

type mockSample struct {
	BCalled chan bool
	BOutput struct {
		Ret0 chan int
		Ret1 chan string
	}
}

func newMockSample() *mockSample {
	m := &mockSample{}
	m.BCalled = make(chan bool, 100)
	m.BOutput.Ret0 = make(chan int, 100)
	m.BOutput.Ret1 = make(chan string, 100)
	return m
}
func (m *mockSample) B() (int, string) {
	m.BCalled <- true
	return <-m.BOutput.Ret0, <-m.BOutput.Ret1
}

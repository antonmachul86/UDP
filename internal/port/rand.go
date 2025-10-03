package port

type Rand interface {
	Read(p []byte) (n int, err error)
	Intn(n int) int
}

package uniprot

/*
Shared interface in the module
*/
type FileCloser interface {
	Close() error
}

type FileWriter interface {
	Write(p []byte) (n int, err error)
	Flush() error
}

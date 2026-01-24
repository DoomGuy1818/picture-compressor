package compressor

type Compressor interface {
	Compress(path string, alias string) (string, error)
}

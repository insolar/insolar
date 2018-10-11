package request

type UpdateNode interface {
	getCurrentVer(address string) (string, error)
	downloadFile(filePath string, url string) error
}

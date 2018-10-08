package request

type RequestUpdateNode interface {
	getCurrentVer(address string) (string, error)
	downloadFile(filePath string, url string) error
}

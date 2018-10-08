package request

type FtpRequestUpdateNode struct {
	RequestUpdateNode
}

func (request FtpRequestUpdateNode) getCurrentVer(address string) (string, error) {

	// todo: return from ftp list, latest version

	return "", nil
}

func (request FtpRequestUpdateNode) downloadFile(filepath string, url string) error {

	// todo: download from ftp latest version

	return nil
}

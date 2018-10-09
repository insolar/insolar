package request

type FtpUpdateNode struct {
	UpdateNode
}

func (request FtpUpdateNode) getCurrentVer(address string) (string, error) {

	// todo: return from ftp list, latest version

	return "", nil
}

func (request FtpUpdateNode) downloadFile(filepath string, url string) error {

	// todo: download from ftp latest version

	return nil
}

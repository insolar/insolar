package request

type SmbUpdateNode struct {
	UpdateNode
}

func (request SmbUpdateNode) getCurrentVer(address string) (string, error) {

	// todo: return from smb list, latest version

	return "", nil
}

func (request SmbUpdateNode) downloadFile(filepath string, url string) error {

	// todo: download from smb latest version

	return nil
}

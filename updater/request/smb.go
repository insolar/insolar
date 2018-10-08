package request

type SmbRequestUpdateNode struct {
	RequestUpdateNode
}

func (request SmbRequestUpdateNode) getCurrentVer(address string) (string, error) {

	// todo: return from smb list, latest version

	return "", nil
}

func (request SmbRequestUpdateNode) downloadFile(filepath string, url string) error {

	// todo: download from smb latest version

	return nil
}

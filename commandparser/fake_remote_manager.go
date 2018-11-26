package commandparser

type FakeRemoteManager struct {
	host     string
	username string
	password string
	FilePath string
}

func (f *FakeRemoteManager) UploadArtifact(filePath string) error {
	f.FilePath = filePath
	return nil
}

func NewFakeRemoteManager() *FakeRemoteManager {
	return &FakeRemoteManager{}
}

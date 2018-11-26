package commandparser

type RemoteManager interface {
	UploadArtifact(filePath string) error
}

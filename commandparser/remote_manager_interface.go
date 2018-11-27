package commandparser

type RemoteManager interface {
	UploadArtifact(source, destination string) error
	ExtractArchive(source, destination string) error
	ExecuteCommand(file string) error
}

package commandparser

type StemcellConstructor struct {
	remoteManager RemoteManager
}

func (p *StemcellConstructor) uploadArtifact(filepath string) error {
	p.remoteManager.UploadArtifact(filepath)
	return nil
}

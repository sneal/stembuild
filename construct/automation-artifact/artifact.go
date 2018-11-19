package artifact

func DownloadAutomationArtifact(version, path string) error {
	refreshToken := "323c5fac66c34dcd85628138d736c8b1-r"

	return nil
}

//func DownloadAutomationArtifact(version, path string) error {
//	refreshToken := "323c5fac66c34dcd85628138d736c8b1-r"
//
//	clientWithToken := NewPivnetClientWithToken(refreshToken, "https://network.pivotal.io")
//	var _ sanitizer.Sanitizer
//
//
//	_, err := os.Create(path + "StemcellAutomation.zip")
//	if err != nil {
//		fmt.Fprintf(os.Stderr, "Can not create the file at given location %s", err.Error())
//	}
//
//
//	_, err = clientWithToken.ReleaseForVersion("stemcells-windows-server", "StemcellAutomation-v0.4")
//	if err != nil {
//		fmt.Fprintf(os.Stderr, "blah %s", err.Error())
//	}
//
//	//err = clientWithToken.DownloadProductFile(file, "stemcells-windows-server", release.ID, 253187, progressWriter)
//	//if err != nil {
//	//	fmt.Fprintf(os.Stderr, "blah %s", err.Error())
//	//}
//
//	return nil
//}

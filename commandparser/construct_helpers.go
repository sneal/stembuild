package commandparser

import (
	//"fmt"
	"os"
	"path/filepath"
)

func IsArtifactInDirectory(directory string, artifactFileName string) (bool, error) {

	if _, directoryErr := os.Stat(directory); os.IsNotExist(directoryErr) {
		return false, directoryErr
	}

	artifactPath := filepath.Join(directory, artifactFileName)

	if _, err := os.Stat(artifactPath); os.IsNotExist(err) {
		return false, nil
	}
	return true, nil
}

func UploadArtifact(rm RemoteManager) error {
	err := rm.UploadArtifact("./LGPO.zip", "C:\\provision\\LGPO.zip")
	if err != nil {
		return err
	}
	err = rm.UploadArtifact("./StemcellAutomation.zip", "C:\\provision\\StemcellAutomation.zip")
	if err != nil {
		return err
	}

	return nil
}

func ExtractArchive(rm RemoteManager) error {
	err := rm.ExtractArchive("C:\\provision\\StemcellAutomation.zip", "C:\\provision\\")
	return err
}
func ExecuteSetupScript(rm RemoteManager) error {
	err := rm.ExecuteCommand("powershell.exe C:\\provision\\Setup.ps1")
	return err
}

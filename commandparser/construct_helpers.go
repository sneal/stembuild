package commandparser

import (
	"github.com/masterzen/winrm"
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

func RemoteExecuteAutomation(username string, password string, ip string) {
	endpoint := winrm.NewEndpoint(ip, 5986, false, true, nil, nil, nil, 0)
	client, err := winrm.NewClient(endpoint, username, password)
	if err != nil {
		panic(err)
	}
	client.Run("Expand-Archive C:\\Provisioners\\StemcellAutomation.zip", os.Stdout, os.Stderr)
	client.Run("./Setup.ps1", os.Stdout, os.Stderr)
}

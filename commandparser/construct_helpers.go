package commandparser

import  (
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
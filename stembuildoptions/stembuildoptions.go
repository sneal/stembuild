package stembuildoptions

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
	"fmt"
	"crypto/sha256"
)

type StembuildOptions struct {
	PatchFile string `yaml:"patch_file"`
	OSVersion string `yaml:"os_version"`
	OutputDir string `yaml:"output_dir"`
	Version   string `yaml:"version"`
	VHDFile   string `yaml:"vhd_file"`
	VMDKFile  string `yaml:"vmdk_file"`
	VHDChecksum string `yaml:"vhd_file_checksum"`
	PatchFileChecksum string `yaml:"patch_file_checksum"`
}

// Copy into `d` the values in `s` which are empty in `d`.
func (d *StembuildOptions) CopyFrom(s StembuildOptions) {
	if d.PatchFile == "" {
		d.PatchFile = s.PatchFile
	}

	if d.OSVersion == "" {
		d.OSVersion = s.OSVersion
	}

	// ignore OutputDir from config file

	if d.Version == "" {
		d.Version = s.Version
	}

	if d.VHDFile == "" {
		d.VHDFile = s.VHDFile
	}

	if d.VMDKFile == "" {
		d.VMDKFile = s.VMDKFile
	}

	if d.VHDChecksum == ""{
		d.VHDChecksum = s.VHDChecksum
	}

	if d.PatchFileChecksum == ""{
		d.PatchFileChecksum = s.PatchFileChecksum
	}

}

func LoadOptionsFromManifest(fileName string, patchArgs *StembuildOptions) error {
	_, err := os.Stat(fileName)
	if err != nil {
		return err
	}

	var patchManifest StembuildOptions
	manifestFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(manifestFile, &patchManifest)
	if err != nil {
		return err
	}

	patchArgs.CopyFrom(patchManifest)

	return nil
}

func (o *StembuildOptions) ValidateChecksum() error {

	errorMessage := ""

	VHDFilePath := o.VHDFile
	VHDFileContents, err := ioutil.ReadFile(VHDFilePath)

	if err != nil {
		errorMessage += fmt.Sprintf("Could not read file or file does not exist: %s\n", VHDFilePath)
	} else {
		VHDChecksumByte := sha256.Sum256(VHDFileContents)

		realVHDChecksum := string(VHDChecksumByte[:32])

		expectedVHDChecksum := o.VHDChecksum

		if realVHDChecksum != expectedVHDChecksum {
			errorMessage += fmt.Sprintf("Actual VHD Checksum %s does not match the expected checksum: %s\n", realVHDChecksum, expectedVHDChecksum)
		}
	}


	PatchFilePath := o.PatchFile
	PatchFileContents, err := ioutil.ReadFile(PatchFilePath)
	if err != nil {
		errorMessage += fmt.Sprintf("Could not read file or file does not exist: %s\n", PatchFilePath)
	} else {
		PatchFileChecksumByte := sha256.Sum256(PatchFileContents)
		realPatchFileChecksum := string(PatchFileChecksumByte[:32])


		expectedPatchFileChecksum := o.PatchFileChecksum


		if realPatchFileChecksum != expectedPatchFileChecksum {
			errorMessage += fmt.Sprintf("Actual PatchFile Checksum %s does not match the expected checksum: %s", realPatchFileChecksum, expectedPatchFileChecksum)
		}
	}

	if errorMessage != "" {
		return fmt.Errorf(errorMessage)
	}

	return nil
}

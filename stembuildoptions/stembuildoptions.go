package stembuildoptions

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type StembuildOptions struct {
	PatchFile         string         `yaml:"patch_file"`
	PatchFileChecksum string `yaml:"patch_file_checksum"`
	OSVersion         string         `yaml:"os_version"`
	OutputDir         string         `yaml:"output_dir"`
	Version           string         `yaml:"version"`
	VHDFile           string         `yaml:"vhd_file"`
	VHDFileChecksum   string `yaml:"vhd_file_checksum"`
	VMDKFile          string         `yaml:"vmdk_file"`
}

// Copy into `d` the values in `s` which are empty in `d`.
func (d *StembuildOptions) CopyFrom(s StembuildOptions) {
	if d.PatchFile == "" {
		d.PatchFile = s.PatchFile
	}

	if d.PatchFileChecksum == "" {
		d.PatchFileChecksum = s.PatchFileChecksum
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

	if d.VHDFileChecksum == "" {
		d.VHDFileChecksum = s.VHDFileChecksum
	}

	if d.VMDKFile == "" {
		d.VMDKFile = s.VMDKFile
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

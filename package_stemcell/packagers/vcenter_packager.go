package packagers

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"

	"github.com/cloudfoundry-incubator/stembuild/filesystem"

	"github.com/cloudfoundry-incubator/stembuild/package_stemcell/config"
)

//go:generate counterfeiter . IaasClient
type IaasClient interface {
	ValidateUrl() error
	ValidateCredentials() error
	FindVM(vmInventoryPath string) error
	ExportVM(vmInventoryPath string, destination string) error
	ListDevices(vmInventoryPath string) ([]string, error)
	RemoveDevice(vmInventoryPath string, deviceName string) error
}

type VCenterPackager struct {
	SourceConfig config.SourceConfig
	OutputConfig config.OutputConfig
	Client       IaasClient
}

func (v VCenterPackager) Package() error {
	fmt.Println(fmt.Sprintf("OS: %s", v.OutputConfig.Os))
	fmt.Println(fmt.Sprintf("Version: %s", v.OutputConfig.StemcellVersion))
	deviceList, err := v.Client.ListDevices(v.SourceConfig.VmInventoryPath)
	if err != nil {
		return err
	}
	for _, deviceName := range deviceList {
		matched, _ := regexp.MatchString("^(floppy-|ethernet-)", deviceName)
		if matched {
			err = v.Client.RemoveDevice(v.SourceConfig.VmInventoryPath, deviceName)
			if err != nil {
				return err
			}
		}
	}

	workingDir, err := ioutil.TempDir(os.TempDir(), "vcenter-packager-working-directory")

	if err != nil {
		return errors.New("failed to create working directory")
	}

	fmt.Println(fmt.Sprintf("working directory: %s", workingDir))
	stemcellDir, err := ioutil.TempDir(os.TempDir(), "vcenter-packager-stemcell-directory")
	if err != nil {
		return errors.New("failed to create stemcell directory")
	}
	fmt.Println(fmt.Sprintf("stemcell directory: %s", stemcellDir))
	err = v.Client.ExportVM(v.SourceConfig.VmInventoryPath, workingDir)

	files, _ := ioutil.ReadDir(workingDir)

	for _, f := range files {
		fmt.Println(fmt.Sprintf("working dir file: %s", f.Name()))
	}

	if err != nil {
		return errors.New("failed to export the prepared VM")
	}

	vmName := path.Base(v.SourceConfig.VmInventoryPath)
	shaSum, err := TarGenerator(filepath.Join(stemcellDir, "image"), filepath.Join(workingDir, vmName))
	manifestContents := CreateManifest(v.OutputConfig.Os, v.OutputConfig.StemcellVersion, shaSum)
	err = WriteManifest(manifestContents, stemcellDir)

	files, _ = ioutil.ReadDir(stemcellDir)

	for _, f := range files {
		fmt.Println(fmt.Sprintf("stemcell dir file: %s", f.Name()))
	}

	if err != nil {
		return errors.New("failed to create stemcell.MF file")
	}

	stemcellFilename := StemcellFilename(v.OutputConfig.StemcellVersion, v.OutputConfig.Os)
	_, err = TarGenerator(filepath.Join(v.OutputConfig.OutputDir, stemcellFilename), stemcellDir)

	return nil
}

func (v VCenterPackager) ValidateFreeSpaceForPackage(fs filesystem.FileSystem) error {
	println(os.Stdout, "WARNING: Please make sure you have enough free disk space for export")
	return nil
}

func (v VCenterPackager) ValidateSourceParameters() error {
	err := v.Client.ValidateUrl()
	if err != nil {
		return errors.New("please provide a valid vCenter URL")
	}

	err = v.Client.ValidateCredentials()
	if err != nil {
		errMsg := fmt.Sprintf("please provide valid credentials for %s", v.SourceConfig.URL)
		return errors.New(errMsg)
	}
	err = v.Client.FindVM(v.SourceConfig.VmInventoryPath)
	if err != nil {
		errorMsg := "VM path is invalid\nPlease make sure to format your inventory path correctly using the 'vm' keyword. Example: /my-datacenter/vm/my-folder/my-vm-name"
		return errors.New(errorMsg)
	}
	return nil
}

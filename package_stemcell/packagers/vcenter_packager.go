package packagers

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/cloudfoundry-incubator/stembuild/filesystem"

	"github.com/cloudfoundry-incubator/stembuild/package_stemcell/config"
	"github.com/cloudfoundry-incubator/stembuild/package_stemcell/iaas_clients"
)

type VCenterPackager struct {
	SourceConfig config.SourceConfig
	OutputConfig config.OutputConfig
	Client       iaas_clients.IaasClient
}

func (v VCenterPackager) Package() error {
	err := v.Client.PrepareVM(v.SourceConfig.VmInventoryPath)
	if err != nil {
		return errors.New("could not prepare the VM for export")
	}
	workingDir, err := ioutil.TempDir(os.TempDir(), "vcenter-packager-working-directory")

	if err != nil {
		return errors.New("failed to create working directory")
	}

	stemcellDir, err := ioutil.TempDir(os.TempDir(), "vcenter-packager-stemcell-directory")
	if err != nil {
		return errors.New("failed to create stemcell directory")
	}

	err = v.Client.ExportVM(v.SourceConfig.VmInventoryPath, workingDir)
	if err != nil {
		return errors.New("failed to export the prepared VM")
	}

	shaSum, err := TarGenerator(filepath.Join(stemcellDir, "image"), workingDir)
	manifestContents := CreateManifest(v.OutputConfig.Os, v.OutputConfig.StemcellVersion, shaSum)
	err = WriteManifest(manifestContents, stemcellDir)
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

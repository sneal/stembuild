package packagers

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/cloudfoundry-incubator/stembuild/filesystem"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/cloudfoundry-incubator/stembuild/package_stemcell/iaas_clients/iaas_clientsfakes"

	"github.com/cloudfoundry-incubator/stembuild/package_stemcell/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("VcenterPackager", func() {

	var outputDir string
	var sourceConfig config.SourceConfig
	var outputConfig config.OutputConfig
	var fakeVcenterClient *iaas_clientsfakes.FakeIaasClient

	BeforeEach(func() {
		outputDir, _ = ioutil.TempDir(os.TempDir(), "vcenter-test-output-dir")
		sourceConfig = config.SourceConfig{Password: "password", URL: "url", Username: "username", VmInventoryPath: "path/valid-vm-name"}
		outputConfig = config.OutputConfig{Os: "2012R2", StemcellVersion: "1200.2", OutputDir: outputDir}
		fakeVcenterClient = &iaas_clientsfakes.FakeIaasClient{}
	})

	AfterEach(func() {
		_ = os.RemoveAll(outputDir)
	})

	Context("ValidateSourceParameters", func() {
		It("returns an error if the vCenter url is invalid", func() {

			fakeVcenterClient.ValidateUrlReturns(errors.New("invalid url"))

			packager := VCenterPackager{SourceConfig: sourceConfig, OutputConfig: outputConfig, Client: fakeVcenterClient}
			err := packager.ValidateSourceParameters()

			Expect(err).To(HaveOccurred())
			Expect(fakeVcenterClient.ValidateUrlCallCount()).To(Equal(1))
			Expect(err.Error()).To(Equal("please provide a valid vCenter URL"))

		})
		It("returns an error if the vCenter credentials are not valid", func() {

			fakeVcenterClient.ValidateCredentialsReturns(errors.New("invalid credentials"))

			packager := VCenterPackager{SourceConfig: sourceConfig, OutputConfig: outputConfig, Client: fakeVcenterClient}

			err := packager.ValidateSourceParameters()

			Expect(err).To(HaveOccurred())
			Expect(fakeVcenterClient.ValidateCredentialsCallCount()).To(Equal(1))
			Expect(err.Error()).To(ContainSubstring("please provide valid credentials for"))
		})

		It("returns an error if VM given does not exist ", func() {
			fakeVcenterClient.FindVMReturns(errors.New("invalid VM path"))

			packager := VCenterPackager{SourceConfig: sourceConfig, OutputConfig: outputConfig, Client: fakeVcenterClient}

			err := packager.ValidateSourceParameters()

			Expect(err).To(HaveOccurred())
			Expect(fakeVcenterClient.FindVMCallCount()).To(Equal(1))
			Expect(err.Error()).To(Equal("VM path is invalid\nPlease make sure to format your inventory path correctly using the 'vm' keyword. Example: /my-datacenter/vm/my-folder/my-vm-name"))
		})
		It("returns no error if all source parameters are valid", func() {

			packager := VCenterPackager{SourceConfig: sourceConfig, OutputConfig: outputConfig, Client: fakeVcenterClient}

			err := packager.ValidateSourceParameters()

			Expect(err).NotTo(HaveOccurred())
		})
	})
	Context("ValidateFreeSpace", func() {
		It("Print a warning to make sure that enough space is available", func() {
			packager := VCenterPackager{SourceConfig: sourceConfig, OutputConfig: outputConfig, Client: fakeVcenterClient}
			err := packager.ValidateFreeSpaceForPackage(&filesystem.OSFileSystem{})

			Expect(err).To(Not(HaveOccurred()))
		})
	})
	Context("Package failure cases", func() {

		It("Package fails if devices were removed unsuccessfully", func() {
			packager := VCenterPackager{SourceConfig: sourceConfig, OutputConfig: outputConfig, Client: fakeVcenterClient}
			fakeVcenterClient.PrepareVMReturns(errors.New("floppy-8000 could not be removed/not found"))
			err := packager.Package()

			Expect(err).To(HaveOccurred())
			Expect(fakeVcenterClient.PrepareVMCallCount()).To(Equal(1))
			args := fakeVcenterClient.PrepareVMArgsForCall(0)
			Expect(args).To(Equal(sourceConfig.VmInventoryPath))
			Expect(err.Error()).To(Equal("could not prepare the VM for export"))
		})

		It("Returns a error message if exporting the VM fails", func() {
			packager := VCenterPackager{SourceConfig: sourceConfig, OutputConfig: outputConfig, Client: fakeVcenterClient}
			fakeVcenterClient.PrepareVMReturns(nil)
			fakeVcenterClient.ExportVMReturns(errors.New(fmt.Sprintf(sourceConfig.VmInventoryPath + " could not be exported")))
			err := packager.Package()

			Expect(err).To(HaveOccurred())
			Expect(fakeVcenterClient.ExportVMCallCount()).To(Equal(1))
			vmPath, _ := fakeVcenterClient.ExportVMArgsForCall(0)
			Expect(vmPath).To(Equal(sourceConfig.VmInventoryPath))
			Expect(err.Error()).To(Equal("failed to export the prepared VM"))
		})
	})

	Context("Package successful case", func() {
		AfterEach(func() {
			_ = os.RemoveAll("./valid-vm-name")
			_ = os.RemoveAll("image")
		})

		It("creates a valid stemcell in the output directory", func() {
			packager := VCenterPackager{SourceConfig: sourceConfig, OutputConfig: outputConfig, Client: fakeVcenterClient}
			fakeVcenterClient.PrepareVMReturns(nil)
			fileContentMap := make(map[string]string)
			fileContentMap["stemcell.MF"] = "file1 content\n"
			fileContentMap["image"] = "file2 content\n"

			fakeVcenterClient.ExportVMStub = func(vmInventoryPath string, destination string) error {
				vmName := path.Base(vmInventoryPath)
				os.Mkdir(filepath.Join(destination, vmName), 0777)

				testOvfName := "valid-vm-name.content"
				err := ioutil.WriteFile(filepath.Join(filepath.Join(destination, vmName), testOvfName), []byte(""), 0777)
				if err != nil {
					log.Fatal(err)
				}
				return nil
			}

			err := packager.Package()

			Expect(err).To(Not(HaveOccurred()))
			stemcellFilename := StemcellFilename(packager.OutputConfig.StemcellVersion, packager.OutputConfig.Os)
			stemcellFile := filepath.Join(packager.OutputConfig.OutputDir, stemcellFilename)
			_, err = os.Stat(stemcellFile)

			Expect(err).NotTo(HaveOccurred())
			var actualStemcellManifestContent string
			expectedManifestContent := `---
name: bosh-vsphere-esxi-windows2012R2-go_agent
version: '1200.2'
sha1: %x
operating_system: windows2012R2
cloud_properties:
  infrastructure: vsphere
  hypervisor: esxi
stemcell_formats:
- vsphere-ovf
- vsphere-ova
`
			var fileReader, _ = os.OpenFile(stemcellFile, os.O_RDONLY, 0777)
			gzr, err := gzip.NewReader(fileReader)
			Expect(err).ToNot(HaveOccurred())
			defer gzr.Close()
			tarfileReader := tar.NewReader(gzr)
			count := 0

			for {
				header, err := tarfileReader.Next()
				if err == io.EOF {
					break
				}

				Expect(err).NotTo(HaveOccurred())

				switch filepath.Base(header.Name) {
				case "stemcell.MF":
					buf := new(bytes.Buffer)
					_, err = buf.ReadFrom(tarfileReader)
					Expect(err).NotTo(HaveOccurred())
					count++

					actualStemcellManifestContent = buf.String()

				case "image":
					count++
					actualSha1 := sha1.New()
					io.Copy(actualSha1, tarfileReader)

					expectedManifestContent = fmt.Sprintf(expectedManifestContent, actualSha1.Sum(nil))

				default:

					Fail(fmt.Sprintf("Found unknown file: %s", filepath.Base(header.Name)))
				}
			}
			Expect(count).To(Equal(2))
			Expect(actualStemcellManifestContent).To(Equal(expectedManifestContent))
		})
	})
})

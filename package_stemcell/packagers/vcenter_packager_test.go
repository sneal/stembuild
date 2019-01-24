package packagers

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/cloudfoundry-incubator/stembuild/filesystem"

	"github.com/cloudfoundry-incubator/stembuild/package_stemcell/iaas_clients/iaas_clientsfakes"

	"github.com/cloudfoundry-incubator/stembuild/package_stemcell/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("VcenterPackager", func() {

	var sourceConfig config.SourceConfig
	var fakeVcenterClient *iaas_clientsfakes.FakeIaasClient

	BeforeEach(func() {
		sourceConfig = config.SourceConfig{Password: "password", URL: "url", Username: "username", VmInventoryPath: "path/valid-vm-name"}
		fakeVcenterClient = &iaas_clientsfakes.FakeIaasClient{}
	})

	Context("ValidateSourceParameters", func() {
		It("returns an error if the vCenter url is invalid", func() {

			fakeVcenterClient.ValidateUrlReturns(errors.New("invalid url"))

			packager := VCenterPackager{SourceConfig: sourceConfig, Client: fakeVcenterClient}
			err := packager.ValidateSourceParameters()

			Expect(err).To(HaveOccurred())
			Expect(fakeVcenterClient.ValidateUrlCallCount()).To(Equal(1))
			Expect(err.Error()).To(Equal("please provide a valid vCenter URL"))

		})
		It("returns an error if the vCenter credentials are not valid", func() {

			fakeVcenterClient.ValidateCredentialsReturns(errors.New("invalid credentials"))

			packager := VCenterPackager{SourceConfig: sourceConfig, Client: fakeVcenterClient}

			err := packager.ValidateSourceParameters()

			Expect(err).To(HaveOccurred())
			Expect(fakeVcenterClient.ValidateCredentialsCallCount()).To(Equal(1))
			Expect(err.Error()).To(ContainSubstring("please provide valid credentials for"))
		})

		It("returns an error if VM given does not exist ", func() {
			fakeVcenterClient.FindVMReturns(errors.New("invalid VM path"))

			packager := VCenterPackager{SourceConfig: sourceConfig, Client: fakeVcenterClient}

			err := packager.ValidateSourceParameters()

			Expect(err).To(HaveOccurred())
			Expect(fakeVcenterClient.FindVMCallCount()).To(Equal(1))
			Expect(err.Error()).To(Equal("VM path is invalid\nPlease make sure to format your inventory path correctly using the 'vm' keyword. Example: /my-datacenter/vm/my-folder/my-vm-name"))
		})
		It("returns no error if all source parameters are valid", func() {

			packager := VCenterPackager{SourceConfig: sourceConfig, Client: fakeVcenterClient}

			err := packager.ValidateSourceParameters()

			Expect(err).NotTo(HaveOccurred())
		})
	})
	Context("ValidateFreeSpace", func() {
		It("Print a warning to make sure that enough space is available", func() {
			packager := VCenterPackager{SourceConfig: sourceConfig, Client: fakeVcenterClient}
			err := packager.ValidateFreeSpaceForPackage(&filesystem.OSFileSystem{})

			Expect(err).To(Not(HaveOccurred()))
		})
	})
	Context("Package failure cases", func() {

		It("Package fails if devices were removed unsuccessfully", func() {
			packager := VCenterPackager{SourceConfig: sourceConfig, Client: fakeVcenterClient}
			fakeVcenterClient.PrepareVMReturns(errors.New("floppy-8000 could not be removed/not found"))
			err := packager.Package()

			Expect(err).To(HaveOccurred())
			Expect(fakeVcenterClient.PrepareVMCallCount()).To(Equal(1))
			args := fakeVcenterClient.PrepareVMArgsForCall(0)
			Expect(args).To(Equal(sourceConfig.VmInventoryPath))
			Expect(err.Error()).To(Equal("could not prepare the VM for export"))
		})

		It("Returns a error message if exporting the VM fails", func() {
			packager := VCenterPackager{SourceConfig: sourceConfig, Client: fakeVcenterClient}
			fakeVcenterClient.PrepareVMReturns(nil)
			fakeVcenterClient.ExportVMReturns(errors.New(fmt.Sprintf(sourceConfig.VmInventoryPath + " could not be exported")))
			err := packager.Package()

			Expect(err).To(HaveOccurred())
			Expect(fakeVcenterClient.ExportVMCallCount()).To(Equal(1))
			args := fakeVcenterClient.ExportVMArgsForCall(0)
			Expect(args).To(Equal(sourceConfig.VmInventoryPath))
			Expect(err.Error()).To(Equal("failed to export the prepared VM"))
		})
	})

	Context("Package successful case", func() {
		AfterEach(func() {
			_ = os.RemoveAll("./valid-vm-name")
			_ = os.RemoveAll("image")
		})

		It("exported VM is a tar named image", func() {
			packager := VCenterPackager{SourceConfig: sourceConfig, Client: fakeVcenterClient}
			fakeVcenterClient.PrepareVMReturns(nil)
			splitVmInventoryPath := strings.Split(sourceConfig.VmInventoryPath, "/")
			vmName := splitVmInventoryPath[len(splitVmInventoryPath)-1]

			fakeVcenterClient.ExportVMStub = func(vmInventoryPath string) error {

				err := os.Mkdir(vmName, 0777)
				if err != nil {
					log.Fatal(err)
				}

				testOvfName := "valid-vm-name.content"
				err = ioutil.WriteFile(fmt.Sprintf(vmName+"/"+testOvfName), []byte(""), 0777)
				if err != nil {
					log.Fatal(err)
				}

				return nil
			}

			err := packager.Package()

			Expect(err).To(Not(HaveOccurred()))
			_, err = os.Stat(fmt.Sprintf("image"))

		})
	})
})

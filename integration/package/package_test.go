package package_test

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha1"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/onsi/gomega/gexec"

	"github.com/cloudfoundry-incubator/stembuild/test/helpers"

	"github.com/vmware/govmomi/govc/cli"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	_ "github.com/vmware/govmomi/govc/vm"
)

var _ = Describe("Package", func() {
	const (
		baseVMName              = "stembuild-package-integration-tests-base-vm"
		stemcellVersion         = "1803.5.3999-manual.1"
		osVersion               = "1803"
		vcenterURLVariable      = "GOVC_URL"
		vcenterUsernameVariable = "GOVC_USERNAME"
		vcenterPasswordVariable = "GOVC_PASSWORD"
		vcenterFolderVariable   = "VM_FOLDER"
		existingVMVariable      = "EXISTING_SOURCE_VM"
	)

	var (
		workingDir   string
		sourceVMName string
		vmPath       string
	)

	It("generates a stemcell with the correct shasum", func() {
		executable, err := helpers.BuildStembuild()
		Expect(err).NotTo(HaveOccurred())

		existingVM := os.Getenv(existingVMVariable)
		vcenterFolder := helpers.EnvMustExist(vcenterFolderVariable)

		rand.Seed(time.Now().UnixNano())
		if existingVM == "" {
			sourceVMName = fmt.Sprintf("stembuild-package-test-%d", rand.Int())
		} else {
			sourceVMName = fmt.Sprintf("%s-%d", existingVM, rand.Int())
		}

		baseVMWithPath := fmt.Sprintf(vcenterFolder + "/" + baseVMName)
		vmPath = strings.Join([]string{vcenterFolder, sourceVMName}, "/")

		cli.Run([]string{
			"vm.clone",
			"-vm", baseVMWithPath,
			"-on=false",
			sourceVMName,
		})

		vcenterURL := helpers.EnvMustExist(vcenterURLVariable)
		vcenterUsername := helpers.EnvMustExist(vcenterUsernameVariable)
		vcenterPassword := helpers.EnvMustExist(vcenterPasswordVariable)

		workingDir, err = ioutil.TempDir(os.TempDir(), "stembuild-package-test")
		if err != nil {
			log.Fatal(err)
		}

		Expect(err).NotTo(HaveOccurred())

		session := helpers.RunCommandInDir(
			workingDir,
			executable, "package",
			"-url", vcenterURL,
			"-username", vcenterUsername,
			"-password", vcenterPassword,
			"-vm-inventory-path", vmPath,
			"-stemcell-version", stemcellVersion,
			"-os", osVersion,
		)

		Eventually(session, 60*time.Minute, 5*time.Second).Should(gexec.Exit(0))
		var out []byte
		session.Out.Write(out)
		fmt.Print(string(out))

		expectedFilename := fmt.Sprintf(
			"bosh-stemcell-%s-vsphere-esxi-windows%s-go_agent.tgz",
			stemcellVersion, osVersion,
		)
		stemcellPath := filepath.Join(workingDir, expectedFilename)

		stemcellZip, err := os.OpenFile(stemcellPath, os.O_RDONLY, 0777)
		Expect(err).NotTo(HaveOccurred())

		gzr, err := gzip.NewReader(stemcellZip)
		Expect(err).ToNot(HaveOccurred())
		defer gzr.Close()

		var imageZip *os.File
		r := tar.NewReader(gzr)
		for {
			header, err := r.Next()
			if err == io.EOF {
				break
			}

			if header.Name == "image" {
				image, err := os.Create(filepath.Join(workingDir, "image"))
				Expect(err).NotTo(HaveOccurred())
				_, err = io.Copy(image, r)
				Expect(err).NotTo(HaveOccurred())
			}
		}

		imageZip, err = os.OpenFile(filepath.Join(workingDir, "image"), os.O_RDONLY, 0777)
		Expect(err).NotTo(HaveOccurred())

		gzr, err = gzip.NewReader(imageZip)
		Expect(err).ToNot(HaveOccurred())
		defer gzr.Close()

		r = tar.NewReader(gzr)
		vmdkSha := sha1.New()
		ovfSha := sha1.New()
		var imageMF *os.File
		for {
			header, err := r.Next()
			if err == io.EOF {
				break
			}

			if strings.Contains(header.Name, ".vmdk") {
				io.Copy(vmdkSha, r)
			}
			if strings.Contains(header.Name, ".ovf") {
				io.Copy(ovfSha, r)
			}
			if strings.Contains(header.Name, ".mf") {
				imageMF, err = os.Create(filepath.Join(workingDir, "image.mf"))
				Expect(err).NotTo(HaveOccurred())
				_, err = io.Copy(imageMF, r)
				Expect(err).NotTo(HaveOccurred())
			}
		}

		imageMFFile, err := helpers.ReadFile(filepath.Join(workingDir, "image.mf"))
		Expect(err).NotTo(HaveOccurred())

		Expect(imageMFFile).To(ContainSubstring(fmt.Sprintf("%x", vmdkSha.Sum(nil))))
		Expect(imageMFFile).To(ContainSubstring(fmt.Sprintf("%x", ovfSha.Sum(nil))))

	})

	AfterEach(func() {
		//os.RemoveAll(workingDir)
		// fix clone vm path (not in winnipeg)
		if vmPath != "" {
			cli.Run([]string{"vm.destroy", "-vm.ipath", vmPath})
		}
	})
})

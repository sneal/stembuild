package package_test

import (
	"archive/tar"
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
		baseVMName      = "stembuild-package-test-source-base"
		stemcellVersion = "1803.5.3999-manual.1"
		osVersion       = "1803"

		expectedStemcellShaSum = "b811888f55e81eb3bee2c68c4ed2986f6b2825ba"

		vcenterURLVariable      = "GOVC_URL"
		vcenterUsernameVariable = "GOVC_USERNAME"
		vcenterPasswordVariable = "GOVC_PASSWORD"
		vcenterFolderVariable   = "GOVC_FOLDER"
		existingVMVariable      = "EXISTING_SOURCE_VM"
	)

	var (
		workingDir   string
		sourceVMName string
		vmPath       string
	)

	FIt("generates a stemcell with the correct shasum", func() {
		existingVM := os.Getenv(existingVMVariable)
		vcenterFolder := helpers.EnvMustExist(vcenterFolderVariable)

		rand.Seed(time.Now().UnixNano())
		if existingVM == "" {
			sourceVMName = "stembuild-package-test-1445168366415952161"//fmt.Sprintf("stembuild-package-test-1445168366415952161", rand.Int())
		} else {
			sourceVMName = fmt.Sprintf("%s-%d", existingVM, rand.Int())
		}
		println(sourceVMName)
		baseVMWithPath := fmt.Sprintf(vcenterFolder + "/" + baseVMName)
		vmPath = strings.Join([]string{vcenterFolder, sourceVMName}, "/")
		cli.Run([]string{"vm.clone", "-vm", baseVMWithPath, "-on=false", sourceVMName})

		executable, err := helpers.BuildStembuild()
		Expect(err).NotTo(HaveOccurred())

		vcenterURL := helpers.EnvMustExist(vcenterURLVariable)
		vcenterUsername := helpers.EnvMustExist(vcenterUsernameVariable)
		vcenterPassword := helpers.EnvMustExist(vcenterPasswordVariable)

		workingDir, err = ioutil.TempDir(os.TempDir(), "stembuild-package-test")
		fmt.Printf("test working directory: %s\n", workingDir)
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(dir)
		fmt.Println("--------------------------------------------------------------")

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

		fmt.Print(session.Out)
		Eventually(session, 30*time.Minute, 5*time.Second).Should(gexec.Exit(0))

		expectedFilename := fmt.Sprintf(
			"bosh-stemcell-%s-vsphere-esxi-windows%s-go_agent.tgz",
			stemcellVersion, osVersion,
		)
		fp := filepath.Join(workingDir, expectedFilename)

		f, err := os.OpenFile(fp, os.O_RDONLY, 0777)
		Expect(err).NotTo(HaveOccurred())

		r := tar.NewReader(f)
		for {
			header, err := r.Next()
			if err == io.EOF {
				break
			}

			fmt.Printf("Header name: %s\n", header.Name)
			if header.Name == "image" {
				f, err = os.Create(filepath.Join(workingDir, "image"))
				Expect(err).NotTo(HaveOccurred())
				_, err = io.Copy(f, r)
				Expect(err).NotTo(HaveOccurred())
			}
		}

		f, err = os.OpenFile(filepath.Join(workingDir, "image"), os.O_RDONLY, 0777)
		Expect(err).NotTo(HaveOccurred())
		r = tar.NewReader(f)
		h := sha1.New()
		for {
			header, err := r.Next()
			if err == io.EOF {
				break
			}

			fmt.Printf("Header name: %s\n", header.Name)
			if strings.Contains(header.Name, ".vmdk") {
				io.Copy(h, r)
			}
		}

		Expect(fmt.Sprintf("%x", h.Sum(nil))).To(Equal(expectedStemcellShaSum))
	})

	AfterEach(func() {
		os.RemoveAll(workingDir)
		//fix clone vm path (not in winnipeg)
		if vmPath != "" {
			cli.Run([]string{"vm.destroy", "-vm.ipath", vmPath})
		}
	})
})

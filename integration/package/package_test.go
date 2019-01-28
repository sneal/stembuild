package package_test

import (
	"crypto/sha1"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
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

		expectedStemcellShaSum = "2c2fb8230e268f2decbd4115eae67e5636afbae6"

		vcenterURLVariable      = "GOVC_URL"
		vcenterUsernameVariable = "GOVC_USERNAME"
		vcenterPasswordVariable = "GOVC_PASSWORD"
		vcenterFolderVariable   = "GOVC_FOLDER"
		existingVMVariable      = "EXISTING_SOURCE_VM"
	)

	var (
		workingDir   string
		sourceVMName string
	)

	It("generates a stemcell with the correct shasum", func() {
		sourceVMName = os.Getenv(existingVMVariable)
		vcenterFolder := helpers.EnvMustExist(vcenterFolderVariable)

		if sourceVMName == "" {
			rand.Seed(time.Now().UnixNano())
			sourceVMName = fmt.Sprintf("stembuild-package-test-%d", rand.Int())
			println(sourceVMName)
			baseVMWithPath := fmt.Sprintf(vcenterFolder + "/" + baseVMName)
			cli.Run([]string{"vm.clone", "-vm", baseVMWithPath, "-on=false", sourceVMName})
		}

		executable, err := helpers.BuildStembuild()
		Expect(err).NotTo(HaveOccurred())

		vcenterURL := helpers.EnvMustExist(vcenterURLVariable)
		vcenterUsername := helpers.EnvMustExist(vcenterUsernameVariable)
		vcenterPassword := helpers.EnvMustExist(vcenterPasswordVariable)

		vmPath := fmt.Sprintf("%s/%s", vcenterFolder, sourceVMName)

		workingDir, err = ioutil.TempDir(os.TempDir(), "stembuild-package-test")
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
		expectedFilePath := filepath.Join(workingDir, expectedFilename)
		stemcell, err := os.Open(expectedFilePath)

		Expect(err).NotTo(HaveOccurred())

		defer stemcell.Close()

		h := sha1.New()
		_, err = io.Copy(h, stemcell)
		Expect(err).NotTo(HaveOccurred())

		Expect(h.Sum(nil)).To(Equal(expectedStemcellShaSum))
		//Assert stemcell content?
		//Assert image tar content?
		//delete test vm in after suite
	})
})

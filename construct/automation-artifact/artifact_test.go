package artifact_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"os"

	"github.com/pivotal-cf-experimental/stembuild/construct/automation-artifact"
	//"github.com/pivotal-cf-experimental/stembuild/pack/options"
	//"github.com/pivotal-cf-experimental/stembuild/pack/stemcell"
)

var _ = Describe("Artifact", func() {
	var tmpDir string
	//var stembuildConfig options.StembuildOptions
	//var c stemcell.Config

	BeforeEach(func() {
		var err error
		tmpDir, err = ioutil.TempDir("", "")
		Expect(err).NotTo(HaveOccurred())

		//stembuildConfig = options.StembuildOptions{
		//	OSVersion: "2012R2",
		//	Version:   "1200.1",
		//}

		//c = stemcell.Config{
		//	Stop:         make(chan struct{}),
		//	Debugf:       func(format string, a ...interface{}) {},
		//	BuildOptions: stembuildConfig,
		//}
	})

	AfterEach(func() {
		Expect(os.RemoveAll(tmpDir)).To(Succeed())
	})

	Describe("Pivnet", func() {

		FContext("when a valid pivnet UAA refresh token is provided", func() {
			It("downloads stemcell automation release at given path ", func() {

				_, err := os.Stat("StemcellAutomation.zip")
				Expect(os.IsNotExist(err)).To(BeTrue())

				err = artifact.DownloadAutomationArtifact("323c5fac66c34dcd85628138d736c8b1-r", tmpDir)

				Expect(err).ToNot(HaveOccurred())
				_, err = os.Stat(tmpDir + "StemcellAutomation.zip")

				Expect(os.IsNotExist(err)).To(BeFalse())

				Expect(err).ToNot(HaveOccurred())

			})
		})
	})
})

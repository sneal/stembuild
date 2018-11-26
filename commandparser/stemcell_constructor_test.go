package commandparser

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"os"
	"path"
)

var _ = Describe("StemcellConstructor", func() {
	Describe("UploadArtifacts", func() {
		Context("When the inputs are accurate", func() {

			FIt("uploads file to the VM", func() {
				manager := NewFakeRemoteManager()
				constructor := &StemcellConstructor{manager}
				wd, _ := os.Getwd()
				file, _ := ioutil.TempFile(wd, "LGPO.exe")
				defer os.Remove(file.Name())
				filepath := path.Join(wd, file.Name())

				constructor.uploadArtifact(filepath)

				Expect(manager.FilePath).To(Equal(filepath))
			})
		})
	})
})

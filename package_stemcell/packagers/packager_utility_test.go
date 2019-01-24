package packagers

import (
	"archive/tar"
	"bytes"
	"crypto/sha1"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Packager Utility", func() {
	Context("TarGenerator", func() {
		var sourceDir string
		var destinationDir string

		BeforeEach(func() {
			sourceDir, _ = ioutil.TempDir(os.TempDir(), "packager-utility-test-source")
			destinationDir, _ = ioutil.TempDir(os.TempDir(), "packager-utility-test-destination")
		})

		It("should tar all files inside provided folder and return its sha1", func() {
			err := ioutil.WriteFile(filepath.Join(sourceDir, "file1"), []byte("file1 content\n"), 0777)
			Expect(err).NotTo(HaveOccurred())
			err = ioutil.WriteFile(filepath.Join(sourceDir, "file2"), []byte("file2 content\n"), 0777)
			Expect(err).NotTo(HaveOccurred())
			fileContentMap := make(map[string]string)
			fileContentMap["file1"] = "file1 content\n"
			fileContentMap["file2"] = "file2 content\n"

			tarball := filepath.Join(destinationDir, "tarball")

			sha1Sum, err := TarGenerator(tarball, sourceDir)

			expectedSha1 := sha1.New()

			Expect(err).NotTo(HaveOccurred())

			_, err = os.Stat(tarball)
			Expect(err).NotTo(HaveOccurred())
			var fileReader, _ = os.OpenFile(tarball, os.O_RDONLY, 0777)

			tarfileReader := tar.NewReader(fileReader)
			count := 0
			for {
				header, err := tarfileReader.Next()
				if err == io.EOF {
					break
				}
				count++
				Expect(err).NotTo(HaveOccurred())
				buf := new(bytes.Buffer)
				_, err = buf.ReadFrom(tarfileReader)
				if err != nil {
					break
				}
				Expect(fileContentMap[filepath.Base(header.Name)]).To(Equal(buf.String()))
			}
			Expect(count).To(Equal(2))

			tarballFile, err := os.Open(tarball)
			defer tarballFile.Close()

			io.Copy(expectedSha1, tarballFile)

			sum := fmt.Sprintf("%x", expectedSha1.Sum(nil))
			Expect(sha1Sum).To(Equal(sum))
		})

		AfterEach(func() {
			os.RemoveAll(sourceDir)
			os.RemoveAll(destinationDir)
		})
	})

	Context("CreateManifest", func() {
		It("Creates a manifest correctly", func() {
			expectedManifest := `---
name: bosh-vsphere-esxi-windows1-go_agent
version: 'version'
sha1: sha1sum
operating_system: windows1
cloud_properties:
  infrastructure: vsphere
  hypervisor: esxi
stemcell_formats:
- vsphere-ovf
- vsphere-ova
`
			result := CreateManifest("1", "version", "sha1sum")
			Expect(result).To(Equal(expectedManifest))
		})
	})
})

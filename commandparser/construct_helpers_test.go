package commandparser_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotal-cf-experimental/stembuild/commandparser"
)

var _ = Describe("Check if necessary files are in directory", func() {
	Context("Directory given is valid", func() {
		Describe("automation artifact", func() {
			filename := "StemcellAutomation.zip"
			Context("artifact is not present", func() {
				dir := "/Users/pivotal/go/src/github.com/pivotal-cf-experimental/stembuild/test/constructData/emptyDir"
				It("should return false with no error", func() {
					present, err := IsArtifactInDirectory(dir, filename)
					Expect(err).ToNot(HaveOccurred())
					Expect(present).To(BeFalse())
				})
			})
			Context("artifact is present", func() {
				dir := "/Users/pivotal/go/src/github.com/pivotal-cf-experimental/stembuild/test/constructData/fullDir"
				It("should return true with no error", func() {
					present, err := IsArtifactInDirectory(dir, filename)
					Expect(err).ToNot(HaveOccurred())
					Expect(present).To(BeTrue())
				})
			})
		})
		Describe("LGPO", func() {
			filename := "LGPO.zip"
			Context("LGPO is not present", func() {
				dir := "/Users/pivotal/go/src/github.com/pivotal-cf-experimental/stembuild/test/constructData/emptyDir"
				It("should return false with no error", func() {
					present, err := IsArtifactInDirectory(dir, filename)
					Expect(err).ToNot(HaveOccurred())
					Expect(present).To(BeFalse())
				})
			})
			Context("artifact is present", func() {
				dir := "/Users/pivotal/go/src/github.com/pivotal-cf-experimental/stembuild/test/constructData/fullDir"
				It("should return true with no error", func() {
					present, err := IsArtifactInDirectory(dir, filename)
					Expect(err).ToNot(HaveOccurred())
					Expect(present).To(BeTrue())
				})
			})
		})
	})
	Context("Directory given is not valid", func() {
		filename := "file"
		It("should return an error", func() {
			dir := "/Users/pivotal/go/src/github.com/pivotal-cf-experimental/stembuild/test/constructData/notExist"
			_, err := IsArtifactInDirectory(dir, filename)
			Expect(err).To(HaveOccurred())
		})
	})
})
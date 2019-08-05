package construct_test

import (
	"errors"

	"github.com/cloudfoundry-incubator/stembuild/construct"
	"github.com/cloudfoundry-incubator/stembuild/construct/constructfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("OsVersionValidator", func() {
	var (
		validator        *construct.OSVersionValidator
		fakeGuestManager *constructfakes.FakeGuestManager
		fakeMessenger    *constructfakes.FakeOSValidatorMessenger
	)

	BeforeEach(func() {
		fakeGuestManager = &constructfakes.FakeGuestManager{}
		fakeMessenger = &constructfakes.FakeOSValidatorMessenger{}

		validator = &construct.OSVersionValidator{
			GuestManager: fakeGuestManager,
			Messenger:    fakeMessenger,
		}

		versionBuffer := gbytes.NewBuffer()
		_, err := versionBuffer.Write([]byte("dev"))
		Expect(err).NotTo(HaveOccurred())

		fakeGuestManager.DownloadFileInGuestReturns(versionBuffer, 3, nil)
	})

	Describe("Validate", func() {
		It("returns nil when the version of stembuild matches the version of the Guest OS VM", func() {
			expectedPid := int64(123)
			fakeGuestManager.StartProgramInGuestReturns(expectedPid, nil)

			err := validator.Validate()

			Expect(err).NotTo(HaveOccurred())

			Expect(fakeGuestManager.StartProgramInGuestCallCount()).To(Equal(1))
			_, _, args := fakeGuestManager.StartProgramInGuestArgsForCall(0)
			Expect(args).To(ContainSubstring("OSVersion.Version.Build"))
			Expect(args).To(ContainSubstring("\\Windows\\Temp"))

			Expect(fakeGuestManager.ExitCodeForProgramInGuestCallCount()).To(Equal(1))
			_, pid := fakeGuestManager.ExitCodeForProgramInGuestArgsForCall(0)
			Expect(pid).To(Equal(expectedPid))
		})

		It("returns an error when the VM's version does not match stembuild's version", func() {
			versionBuffer := gbytes.NewBuffer()
			_, err := versionBuffer.Write([]byte("BAD VERSION"))
			Expect(err).NotTo(HaveOccurred())

			fakeGuestManager.DownloadFileInGuestReturns(versionBuffer, 11, nil)

			err = validator.Validate()

			Expect(err).To(HaveOccurred())
		})

		It("returns nil even if get OS version file creation fails", func() {
			fakePid := 123
			fakeGuestManager.StartProgramInGuestReturnsOnCall(0, int64(fakePid), errors.New("failed to create blah"))

			err := validator.Validate()

			Expect(err).NotTo(HaveOccurred())
			Expect(fakeGuestManager.DownloadFileInGuestCallCount()).To(Equal(0))
			Expect(fakeMessenger.OSVersionFileCreationFailedCallCount()).To(Equal(1))
		})

		It("returns nil if the exit code for OS version file creation process cannot be retrieved", func() {
			fakeExitCode := 123
			fakeGuestManager.ExitCodeForProgramInGuestReturnsOnCall(0, int32(fakeExitCode), errors.New("failed to get exit code for process"))

			err := validator.Validate()

			Expect(err).NotTo(HaveOccurred())
			Expect(fakeGuestManager.DownloadFileInGuestCallCount()).To(Equal(0))
			Expect(fakeMessenger.ExitCodeRetrievalFailedCallCount()).To(Equal(1))
		})

		It("returns nil if the exit code for OS version file creation process is non-zero", func() {
			fakeExitCode := 123
			fakeGuestManager.StartProgramInGuestReturns(123, nil)
			fakeGuestManager.ExitCodeForProgramInGuestReturnsOnCall(0, int32(fakeExitCode), nil)

			err := validator.Validate()

			Expect(err).NotTo(HaveOccurred())
			Expect(fakeGuestManager.DownloadFileInGuestCallCount()).To(Equal(0))
			Expect(fakeMessenger.OSVersionFileCreationFailedCallCount()).To(Equal(1))
		})

		It("returns nil if the os version file that was created cannot be downloaded", func() {
			fakeGuestManager.DownloadFileInGuestReturnsOnCall(0, nil, 0, errors.New("could not download"))

			err := validator.Validate()

			Expect(err).NotTo(HaveOccurred())
			Expect(fakeMessenger.DownloadFileFailedCallCount()).To(Equal(1))
		})
	})

})

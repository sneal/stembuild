package construct_test

import (
	"errors"

	"github.com/cloudfoundry-incubator/stembuild/construct"
	"github.com/cloudfoundry-incubator/stembuild/construct/constructfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("VmAuthenticationValidator", func() {
	var (
		validator        *construct.VMAuthValidator
		fakeGuestManager *constructfakes.FakeGuestManager
		//fakeMessenger *constructfakes.mess
	)

	BeforeEach(func() {
		fakeGuestManager = &constructfakes.FakeGuestManager{}
		validator = &construct.VMAuthValidator{
			GuestManager: fakeGuestManager,
		}
	})

	Describe("IsValidAuth", func() {
		It("should return False when the username / password is incorrect", func() {
			fakeGuestManager.StartProgramInGuestReturns(1, errors.New("test error: Failed to authenticate with the guest operating system"))

			isValidAuth, err := validator.IsValidAuth()

			Expect(err).To(Not(HaveOccurred()))
			Expect(isValidAuth).To(BeFalse())
		})

		It("should return True when the username / password is correct", func() {
			fakeGuestManager.StartProgramInGuestReturns(0, nil)

			isValidAuth, err := validator.IsValidAuth()

			Expect(err).To(Not(HaveOccurred()))
			Expect(isValidAuth).To(BeTrue())
		})

		It("should return an error when non-auth related errors occur", func() {
			fakeGuestManager.StartProgramInGuestReturns(1, errors.New("DIFFERENT ERROR"))

			isValidAuth, err := validator.IsValidAuth()

			Expect(err).To(HaveOccurred())
			Expect(isValidAuth).To(BeFalse())
		})
	})

})

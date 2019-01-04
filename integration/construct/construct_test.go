package construct_test

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/cloudfoundry-incubator/stembuild/test/helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("stembuild construct", func() {
	It("transfers LGPO and StemcellAutomation archives, unarchive them and execute automation script", func() {
		pwd, err := os.Getwd()
		if err != nil {
			fmt.Fprint(os.Stderr, "unable to find current working directory", err)
		}
		CopyFile(fmt.Sprintf(pwd+"/assets/StemcellAutomation.zip"), fmt.Sprintf(pwd+"/StemcellAutomation.zip"))
		CopyFile(fmt.Sprintf(pwd+"/assets/LGPO.zip"), fmt.Sprintf(pwd+"/LGPO.zip"))

		session := helpers.Stembuild(stembuildExecutable, "construct", "-winrm-ip", conf.TargetIP, "-stemcell-version", "1709.1", "-winrm-username", conf.VMUsername, "-winrm-password", conf.VMPassword)

		Eventually(session, 20).Should(Exit(0))
		Eventually(session.Out).Should(Say(`mock stemcell automation script executed`))
	})

	It("fails with an appropriate error when LGPO and/or StemcellAutomation is missing", func() {
		session := helpers.Stembuild(stembuildExecutable, "construct", "-winrm-ip", conf.TargetIP, "-stemcell-version", "1803.1", "-winrm-username", conf.VMUsername, "-winrm-password", conf.VMPassword)

		Eventually(session, 20).Should(Exit(1))
		Eventually(session.Err).Should(Say(`automation artifact not found in current directory`))
	})

	AfterEach(func() {
		pwd, err := os.Getwd()
		if err != nil {
			fmt.Fprint(os.Stderr, "unable to find current working directory", err)
		}

		os.Remove(fmt.Sprintf(pwd + "/StemcellAutomation.zip"))
		os.Remove(fmt.Sprintf(pwd + "/LGPO.zip"))
	})
})

func CopyFile(src string, dest string) error {
	input, err := ioutil.ReadFile(src)
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = ioutil.WriteFile(dest, input, 0644)
	if err != nil {
		fmt.Println("Error creating file")
		fmt.Println(err)
		return err
	}

	return err
}

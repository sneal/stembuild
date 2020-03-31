package construct_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/session"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/methods"
	"github.com/vmware/govmomi/vim25/soap"
	"github.com/vmware/govmomi/vim25/types"

	"github.com/cloudfoundry-incubator/stembuild/test/helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

const(
	vCenterUsername = "USER"
	vCenterPassword = "PASS"
)

func buildSoapClient() (*soap.Client, error) {
	vCenterServer := "https://vcenter.wild.cf-app.com"
	username := vCenterUsername
	password := vCenterPassword
	//rootCACertPath :=

	vCenterURL, err := soap.ParseURL(vCenterServer)
	if err != nil {
		return nil, err
	}
	credentials := url.UserPassword(username, password)
	vCenterURL.User = credentials

	soapClient := soap.NewClient(vCenterURL, false)

	//if rootCACertPath != "" {
	//	err = soapClient.SetRootCAs(rootCACertPath)
	//	if err != nil {
	//		return nil, err
	//	}
	//}

	return soapClient, nil
}

func instantCloneVm(sourceVmInventoryPath, targetVmName string) error {

	ctx := context.Background()
	// login

	// managerFactory.soapClient() creates a SC
	// soapClient(ctx, sc) calls NewClient

	// NewClient(ctx, rt) returns a vim25.Client
	// vim25.Client is a soap.RoundTripper
	//mf := vcenter_client_factory.ManagerFactory{
	//	vcenter_client_factory.FactoryConfig{
	//
	//	}
	//}
	soapClient, err := buildSoapClient()

	vim25Client, err := vim25.NewClient(ctx, soapClient)
	if err != nil {
		return fmt.Errorf("error building vim25 client: %s", err)
	}
	govmomiClient := govmomi.Client{
		Client:         vim25Client,
		SessionManager: session.NewManager(vim25Client),
	}

	creds := url.UserPassword(vCenterUsername, vCenterPassword)

	err = govmomiClient.Login(ctx, creds)
	if err != nil {
		return fmt.Errorf("error logging in to vcenter: %s", err)
	}

	// get vm to clone
	recurse := false
	finder := find.NewFinder(vim25Client, recurse)

	vm, err := finder.VirtualMachine(ctx, sourceVmInventoryPath)
	if err != nil {
		return fmt.Errorf("could not find VM: %s", err)
	}

	//cloneConfig :=
	//ipAddressConfig := []types.BaseOptionValue{
	//	&types.OptionValue{
	//		Key: "ipAddress",
	//		Value: "10.85.81.5",
	//	},
	//}
	req := types.InstantClone_Task{
		This: vm.Reference(),
		Spec: types.VirtualMachineInstantCloneSpec{
			Name:     targetVmName,
			Location: types.VirtualMachineRelocateSpec{},
			//Config: ipAddressConfig,
		},
	}

	// v.c is off a client
	_, err = methods.InstantClone_Task(ctx, vim25Client, &req)
	if err != nil {
		return fmt.Errorf("failed to instant-clone: %s", err)
	}

	//return NewTask(v.c, res.Returnval), nil
	return nil
}

var _ = Describe("stembuild construct", func() {
	var workingDir string

	BeforeEach(func() {
		var err error
		workingDir, err = os.Getwd()
		Expect(err).ToNot(HaveOccurred())

	})

	const constructOutputTimeout = 60 * time.Second
	Context("run successfully", func() {

		It("successfully exits when vm becomes powered off", func() {
			err := CopyFile(filepath.Join(workingDir, "assets", "LGPO.zip"), filepath.Join(workingDir, "LGPO.zip"))
			Expect(err).ToNot(HaveOccurred())

			session := helpers.Stembuild(stembuildExecutable, "construct", "-vm-ip", conf.TargetIP, "-vm-username", conf.VMUsername, "-vm-password", conf.VMPassword, "-vcenter-url", conf.VCenterURL, "-vcenter-username", conf.VCenterUsername, "-vcenter-password", conf.VCenterPassword, "-vm-inventory-path", conf.VMInventoryPath)

			shutdownTimeout := 3 * time.Minute
			Eventually(session, shutdownTimeout).Should(Exit(0))
		})

		It("transfers LGPO and StemcellAutomation archives, unarchive them and execute automation script", func() {
			err := CopyFile(filepath.Join(workingDir, "assets", "LGPO.zip"), filepath.Join(workingDir, "LGPO.zip"))
			Expect(err).ToNot(HaveOccurred())

			session := helpers.Stembuild(stembuildExecutable, "construct", "-vm-ip", conf.TargetIP, "-vm-username", conf.VMUsername, "-vm-password", conf.VMPassword, "-vcenter-url", conf.VCenterURL, "-vcenter-username", conf.VCenterUsername, "-vcenter-password", conf.VCenterPassword, "-vm-inventory-path", conf.VMInventoryPath)

			Eventually(session.Out, constructOutputTimeout).Should(Say(`mock stemcell automation script executed`))
		})

		It("extracts the WinRM BOSH powershell script and executes it successfully on the guest VM", func() {
			err := CopyFile(filepath.Join(workingDir, "assets", "LGPO.zip"), filepath.Join(workingDir, "LGPO.zip"))
			Expect(err).ToNot(HaveOccurred())

			session := helpers.Stembuild(stembuildExecutable, "construct", "-vm-ip", conf.TargetIP, "-vm-username", conf.VMUsername, "-vm-password", conf.VMPassword, "-vcenter-url", conf.VCenterURL, "-vcenter-username", conf.VCenterUsername, "-vcenter-password", conf.VCenterPassword, "-vm-inventory-path", conf.VMInventoryPath)

			Eventually(session.Out, constructOutputTimeout).Should(Say(`Attempting to enable WinRM on the guest vm...WinRm enabled on the guest VM`))

		})

		It("handles special characters", func() {
			isAlphaNumeric, err := regexp.Compile("[a-zA-Z0-9]+")
			Expect(err).ToNot(HaveOccurred())

			if isAlphaNumeric.MatchString(conf.VCenterUsername) && isAlphaNumeric.MatchString(conf.VCenterPassword) {
				Skip("vCenter username or password must contain special characters")
			}
			err = CopyFile(filepath.Join(workingDir, "assets", "LGPO.zip"), filepath.Join(workingDir, "LGPO.zip"))
			Expect(err).ToNot(HaveOccurred())

			session := helpers.Stembuild(stembuildExecutable, "construct", "-vm-ip", conf.TargetIP, "-vm-username", conf.VMUsername, "-vm-password", conf.VMPassword, "-vcenter-url", conf.VCenterURL, "-vcenter-username", conf.VCenterUsername, "-vcenter-password", conf.VCenterPassword, "-vm-inventory-path", conf.VMInventoryPath)

			Eventually(session, constructOutputTimeout).Should(Exit(0))
			Eventually(session.Out).Should(Say(`mock stemcell automation script executed`))
		})

		FIt("successfully runs even when a user has logged in", func() {
			//endpoint := winrm.NewEndpoint(conf.TargetIP, 5985, false, true, nil, nil, nil, 0)
			// new client -> visually: not logged in. test behavior: ?
			//client, err := winrm.NewClient(endpoint, conf.VMUsername, conf.VMPassword)
			//
			//// new shell on the client: visually: not logged in. test behavior: ?
			//shell, err := client.CreateShell()
			//
			//// execute something
			//
			//// execute a long-running something
			//// we tried shell.Execute(timeout) but no error occurred (so does not
			//// accurately simulate user is logged in
			//shell.Execute("timeout 600 /nobreak")

			//
			// can we login some other way (send Ctrl-Alt-Del, etc.): govc?

			//Fail("call Instant clone vm correctly, and stembuild against a vsphere 6.7 environment")

			//err := instantCloneVm("/private/vm/cove/stembuild_base_vm_2019_22", "cloned-vm-from-test")
			//Expect(err).NotTo(HaveOccurred())

			//Fail("cloned VM")

			fmt.Printf("please suspend source vm\n")
			time.Sleep(30 * time.Second)

			// run normal stembuild construct command, like we do in prev. test
			err := CopyFile(filepath.Join(workingDir, "assets", "LGPO.zip"), filepath.Join(workingDir, "LGPO.zip"))
			Expect(err).ToNot(HaveOccurred())

			session := helpers.Stembuild(stembuildExecutable, "construct",
				"-vm-ip", "10.85.81.22",
				"-vm-username", "Pivotal",
				"-vm-password", "password123!",
				"-vcenter-url", "vcenter.wild.cf-app.com",
				"-vcenter-username", vCenterUsername,
				"-vcenter-password", vCenterPassword,
				"-vm-inventory-path", "/private/vm/cove/stembuild_base_vm_2019_22")
			//"-vm-inventory-path", "/private/vm/cove/cloned-vm-from-test")

			// assuming old, pre-story state
			// expect timeout
			shutdownTimeout := 3 * time.Minute
			Eventually(session, shutdownTimeout).Should(Exit(0))
			Expect(err).NotTo(HaveOccurred())
			//time.Sleep(time.Duration(1 * time.Minute))
		})
	})

	It("fails with an appropriate error when LGPO is missing", func() {
		session := helpers.Stembuild(stembuildExecutable, "construct", "-vm-ip", conf.TargetIP, "-vm-username", conf.VMUsername, "-vm-password", conf.VMPassword, "-vcenter-url", conf.VCenterURL, "-vcenter-username", conf.VCenterUsername, "-vcenter-password", conf.VCenterPassword, "-vm-inventory-path", conf.VMInventoryPath)

		Eventually(session, constructOutputTimeout).Should(Exit(1))
		Eventually(session.Err).Should(Say(`Could not find LGPO.zip in the current directory`))
	})

	It("does not exit when the target VM has not powered off", func() {
		err := CopyFile(filepath.Join(workingDir, "assets", "LGPO.zip"), filepath.Join(workingDir, "LGPO.zip"))
		Expect(err).ToNot(HaveOccurred())

		fakeStemcellAutomationShutdownDelay := 45 * time.Second

		session := helpers.Stembuild(stembuildExecutable, "construct", "-vm-ip", conf.TargetIP, "-vm-username", conf.VMUsername, "-vm-password", conf.VMPassword, "-vcenter-url", conf.VCenterURL, "-vcenter-username", conf.VCenterUsername, "-vcenter-password", conf.VCenterPassword, "-vm-inventory-path", conf.VMInventoryPath)

		Consistently(session, fakeStemcellAutomationShutdownDelay-5*time.Second).Should(Not(Exit()))
	})

	AfterEach(func() {
		_ = os.Remove(filepath.Join(workingDir, "LGPO.zip"))
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

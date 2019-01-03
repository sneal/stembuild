package construct_test

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/cloudfoundry-incubator/stembuild/test/helpers"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"text/template"
	"time"

	"github.com/masterzen/winrm"

	"github.com/concourse/pool-resource/out"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vmware/govmomi/govc/cli"
	_ "github.com/vmware/govmomi/govc/importx"
	_ "github.com/vmware/govmomi/govc/vm"
)

func TestConstruct(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Construct Suite")
}

const (
	NetworkGatewayVariable = "CONSTRUCT_TEST_GATEWAY"
	SubnetMaskVariable     = "CONSTRUCT_TEST_SUBNET_MASK"
	OvaFileVariable        = "OVA_FILE"
	VMNamePrefixVariable   = "VM_NAME_PREFIX"
	VMUsernameVariable     = "VM_USERNAME"
	VMPasswordVariable     = "VM_PASSWORD"
	ExistingVmIPVariable   = "EXISTING_VM_IP"
	UserProvidedIPVariable = "USER_PROVIDED_IP"
	LockPrivateKeyVariable = "LOCK_PRIVATE_KEY"
)

var (
	conf                config
	tmpDir              string
	lockPool            out.LockPool
	lockDir             string
	stembuildExecutable string
)

//TODO: separate VM attributes from VCenter configuration
type config struct {
	TargetIP       string
	NetworkGateway string
	SubnetMask     string
	VMUsername     string
	VMPassword     string
	VMName         string
}

func envMustExist(variableName string) string {
	result := os.Getenv(variableName)
	if result == "" {
		Fail(fmt.Sprintf("%s must be set", variableName))
	}

	return result
}

func claimAvailableIP() string {

	lockPrivateKey := envMustExist(LockPrivateKeyVariable)
	keyFile, err := ioutil.TempFile(os.TempDir(), "keyfile")
	Expect(err).NotTo(HaveOccurred())
	keyFile.Write([]byte(lockPrivateKey))
	keyFile.Chmod(0600)

	err = exec.Command("ssh-add", keyFile.Name()).Run()
	Expect(err).NotTo(HaveOccurred())

	poolSource := out.Source{
		URI:        "git@github.com:pivotal-cf-experimental/Bosh-Windows-Locks.git",
		Branch:     "master",
		Pool:       "vcenter-ips",
		RetryDelay: 5 * time.Second,
	}

	buffer := bytes.Buffer{}
	writer := bufio.NewWriter(&buffer)

	lockPool = out.NewLockPool(poolSource, writer)

	ip, _, err := lockPool.AcquireLock()
	Expect(err).NotTo(HaveOccurred())
	Expect(ip).NotTo(Equal(""))

	lockDir, err = ioutil.TempDir("", "acquired-lock")
	Expect(err).NotTo(HaveOccurred())
	err = ioutil.WriteFile(filepath.Join(lockDir, "name"), []byte(ip), os.ModePerm)
	Expect(err).NotTo(HaveOccurred())

	return ip
}

var _ = SynchronizedBeforeSuite(func() []byte {
	var err error
	stembuildExecutable, err = helpers.BuildStembuild()
	Expect(err).NotTo(HaveOccurred())

	// Build a VM and wait for it's IP

	networkGateway := envMustExist(NetworkGatewayVariable)
	subnetMask := envMustExist(SubnetMaskVariable)
	ovaFile := envMustExist(OvaFileVariable)

	vmNamePrefix := envMustExist(VMNamePrefixVariable)
	vmUsername := envMustExist(VMUsernameVariable)
	vmPassword := envMustExist(VMPasswordVariable)

	targetIP := os.Getenv(ExistingVmIPVariable) //TODO: make a boolean if existing machine should be used for readaibility
	conf = config{
		TargetIP:       targetIP,
		NetworkGateway: networkGateway,
		SubnetMask:     subnetMask,
		VMUsername:     vmUsername,
		VMPassword:     vmPassword,
	}

	if targetIP == "" {

		fmt.Println("No existing VM IP given")

		givenIP := os.Getenv(UserProvidedIPVariable)
		if givenIP != "" {
			targetIP = givenIP
		}

		fmt.Println("No user-provided IP given. Finding available IP...")
		targetIP = claimAvailableIP()
		conf.TargetIP = targetIP
		fmt.Printf("Target ip is %s\n", targetIP)

		vmNameSuffix := strings.Split(targetIP, ".")[3]
		vmName := fmt.Sprintf("%s%s", vmNamePrefix, vmNameSuffix)
		conf.VMName = vmName

		templateFile, err := filepath.Abs("assets/ova_options.json.template")
		Expect(err).NotTo(HaveOccurred())
		tmpl, err := template.New("ova_options.json.template").ParseFiles(templateFile)

		tmpDir, err = ioutil.TempDir("", "construct-test")
		Expect(err).NotTo(HaveOccurred())

		tmpFile, err := ioutil.TempFile(tmpDir, "ova_options*.json")
		Expect(err).NotTo(HaveOccurred())

		err = tmpl.Execute(tmpFile, conf)
		Expect(err).NotTo(HaveOccurred())

		opts := []string{
			"import.ova",
			fmt.Sprintf("--options=%s", tmpFile.Name()),
			fmt.Sprintf("--name=%s", vmName),
			"--folder=/canada-dc/vm/winnipeg",
			ovaFile,
		}

		fmt.Printf("Opts are %s", opts)

		exitCode := cli.Run(opts)
		Expect(exitCode).To(BeZero())
	}

	endpoint := winrm.NewEndpoint(targetIP, 5985, false, true, nil, nil, nil, 0)
	client, err := winrm.NewClient(endpoint, conf.VMUsername, conf.VMPassword)
	Expect(err).NotTo(HaveOccurred())

	var shell *winrm.Shell
	Eventually(func() error {
		shell, err = client.CreateShell()
		return err
	}, 3*time.Minute).Should(BeNil()) //TODO: Should this be different depending on whether we just created the VM or using an existing one Will need more time?
	shell.Close()
	fmt.Println("VM created and connected")

	return nil
}, func(_ []byte) {
})

var _ = SynchronizedAfterSuite(func() {
	os.RemoveAll(tmpDir)
	return
	if conf.TargetIP != "" {

		delete_command := []string{"vm.destroy", fmt.Sprintf("-vm.ip=%s", conf.TargetIP)}
		Eventually(func() int {
			return cli.Run(delete_command)
		}, 3*time.Minute).Should(BeZero())
		fmt.Println("VM destroyed")
		if lockDir != "" {
			_, _, err := lockPool.ReleaseLock(lockDir)
			Expect(err).NotTo(HaveOccurred())

			tmpDir := os.TempDir()
			childItems, err := ioutil.ReadDir(tmpDir)
			Expect(err).NotTo(HaveOccurred())

			for _, item := range childItems {
				if item.IsDir() && strings.HasPrefix(filepath.Base(item.Name()), "pool-resource") {
					fmt.Printf("Cleaning up temporary pool resource %s\n", item.Name())
					os.RemoveAll(item.Name())
				}
			}
		}
	}
}, func() {
	Expect(os.RemoveAll(stembuildExecutable)).To(Succeed())
})

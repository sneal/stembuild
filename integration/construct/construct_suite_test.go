package construct_test

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"text/template"
	"time"

	"github.com/concourse/pool-resource/out"

	"github.com/masterzen/winrm"

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
	TargetIPRangeVariable    = "CONSTRUCT_TEST_IP_RANGE"
	NetworkGatewayVariable   = "CONSTRUCT_TEST_GATEWAY"
	SubnetMaskVariable       = "CONSTRUCT_TEST_SUBNET_MASK"
	OvaFileVariable          = "OVA_FILE"
	VMNamePrefixVariable     = "VM_NAME_PREFIX"
	VMUsernameVariable       = "VM_USERNAME"
	VMPasswordVariable       = "VM_PASSWORD"
	ExistingTargetIPVariable = "EXISTING_TARGET_IP"
	TargetIPVariable         = "TARGET_IP"
	LockPrivateKeyVariable   = "LOCK_PRIVATE_KEY"
)

var (
	conf     config
	tmpDir   string
	lockPool out.LockPool
)

type config struct {
	TargetIP       string
	NetworkGateway string
	SubnetMask     string
	VMName         string
}

func envMustExist(variableName string) string {
	result := os.Getenv(variableName)
	if result == "" {
		Fail(fmt.Sprintf("%s must be set", variableName))
	}

	return result
}

func findAvailableIP(cidr string) string {

	//If specific IP is given
	//TODO: maybe if IP is given, don't enter this function or rename function name
	givenIP := os.Getenv(TargetIPVariable)
	if givenIP != "" {
		return givenIP
	}

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
	return ip
}

var _ = BeforeSuite(func() {
	// Build a VM and wait for it's IP

	targetIPRange := envMustExist(TargetIPRangeVariable)
	networkGateway := envMustExist(NetworkGatewayVariable)
	subnetMask := envMustExist(SubnetMaskVariable)
	ovaFile := envMustExist(OvaFileVariable)

	vmNamePrefix := envMustExist(VMNamePrefixVariable)
	vmUsername := envMustExist(VMUsernameVariable)
	vmPassword := envMustExist(VMPasswordVariable)

	targetIP := os.Getenv(ExistingTargetIPVariable)

	if targetIP == "" {

		targetIP = findAvailableIP(targetIPRange)
		//TODO: Delete below. Just for testing until we fix the fping issue
		fmt.Printf("Target ip is %s\n", targetIP)

		vmNameSuffix := strings.Split(targetIP, ".")[3]
		vmName := fmt.Sprintf("%s%s", vmNamePrefix, vmNameSuffix)

		conf = config{
			TargetIP:       targetIP,
			NetworkGateway: networkGateway,
			SubnetMask:     subnetMask,
			VMName:         vmName,
		}

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
	client, err := winrm.NewClient(endpoint, vmUsername, vmPassword)

	fmt.Println(endpoint)
	fmt.Println(client)
	Expect(err).NotTo(HaveOccurred())

	_, err = client.CreateShell()
	Expect(err).NotTo(HaveOccurred())

	//Eventually(func() error {
	//}, 5*time.Second).Should(BeNil()) //TODO: Should this be different depending on whether we just created the VM or using an existing one Will need more time?

})

var _ = AfterSuite(func() {
	//os.RemoveAll(tmpDir)
	//TODO: power off VM, delete from disk, and unclaim IP from locks!

	if conf.TargetIP != "" {
		exitCode := cli.Run([]string{"vm.destroy", fmt.Sprintf("-vm.ip=%s", conf.TargetIP)})
		Expect(exitCode).To(BeZero())

		_, _, err := lockPool.ReleaseLock(conf.TargetIP)
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
})

package cleanup_vm

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/concourse/pool-resource/out"

	"github.com/cloudfoundry-incubator/stembuild/test/helpers"

	. "github.com/onsi/gomega"
)

type TestResources struct {
	LockDir       string
	LockPool      out.LockPool
	LockParentDir string
	tr
}

func Cleanup(vcenterCredentialUrl string, targetIP string, locks TestResources) {
	println("cleanup called")

	deleteCommand := []string{
		"vm.destroy",
		fmt.Sprintf("-u=%s", vcenterCredentialUrl),
		fmt.Sprintf("-vm.ip=%s", targetIP),
	}
	Eventually(func() int {
		return helpers.RunIgnoringOutput(deleteCommand)
	}, 3*time.Minute, 10*time.Second).Should(BeZero())
	fmt.Println("VM destroyed")
	if locks.LockDir != "" {
		_, _, err := locks.LockPool.ReleaseLock(locks.LockDir)
		Expect(err).NotTo(HaveOccurred())

		childItems, err := ioutil.ReadDir(locks.LockParentDir)
		Expect(err).NotTo(HaveOccurred())

		for _, item := range childItems {
			if item.IsDir() && strings.HasPrefix(filepath.Base(item.Name()), "pool-resource") {
				fmt.Printf("Cleaning up temporary pool resource %s\n", item.Name())
				_ = os.RemoveAll(item.Name())
			}
		}
	}
}

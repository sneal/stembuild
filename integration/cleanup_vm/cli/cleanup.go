package main

import "github.com/cloudfoundry-incubator/stembuild/integration/cleanup_vm"

// receive
func main() {
	cleanup_vm.Cleanup()
}

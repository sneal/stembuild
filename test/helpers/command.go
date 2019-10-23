package helpers

import (
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/vmware/govmomi/govc/cli"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
)

const (
	DebugCommandPrefix = "\nCMD>"
	DebugOutPrefix     = "OUT: "
	DebugErrPrefix     = "ERR: "
)

func Stembuild(command string, args ...string) *Session {
	WriteCommand(command, args)
	session, err := Start(
		exec.Command(command, args...),
		NewPrefixedWriter(DebugOutPrefix, GinkgoWriter),
		NewPrefixedWriter(DebugErrPrefix, GinkgoWriter))
	Expect(err).NotTo(HaveOccurred())
	return session
}

func RunCommandInDir(workingDir, command string, args ...string) *Session {
	WriteCommand(command, args)
	session, err := Start(
		&exec.Cmd{
			Path: command,
			Args: append([]string{command}, args...),
			Dir:  workingDir,
		},
		NewPrefixedWriter(DebugOutPrefix, GinkgoWriter),
		NewPrefixedWriter(DebugErrPrefix, GinkgoWriter))
	Expect(err).NotTo(HaveOccurred())
	return session
}

func RunIgnoringOutput(args []string) int {
	oldStderr := os.Stderr
	oldStdout := os.Stdout

	_, w, _ := os.Pipe()

	defer w.Close()

	os.Stderr = w
	os.Stdout = w

	os.Stderr = os.NewFile(uintptr(syscall.Stderr), "/dev/null")

	exitCode := cli.Run(args)

	os.Stderr = oldStderr
	os.Stdout = oldStdout

	return exitCode
}

func WriteCommand(command string, args []string) {
	display := append([]string{DebugCommandPrefix, command}, args...)
	GinkgoWriter.Write([]byte(strings.Join(append(display, "\n"), " ")))
}

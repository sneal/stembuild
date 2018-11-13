package commandparser

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/subcommands"
	"github.com/pivotal-cf-experimental/stembuild/colorlogger"
)

type ConstructCmd struct {
	version     string
	GlobalFlags *GlobalFlags
}

func (*ConstructCmd) Name() string     { return "construct" }
func (*ConstructCmd) Synopsis() string { return "Download stemcell automation articfact and lgpo" }
func (*ConstructCmd) Usage() string {
	return fmt.Sprintf(`%[1]s construct -version <stemcell version>

Create a BOSH Stemcell from a VMDK file

The [vmdk], [version], and [os] flags must be specified.  If the [output] flag is
not specified the stemcell will be created in the current working directory.

Requirements:
	The VMware 'ovftool' binary must be on your path or Fusion/Workstation
	must be installed (both include the 'ovftool').

Examples:
	%[1]s -vmdk disk.vmdk -version 1.2 -os 1803

	Will create an Windows 1803 stemcell using [vmdk] 'disk.vmdk', and set the stemcell version to 1.2.
	The final stemcell will be found in the current working directory.

Flags:
`, filepath.Base(os.Args[0]))
}

func (p *ConstructCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&p.version, "version", "", "Stemcell version in the form of [DIGITS].[DIGITS] (e.g. 123.01)")
	f.StringVar(&p.version, "v", "", "Stemcell version (shorthand)")
}
func (p *ConstructCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	logLevel := colorlogger.NONE
	if p.GlobalFlags.Debug {
		logLevel = colorlogger.DEBUG
	}
	logger := colorlogger.ConstructLogger(logLevel, p.GlobalFlags.Color, os.Stderr)

	logger.Debugf("hello, world.")
	if !IsValidVersion(p.version) {
		fmt.Fprintf(os.Stderr, "invalid version (%s) expected format [NUMBER].[NUMBER] or "+
			"[NUMBER].[NUMBER].[NUMBER]\n", p.version)

		return subcommands.ExitFailure
	}

	return subcommands.ExitSuccess
}

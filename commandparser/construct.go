package commandparser

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cloudfoundry-incubator/stembuild/colorlogger"
	"github.com/google/subcommands"
)

type ConstructCmd struct {
	stemcellVersion string
	winrmUsername   string
	winrmPassword   string
	winrmIP         string
	GlobalFlags     *GlobalFlags
}

func (*ConstructCmd) Name() string     { return "construct" }
func (*ConstructCmd) Synopsis() string { return "Transfer automation artifact and LGPO to vCenter" }

//TODO: REWRITE USAGE
func (*ConstructCmd) Usage() string {
	return fmt.Sprintf(`%[1]s construct -stemcellVersion <stemcell stemcellVersion>

The [stemcellVersion], [winrmUsername], [winrmPassword], [winrmIP] flags must be specified.

Requirements:
	The VMware 'ovftool' binary must be on your path or Fusion/Workstation
	must be installed (both include the 'ovftool').

Examples:
	%[1]s -vmdk disk.vmdk -stemcellVersion 1.2 -os 1803

	Will create an Windows 1803 stemcell using [vmdk] 'disk.vmdk', and set the stemcell stemcellVersion to 1.2.
	The final stemcell will be found in the current working directory.

Flags:
`, filepath.Base(os.Args[0]))
}

func (p *ConstructCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&p.stemcellVersion, "stemcellVersion", "", "Stemcell version in the form of [DIGITS].[DIGITS] (e.g. 123.01)")
	f.StringVar(&p.stemcellVersion, "s", "", "Stemcell version (shorthand)")
	f.StringVar(&p.winrmUsername, "winrmUsername", "", "Example: ")
	f.StringVar(&p.winrmUsername, "u", "", "winrmUsername (shorthand)")
	f.StringVar(&p.winrmPassword, "winrmPassword", "", "Example: ")
	f.StringVar(&p.winrmPassword, "p", "", "winrmPassword (shorthand)")
	f.StringVar(&p.winrmIP, "winrmIP", "", "Example: ")
	f.StringVar(&p.winrmIP, "ip", "", "winrmIP (shorthand)")

}
func (p *ConstructCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	logLevel := colorlogger.NONE
	if p.GlobalFlags.Debug {
		logLevel = colorlogger.DEBUG
	}
	logger := colorlogger.ConstructLogger(logLevel, p.GlobalFlags.Color, os.Stderr)
	logger.Debugf("hello, world.")
	if !IsValidStemcellVersion(p.stemcellVersion) {
		fmt.Fprintf(os.Stderr, "invalid stemcellVersion (%s) expected format [NUMBER].[NUMBER] or "+
			"[NUMBER].[NUMBER].[NUMBER]\n", p.stemcellVersion)

		return subcommands.ExitFailure
	}

	pwd, err := os.Getwd()
	if err != nil {
		fmt.Fprint(os.Stderr, "unable to find current working directory", err)
		return subcommands.ExitFailure
	}
	automationArtifactPresent, err := IsArtifactInDirectory(pwd, "StemcellAutomation.zip")
	if !automationArtifactPresent {
		fmt.Fprintf(os.Stderr, "automation artifact not found in current directory")
		return subcommands.ExitFailure
		//TODO: Download automation Artifact
	}
	lgpoPresent, err := IsArtifactInDirectory(pwd, "LGPO.zip")
	if !lgpoPresent {
		fmt.Fprintf(os.Stderr, "lgpo not found in current directory")
		return subcommands.ExitFailure
		//TODO: Download LGPO
	}

	remoteManager := NewWinRM(p.winrmIP, p.winrmUsername, p.winrmPassword)
	fmt.Printf("upload artifact...")
	err = UploadArtifact(remoteManager)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		return subcommands.ExitFailure
	}
	fmt.Printf("extract artifact...")
	err = ExtractArchive(remoteManager)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		return subcommands.ExitFailure
	}
	fmt.Printf("execute script...")
	err = ExecuteSetupScript(remoteManager)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		return subcommands.ExitFailure
	}

	return subcommands.ExitSuccess
}

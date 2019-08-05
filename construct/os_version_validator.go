package construct

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/cloudfoundry-incubator/stembuild/version"
)

//go:generate counterfeiter . OSValidatorMessenger
type OSValidatorMessenger interface {
	OSVersionFileCreationFailed(errorMessage string)
	ExitCodeRetrievalFailed(errorMessage string)
	DownloadFileFailed(errorMessage string)
}

type OSVersionValidator struct {
	GuestManager GuestManager
	Messenger    OSValidatorMessenger
}

func (v *OSVersionValidator) Validate() error {
	//fmt.Println("before startprogram")
	pid, err := v.GuestManager.StartProgramInGuest(
		context.Background(),
		powershell,
		"[System.Environment]::OSVersion.Version.Build > C:\\Windows\\Temp\\version.log",
	)
	if err != nil {
		v.Messenger.OSVersionFileCreationFailed(err.Error())
		return nil
	}

	fmt.Println("before exit code")
	exitCode, err := v.GuestManager.ExitCodeForProgramInGuest(context.Background(), pid)
	if err != nil {
		v.Messenger.ExitCodeRetrievalFailed(err.Error())
		return nil
	}
	fmt.Println("no exit code error")
	if exitCode != 0 {
		v.Messenger.OSVersionFileCreationFailed(fmt.Sprintf("OS version file creation failed with non-zero exit code: %d", exitCode))
		return nil
	}

	fmt.Println("before download")

	fileReader, _, err := v.GuestManager.DownloadFileInGuest(context.Background(), "C:\\Windows\\Temp\\version.log")
	if err != nil {
		v.Messenger.DownloadFileFailed(err.Error())
		return nil
	}
	buf, err := ioutil.ReadAll(fileReader)
	if err != nil {
		v.Messenger.DownloadFileFailed(err.Error())
		return nil
	}
	fmt.Println("before check")
	guestOSVersion := version.GetOSVersionFromBuildNumber(string(buf))
	fmt.Printf("Guest OS Version:'%s', Stembuild Version:'%s'\n", guestOSVersion, version.Version)
	if guestOSVersion != version.Version {
		fmt.Println("error vesrion mismatch")
		return fmt.Errorf("OS version of stembuild and guest OS VM do not match. Guest OS Version:'%s', Stembuild Version:'%s'", guestOSVersion, version.Version)
	}

	fmt.Println("version match")
	return nil
}

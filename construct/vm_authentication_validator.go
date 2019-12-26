package construct

import (
	"context"
	"strings"
)

type VMAuthValidator struct {
	GuestManager GuestManager
}

func (v *VMAuthValidator) IsValidAuth() (bool, error) {
	err := v.runThrowawayCommand()
	if err != nil {
		if strings.Contains(err.Error(), "Failed to authenticate with the guest operating system") {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}

func (v *VMAuthValidator) runThrowawayCommand() error {
	_, err := v.GuestManager.StartProgramInGuest(context.Background(), powershell, "ls")
	return err
}

package commandparser

import (
	"github.com/cloudfoundry-incubator/winrmcp/winrmcp"
	"time"
)

type WinRm struct {
	host     string
	username string
	password string
}

func NewWinRm() *WinRm {
	return &WinRm{}
}

func (w * WinRm) UploadArtifact(sourceFilePath, destinationFilePath string) error {
	client, err := winrmcp.New(w.host, &winrmcp.Config{
		Auth:                  winrmcp.Auth{User: w.username, Password: w.password},
		Https:                 false,
		Insecure:              true,
		OperationTimeout:      time.Second * 60,
		MaxOperationsPerShell: 15,
	})


	if err != nil {
		return err
	}
	return client.Copy(sourceFilePath, destinationFilePath)
}
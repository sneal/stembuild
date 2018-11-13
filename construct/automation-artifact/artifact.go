package artifact

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/google/go-github/github"
)

func DownloadAutomationArtifact(version, path string) error {
	owner := "cloudfoundry-incubator"
	repositoryName := "bosh-windows-stemcell-automation"

	client := github.NewClient(nil)

	assetID, err := getReleaseAssetID(client, owner, repositoryName, version)

	readCloser, redirectURL := getRelease(client, assetID)

	if redirectURL == "" {
		getReleaseFrom(readCloser, path, err)
	} else {
		getReleaseByRedirectURL(path, redirectURL)
	}

	return nil
}

func getRelease(client *github.Client, assetID int64) (io.ReadCloser, string) {
	readCloser, redirectURL, err := client.Repositories.DownloadReleaseAsset(context.Background(), "cloudfoundry-incubator", "bosh-windows-stemcell-automation", assetID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "download release asset failed %s", err.Error())
	}
	return readCloser, redirectURL
}

func getReleaseAssetID(client *github.Client, owner string, repositoryName string, version string) (int64, error) {
	if version != "" {
		release, _, err := client.Repositories.GetReleaseByTag(context.Background(), owner, repositoryName, version)
		if err != nil {
			fmt.Fprintf(os.Stderr, "get release by tag failed %s", err.Error())
			println("")
		}
		return *release.Assets[0].ID, err
	}else{
		release, _, err := client.Repositories.GetLatestRelease(context.Background(), owner, repositoryName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "get release by tag failed %s", err.Error())
			println("")
		}
		return *release.Assets[0].ID, err
	}
}

func getReleaseByRedirectURL(path string, redirectURL string) {
	filepath := path + "StemcellAutomation.zip"
	out, err := os.Create(filepath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not craete file %s", err.Error())
	}
	defer out.Close()
	// Get the data
	resp, err := http.Get(redirectURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not get contetnt from url %s", err.Error())
	}
	defer resp.Body.Close()
	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not write to file %s", err.Error())
	}
}

func getReleaseFrom(readCloser io.ReadCloser, path string, err error) {
	var ba []byte
	readCloser.Read(ba)
	ioutil.WriteFile(path, ba, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not copy file %s", err.Error())
	}
}

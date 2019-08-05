package version

import (
	"strings"
)

var Version = "dev"
var oSBuildNumberToVersionMap = map[string]string{"17763": "2019", "17134": "1803", "dev":"dev"}


func GetVersions(mainVersion string) (string, string) {
	stringArr := strings.Split(mainVersion, ".")

	os := stringArr[0]
	switch os {
	case "1709":
		os = "2016"
	case "1200":
		os = "2012R2"
	}
	stemcellVersion := strings.Join(stringArr[0:2], ".")

	return os, stemcellVersion
}

func GetOSVersionFromBuildNumber(version string) string {
	return oSBuildNumberToVersionMap[version]
}

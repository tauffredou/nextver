package provider

import (
	"fmt"
	"regexp"
	"strconv"
)

const (
	UNDEFINED = 0
	PATCH     = 1
	MINOR     = PATCH << 1
	MAJOR     = MINOR << 1
)
const (
	SemverRegex              = `^v?(\d+)(\.(\d)+)?(\.(\d)+)?`
	DateRegexp               = `\d{4}-\d{2}-\d{2}-\d{6}`
	ConventionalCommitRegexp = `^([a-zA-Z-_]+)(\(([^\):]+)\))?:? ?(.*)$`
	FirstVersion             = "0.0.0"
)

type Provider interface {
	GetLatestRelease() Release
}

func ReadSemver(v string) ([]int64, error) {
	re := regexp.MustCompile(SemverRegex)
	if re.MatchString(v) {
		data := re.FindStringSubmatch(v)
		major, _ := strconv.ParseInt(data[1], 10, 0)
		minor, _ := strconv.ParseInt(data[3], 10, 0)
		patch, _ := strconv.ParseInt(data[5], 10, 0)
		return []int64{major, minor, patch}, nil
	} else {
		return nil, fmt.Errorf("cannot read version")
	}

}

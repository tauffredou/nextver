package model

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const DATE_FORMAT = "2006-01-02-150405"

func SemverCalculator(r *Release) (string, error) {
	mmp, err := ReadSemver(r.CurrentVersion)
	if err != nil {
		return "", fmt.Errorf("cannot calculate next version. version: %s", r.CurrentVersion)
	}

	var mask byte = 0
	for i := range r.Changelog {
		mask = mask | r.Changelog[i].Level
	}

	switch {
	case mask&MAJOR == MAJOR:
		mmp[0] += 1
		mmp[1] = 0
		mmp[2] = 0
	case mask&MINOR == MINOR:
		mmp[1] += 1
		mmp[2] = 0
	case mask&PATCH == PATCH:
		mmp[2] += 1
	}
	version := fmt.Sprintf("%d.%d.%d", mmp[0], mmp[1], mmp[2])
	return strings.ReplaceAll(r.VersionPattern, "SEMVER", version), nil
}

func DateVersionCalculator(r *Release) (string, error) {
	t := time.Now()

	date := t.Format(DATE_FORMAT)

	return strings.ReplaceAll(r.VersionPattern, "DATE", date), nil
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

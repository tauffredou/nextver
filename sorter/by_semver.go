package sorter

import (
	"github.com/tauffredou/nextver/model"
)

type BySemver []model.Release

func (a BySemver) Len() int      { return len(a) }
func (a BySemver) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a BySemver) Less(i, j int) bool {
	semverI, errI := model.ReadSemver(a[i].CurrentVersion)
	semverJ, errJ := model.ReadSemver(a[j].CurrentVersion)

	if errI != nil || errJ != nil {
		return a[i].CurrentVersion > a[j].CurrentVersion
	}

	return semverI[0] > semverJ[0] ||
		semverI[1] > semverJ[1] ||
		semverI[2] > semverJ[2]
}

package model_test

import (
	"github.com/stretchr/testify/assert"
	"testing"
	  "github.com/tauffredou/nextver/model"
"fmt"
)

func TestReadSemver(t *testing.T) {
	tests := []struct {
		name    string
		version    string
		want    interface{}
		wantErr error
	}{
		{name: "semver prefix", version: "v1.2.0", want: []int64{1,2,0}},
		{name: "semver", version: "1.2.0", want: []int64{1,2,0}},
		{name: "semver", version: "v1.0.1", want: []int64{1,0,1}},
		{name: "date", version: "2006-01-02-150405", wantErr: fmt.Errorf("cannot read version")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			semver, err := model.ReadSemver(test.version)

			if test.wantErr != nil {
				assert.Equal(t, test.wantErr, err)
			} else {
				assert.Equal(t, test.want, semver)
			}
		})
	}

}


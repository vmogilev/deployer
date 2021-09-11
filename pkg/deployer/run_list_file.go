package deployer

import (
	"fmt"
	"strconv"
	"strings"
)

type RunListFile struct {
	Name      string
	id        int
	directive string
}

const (
	RunListFileExtension = ".json"
	RunListFileSeparator = "_"
)

func (f *RunListFile) Parse() error {
	if !strings.HasSuffix(f.Name, RunListFileExtension) {
		return fmt.Errorf("run-list file %s has wrong extension", f.Name)
	}

	n := strings.TrimSuffix(f.Name, RunListFileExtension)
	pp := strings.Split(n, RunListFileSeparator)
	if len(pp) != 2 {
		return fmt.Errorf("run-list file %s has wrong format (want: {ID}%s{DirectiveName})",
			n,
			RunListFileSeparator,
		)
	}

	var err error
	f.id, err = strconv.Atoi(pp[0])
	if err != nil {
		return fmt.Errorf("run-list file %s {ID} is not numeric: %s: %w",
			n,
			pp[0],
			err,
		)
	}

	f.directive = pp[1]
	return nil
}

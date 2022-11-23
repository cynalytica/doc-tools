package flags

import (
	"os"
	"strings"
)

type StringArrayFile struct {
	isSet  bool
	values []string
}

func (s *StringArrayFile) Set(filePath string) error {
	if filePath != "" {
		s.isSet = true
		if data, err := os.ReadFile(filePath); err == nil {
			s.values = strings.Split(string(data), "\n")
		} else {
			return err
		}
	}
	return nil
}
func (s StringArrayFile) IsSet() bool {
	return s.isSet
}
func (s StringArrayFile) Values() []string {
	return s.values
}
func (s StringArrayFile) String() string {
	return ""
}

package version

import "fmt"

var (
	Version      = "0.0.0"
	GitCommit    = "undef"
	GitTreeState = "undef"
	BuildNumber  = "undef"
	BuildDate    = "undef"
	PackageName  = "undef"
)

func String() string {
	return fmt.Sprintf("version %v (GIT commit '%v', treestate '%v' | BUILD number '%v', date '%v' | PACKAGE name '%v')\n", Version, GitCommit, GitTreeState, BuildNumber, BuildDate, PackageName)
}

func Map() map[string]string {
	return map[string]string{
		"version":      Version,
		"gitcommit":    GitCommit,
		"gittreestate": GitTreeState,
		"buildnumber":  BuildNumber,
		"builddate":    BuildDate,
		"packagename":  PackageName,
	}
}

func MapInf() map[string]interface{} {
	m := map[string]interface{}{}

	for k, v := range Map() {
		m[k] = v
	}

	return m
}

func Slice() []string {
	s := []string{}

	for k, v := range Map() {
		s = append(s, k, v)
	}

	return s
}

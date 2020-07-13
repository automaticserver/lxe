package cli

import (
	"testing"
)

func Test_confShowCmdRun(t *testing.T) {
	compareGoldenFile(t, confShowCmd, []string{"json"}, "testdata/confShowCmdRun.golden.json")
}

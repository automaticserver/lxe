package cli

import (
	"testing"
)

func Test_cmplBashCmdRunE(t *testing.T) {
	compareGoldenFile(t, cmplBashCmd, nil, "testdata/cmplBashCmdRun.golden.sh")
}

func Test_cmplZshCmdRunE(t *testing.T) {
	compareGoldenFile(t, cmplZshCmd, nil, "testdata/cmplZshCmdRun.golden.zsh")
}

func Test_cmplPwrshCmdRunE(t *testing.T) {
	compareGoldenFile(t, cmplPwrshCmd, nil, "testdata/cmplPwrshCmdRun.golden.ps1")
}

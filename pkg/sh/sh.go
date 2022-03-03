package sh

import (
	"os"
	"os/exec"
	"strings"

	"github.com/sirupsen/logrus"
)

func Run(scmd string) error {
	logrus.Debug("Running: ", scmd)
	cmds := []string{"-c", scmd}

	cmd := exec.Command("bash", cmds...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func RunMany(scmds []string) error {
	return Run(strings.Join(scmds, ";"))
}

func RunIO(scmd string) ([]byte, error) {
	logrus.Debug("Running: ", scmd)
	cmds := []string{"-c", scmd}

	return exec.Command("bash", cmds...).Output()
}

package aoj

import (
	"io"
	"io/ioutil"
	"os/exec"
)

func Check(path string, input, expected io.Reader) (bool, error) {
	be, err := ioutil.ReadAll(expected)
	if err != nil {
		return false, err
	}

	cmd := exec.Command(path)
	cmd.Stdin = input

	bo, err := cmd.Output()
	if err != nil {
		return false, err
	}

	logdbg("checking testcase output")
	return fuzzyCompare(be, bo), nil
}

func (t *Testcase) CheckCase(caseID int, path string) (bool, error) {
	input, err := t.CaseInput(caseID)
	if err != nil {
		return false, err
	}
	defer input.Close()

	expected, err := t.CaseOutput(caseID)
	if err != nil {
		return false, err
	}
	defer expected.Close()

	return Check(path, input, expected)
}

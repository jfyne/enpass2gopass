package enpass2gopass

import (
	"os/exec"
	"strings"
)

// Insert into gopass.
func Insert(path string, data []string) error {
	cmd := exec.Command("gopass", "insert", "-m", path)
	in, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	defer in.Close()
	if err := cmd.Start(); err != nil {
		return err
	}

	input := strings.Join(data, "\n")
	if _, err := in.Write([]byte(input + "\n")); err != nil {
		return err
	}

	return nil
}

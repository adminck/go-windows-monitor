package comm

import (
	"os"
	"os/exec"
	"path/filepath"
)

func GetCurrentProcessPath() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	return path
}

func GetCurrentProcessName() string {
	return filepath.Base(GetCurrentProcessPath())
}



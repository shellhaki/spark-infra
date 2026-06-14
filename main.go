package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		panic("Usage: sudo go run main.go run <command>")
	}

	switch os.Args[1] {
	case "run":
		run()
	case "child":
		child()
	default:
		panic("Unknown command")
	}
}

func run() {
	targetCmd := os.Args[2:]

	cellID := "cell_" + strconv.FormatInt(time.Now().Unix(), 10)
	cellPath := filepath.Join(os.Getenv("HOME"), "spark_cells", cellID)
	must(os.MkdirAll(cellPath, 0755))

	archivePath := filepath.Join(os.Getenv("HOME"), "Downloads", "base_linux_environment.gz")
	if _, err := os.Stat(archivePath); os.IsNotExist(err) {
		archivePath = filepath.Join(os.Getenv("HOME"), "base_linux_environment.gz")
		if _, err := os.Stat(archivePath); os.IsNotExist(err) {
			panic("Could not find base_linux_environment.gz in local Downloads or Home folder")
		}
	}

	untar(archivePath, cellPath)

	must(os.MkdirAll(filepath.Join(cellPath, "proc"), 0755))
	must(os.MkdirAll(filepath.Join(cellPath, "usr/lib/postgresql"), 0755))

	cmd := exec.Command("/proc/self/exe", append([]string{"child", cellPath}, targetCmd...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWNET,
	}

	must(cmd.Start())

	fmt.Printf("\nSUCCESS! Cell spawned completely isolated from host.\n")
	fmt.Printf("..TRACKING PID ON HOST: %d\n", cmd.Process.Pid)
	fmt.Printf("..FILESYSTEM JAIL DIRECTORY: %s\n\n", cellPath)

	must(cmd.Wait())
}

func child() {
	cellPath := os.Args[2]
	targetCmd := os.Args[3]
	targetArgs := os.Args[4:]

	must(syscall.Mount("/usr/lib/postgresql", filepath.Join(cellPath, "usr/lib/postgresql"), "", syscall.MS_BIND, ""))

	must(syscall.Sethostname([]byte("spark-cell")))
	must(syscall.Chroot(cellPath))
	must(os.Chdir("/"))

	must(syscall.Mount("proc", "proc", "proc", 0, ""))

	cmd := exec.Command(targetCmd, targetArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	must(cmd.Run())

	_ = syscall.Unmount("proc", 0)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func untar(archive string, dest string) {
	cmd := exec.Command("tar", "-xzf", archive, "-C", dest)
	must(cmd.Run())
}

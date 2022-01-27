package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

var boldGreen = color.New(color.Bold, color.FgGreen).SprintfFunc()

var (
	fastbootPath = flag.String("f", "fastboot", "Path to the fastboot executable")
	firmwareDir  = flag.String("D", "", "Path to firmware files")
	skipCheck    = flag.Bool("C", false, "Skip file checksum checks")
	dryRun       = flag.Bool("d", false, "Do not run any commands")
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("Error: %s\n", err)
		for {
			err = errors.Unwrap(err)
			if err == nil {
				break
			}

			fmt.Printf("       Caused by: %s\n", err)
		}
	}
}

func run() error {
	flag.Parse()

	path := flag.Arg(0)
	if path == "" {
		return errors.New("path to flash file not provided")
	}

	dir := *firmwareDir
	if dir == "" {
		dir = filepath.Dir(path)
	}

	file, err := os.Open(path)
	if err != nil {
		return err
	}

	var f Flashfile
	if err := xml.NewDecoder(file).Decode(&f); err != nil {
		return err
	}

	type Command struct {
		Step    *Step
		Command *exec.Cmd
	}

	var cmds []Command
	for _, step := range f.Steps {
		var args []string
		if step.Operation != "" {
			args = append(args, step.Operation)
		}
		if step.Var != "" {
			args = append(args, step.Var)
		}
		if step.Partition != "" {
			args = append(args, step.Partition)
		}
		if step.Filename != "" {
			if !*skipCheck {
				if step.MD5 != "" {
					fmt.Printf(boldGreen("%15s", "Verifying")+" %s\n", step.Filename)

					hash := md5.New()
					file, err := os.Open(filepath.Join(dir, step.Filename))
					if err != nil {
						return err
					}

					if _, err := io.Copy(hash, file); err != nil {
						return err
					}

					md5sum := hex.EncodeToString(hash.Sum(nil))
					if md5sum != step.MD5 {
						return errors.New("checksum failure, file may be corrupted")
					}
				}
			}

			args = append(args, step.Filename)
		}

		cmd := exec.Command(*fastbootPath, args...)
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmds = append(cmds, Command{step, cmd})
	}

	for _, cmd := range cmds {
		var action string
		var args string
		switch cmd.Step.Operation {
		case "flash":
			action = " Flashing"
			args = fmt.Sprintf("%-10s (%s)", cmd.Step.Partition, cmd.Step.Filename)
		case "erase":
			action = "  Erasing"
			args = cmd.Step.Partition
		default:
			action = "  Running"
			args = strings.Join(cmd.Command.Args[1:], " ")
		}

		fmt.Printf(boldGreen("%15s", action)+" %s\n", args)
		if !*dryRun {
			if err := cmd.Command.Run(); err != nil {
				return err
			}
		}
	}

	return nil
}

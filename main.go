package main

import (
	"crypto/md5"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

var boldGreen = color.New(color.Bold, color.FgGreen).SprintFunc()

func main() {
	if err := run(os.Args); err != nil {
		fmt.Printf("Error: %s\n", err)
		for {
			err = errors.Unwrap(err)
			if err != nil {
				fmt.Printf("Caused by: %s\n", err)
				continue
			}

			break
		}
	}
}

func run(args []string) error {
	path := args[1]
	dir := filepath.Dir(path)

	file, err := os.Open(path)
	if err != nil {
		return err
	}

	var f Flashfile
	if err := xml.NewDecoder(file).Decode(&f); err != nil {
		return err
	}

	var cmds []*exec.Cmd
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
			if step.MD5 != "" {
				fmt.Printf("    %s %s\n", boldGreen("Verifying"), step.Filename)

				hash := md5.New()
				file, err := os.Open(filepath.Join(dir, step.Filename))
				if err != nil {
					return err
				}

				if _, err := io.Copy(hash, file); err != nil {
					return err
				}

				md5sum := fmt.Sprintf("%x", hash.Sum(nil))
				if md5sum != step.MD5 {
					return errors.New("checksum failure, file may be corrupted")
				}
			}

			args = append(args, step.Filename)
		}

		cmd := exec.Command("fastboot", args...)
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmds = append(cmds, cmd)
	}

	for _, cmd := range cmds {
		fmt.Printf("      %s %s\n", boldGreen("Running"), strings.Join(cmd.Args[1:], " "))
		if err := cmd.Run(); err != nil {
			return err
		}
	}

	return nil
}

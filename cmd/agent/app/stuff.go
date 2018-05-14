// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2018 Datadog, Inc.

// +build cpython

package app

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/DataDog/datadog-agent/pkg/util/executable"
	"github.com/spf13/cobra"
)

const (
	constraintsFile = "agent_requirements.txt"
	tufConfigFile   = "public-tuf-config.json"
	tufPyPiServer   = "https://integrationsproxy.azurewebsites.net/simple/"
)

var (
	withTuf bool
)

func init() {
	AgentCmd.AddCommand(stuffCmd)
	stuffCmd.AddCommand(installCmd)
	stuffCmd.AddCommand(removeCmd)
	stuffCmd.AddCommand(searchCmd)
	stuffCmd.Flags().BoolVarP(&withTuf, "tuf", "t", true, "use TUF repo")
}

var stuffCmd = &cobra.Command{
	Use:   "stuff [command]",
	Short: "Datadog integration/package manager",
	Long:  ``,
}

var installCmd = &cobra.Command{
	Use:   "install [package]",
	Short: "Install Datadog integration/extra packages",
	Long:  ``,
	RunE:  installStuff,
}

var removeCmd = &cobra.Command{
	Use:   "remove [package]",
	Short: "Remove Datadog integration/extra packages",
	Long:  ``,
	RunE:  removeStuff,
}

var searchCmd = &cobra.Command{
	Use:   "search [package]",
	Short: "Search Datadog integration/extra packages",
	Long:  ``,
	RunE:  searchStuff,
}

func getInstrumentedPipPath() (string, error) {
	here, _ := executable.Folder()
	pipPath := filepath.Join(here, relPipPath)

	if _, err := os.Stat(pipPath); err != nil {
		if os.IsNotExist(err) {
			return pipPath, errors.New("unable to find pip executable")
		}
	}

	return pipPath, nil
}

func getConstraintsFilePath() (string, error) {
	here, _ := executable.Folder()
	cPath := filepath.Join(here, relConstraintsPath)

	if _, err := os.Stat(cPath); err != nil {
		if os.IsNotExist(err) {
			return cPath, err
		}
	}

	return cPath, nil
}

func getTUFConfigFilePath() (string, error) {
	here, _ := executable.Folder()
	tPath := filepath.Join(here, relTufConfigFilePath)

	if _, err := os.Stat(tPath); err != nil {
		if os.IsNotExist(err) {
			return tPath, err
		}
	}

	return tPath, nil
}

func stuff(args []string) error {
	pipPath, err := getInstrumentedPipPath()
	if err != nil {
		return err
	}
	tufPath, err := getTUFConfigFilePath()
	if err != nil && withTuf {
		return err
	}

	stuffCmd := exec.Command(pipPath, args...)

	var stdout, stderr bytes.Buffer
	stuffCmd.Stdout = &stdout
	stuffCmd.Stderr = &stderr
	if withTuf {
		stuffCmd.Env = append(os.Environ(),
			fmt.Sprintf("TUF_CONFIG_FILE=%s", tufPath),
		)
	}

	err = stuffCmd.Run()
	if err != nil {
		fmt.Printf("error running command: %v", stderr.String())
	} else {
		fmt.Printf("%v", stdout.String())
	}

	return err
}

func installStuff(cmd *cobra.Command, args []string) error {
	constraintsPath, err := getConstraintsFilePath()
	if err != nil {
		return err
	}

	stuffArgs := []string{
		"install",
		"-c", constraintsPath,
	}
	if withTuf {
		stuffArgs = append(stuffArgs, "--index-url", tufPyPiServer)
	}
	stuffArgs = append(stuffArgs, args...)

	return stuff(stuffArgs)
}

func removeStuff(cmd *cobra.Command, args []string) error {
	stuffArgs := []string{
		"uninstall",
	}
	stuffArgs = append(stuffArgs, args...)

	return stuff(stuffArgs)
}

func searchStuff(cmd *cobra.Command, args []string) error {

	stuffArgs := []string{
		"search",
	}
	if withTuf {
		stuffArgs = append(stuffArgs, "--index-url", tufPyPiServer)
	}
	stuffArgs = append(stuffArgs, args...)

	return stuff(stuffArgs)
}

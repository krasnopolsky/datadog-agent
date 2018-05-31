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
	pyPiServer      = "https://pypi.org/simple/"
)

var (
	withTuf   bool
	nativePkg bool
)

func init() {
	AgentCmd.AddCommand(tufCmd)
	tufCmd.AddCommand(installCmd)
	tufCmd.AddCommand(removeCmd)
	tufCmd.AddCommand(searchCmd)
	tufCmd.Flags().BoolVarP(&withTuf, "tuf", "t", true, "use TUF repo")
	tufCmd.Flags().BoolVarP(&nativePkg, "pip-package", "p", false, "providing native pip package name")
}

var tufCmd = &cobra.Command{
	Use:   "integration [command]",
	Short: "Datadog integration/package manager",
	Long:  ``,
}

var installCmd = &cobra.Command{
	Use:   "install [package]",
	Short: "Install Datadog integration/extra packages",
	Args:  cobra.ArbitraryArgs,
	Long:  ``,
	RunE:  installTuf,
}

var removeCmd = &cobra.Command{
	Use:   "remove [package]",
	Short: "Remove Datadog integration/extra packages",
	Args:  cobra.ArbitraryArgs,
	Long:  ``,
	RunE:  removeTuf,
}

var searchCmd = &cobra.Command{
	Use:   "search [package]",
	Short: "Search Datadog integration/extra packages",
	Args:  cobra.ArbitraryArgs,
	Long:  ``,
	RunE:  searchTuf,
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

func tuf(args []string) error {
	pipPath, err := getInstrumentedPipPath()
	if err != nil {
		return err
	}
	tufPath, err := getTUFConfigFilePath()
	if err != nil && withTuf {
		return err
	}

	tufCmd := exec.Command(pipPath, args...)

	var stdout, stderr bytes.Buffer
	tufCmd.Stdout = &stdout
	tufCmd.Stderr = &stderr
	if withTuf {
		tufCmd.Env = append(os.Environ(),
			fmt.Sprintf("TUF_CONFIG_FILE=%s", tufPath),
		)
	}

	err = tufCmd.Run()
	if err != nil {
		fmt.Printf("error running command: %v", stderr.String())
	} else {
		fmt.Printf("%v", stdout.String())
	}

	return err
}

func installTuf(cmd *cobra.Command, args []string) error {
	constraintsPath, err := getConstraintsFilePath()
	if err != nil {
		return err
	}

	tufArgs := []string{
		"install",
		"-c", constraintsPath,
	}

	tufArgs = append(tufArgs, args...)
	if withTuf {
		tufArgs = append(tufArgs, "--index-url", tufPyPiServer)
		tufArgs = append(tufArgs, "--extra-index-url", pyPiServer)
		tufArgs = append(tufArgs, "--disable-pip-version-check")
	}

	return tuf(tufArgs)
}

func removeTuf(cmd *cobra.Command, args []string) error {
	tufArgs := []string{
		"uninstall",
	}
	tufArgs = append(tufArgs, args...)
	tufArgs = append(tufArgs, "-y")

	return tuf(tufArgs)
}

func searchTuf(cmd *cobra.Command, args []string) error {

	tufArgs := []string{
		"search",
	}
	tufArgs = append(tufArgs, args...)
	if withTuf {
		tufArgs = append(tufArgs, "--index-url", tufPyPiServer)
		tufArgs = append(tufArgs, "--extra-index-url", pyPiServer)
		tufArgs = append(tufArgs, "--disable-pip-version-check")
	}

	return tuf(tufArgs)
}

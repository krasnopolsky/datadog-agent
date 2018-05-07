// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2018 Datadog, Inc.

// +build cpython

package app

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

const (
	constraints = "agent_requirements.txt"
)

func init() {
	AgentCmd.AddCommand(stuffCmd)
	stuffCmd.AddCommand(installCmd)
	stuffCmd.AddCommand(removeCmd)
	stuffCmd.AddCommand(searchCmd)
	stuffCmd.Flags().BoolVarP(&jsonStatus, "verbose", "v", false, "verbose output")
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
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

var searchCmd = &cobra.Command{
	Use:   "search [package]",
	Short: "Search Datadog integration/extra packages",
	Long:  ``,
	RunE:  searchStuff,
}

func stuff(args []string) error {
	pipPath, err := getInstrumentedPipPath()
	if err != nil {
		return err
	}

	stuffCmd := exec.Command(pipPath, args...)

	var stdout, stderr bytes.Buffer
	stuffCmd.Stdout = &stdout
	stuffCmd.Stderr = &stderr

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
	stuffArgs = append(stuffArgs, args...)

	return stuff(stuffArgs)
}

func removeStuff(cmd *cobra.Command, args []string) error {
	return nil
}

func searchStuff(cmd *cobra.Command, args []string) error {
	stuffArgs := []string{
		"search",
	}

	stuffArgs = append(stuffArgs, args...)
	return stuff(stuffArgs)
}

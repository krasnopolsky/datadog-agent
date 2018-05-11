// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2018 Datadog, Inc.

// +build windows
// +build cpython

package app

import (
	"os"
	"path/filepath"

	"github.com/DataDog/datadog-agent/pkg/util/executable"
)

const (
	pip = "pip.exe"
)

func getInstrumentedPipPath() (string, error) {
	here, _ := executable.Folder()
	pipPath := filepath.Join(here, "Scripts", pip)

	if _, err := os.Stat(pipPath); err != nil {
		if os.IsNotExist(err) {
			return pipPath, errors.New("unable to find pip executable")
		}
	}

	return pipPath, nil
}

func getConstraintsFilePath() (string, error) {
	here, _ := executable.Folder()
	cPath := filepath.Join(here, "..", constraintsFile)

	if _, err := os.Stat(cPath); err != nil {
		if os.IsNotExist(err) {
			return nil, err
		}
	}

	return cPath, nil
}

func getTUFConfigFilePath() (string, error) {
	here, _ := executable.Folder()
	tPath := filepath.Join(here, "..", "..", tufConfigFile)

	if _, err := os.Stat(tPath); err != nil {
		if os.IsNotExist(err) {
			return nil, err
		}
	}

	return tPath, nil
}

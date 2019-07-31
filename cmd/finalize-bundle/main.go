package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/projectriff/k8s-manifest-scanner/pkg/bundle"
)

func main() {
	bundlePath, manifestPath, err := verifyCommandLineArgs(os.Args)
	if err != nil {
		fmt.Printf("error validating arguments %v\n", err)
		os.Exit(1)
	}

	err = bundle.Finalize(bundlePath, manifestPath)
	if err != nil {
		fmt.Printf("error updating bundle: %v\n", err)
		os.Exit(1)
	}
}

func verifyCommandLineArgs(args []string) (bundleTemplatePath, manifestPath string, err error) {
	if len(args) == 1 {
		return "duffle.json", "./cnab/app/kab/manifest.yaml", nil
	}
	if len(args) != 3 {
		return "", "", errors.New("usage: ./finalize-bundle <path/to/duffle.json> </path/to/kab-manifest.yaml>")
	}
	return args[1], args[2], nil
}

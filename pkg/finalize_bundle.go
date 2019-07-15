package cnab_riff

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/deislabs/cnab-go/bundle"
	"github.com/deislabs/duffle/pkg/duffle/manifest"
	"github.com/ghodss/yaml"
	"github.com/pivotal/go-ape/pkg/furl"
	"github.com/pivotal/image-relocation/pkg/image"
	"github.com/projectriff/cnab-k8s-installer-base/pkg/apis/kab/v1alpha1"
	"github.com/projectriff/k8s-manifest-scanner/pkg/scan"
)

// this performs two tasks:
// 1. inlines the content of the resource url into the bundle
// 2. adds images to duffle.json by scanning the resource content
func FinalizeBundle(bundlePath, kabManifestPath string) error {
	mfst := &manifest.Manifest{}
	err := unmarshallFile(bundlePath, mfst)
	if err != nil {
		return err
	}

	err = InlineContentInKabManifest(kabManifestPath)
	if err != nil {
		return err
	}

	images, err := GetImagesFromKabManifest(kabManifestPath)
	if err != nil {
		return err
	}

	mfst.Images = map[string]bundle.Image{}
	for _, img := range images {
		name, err := image.NewName(img)
		if err != nil {
			fmt.Printf("err %v\n", err)
		}
		n := strings.ReplaceAll(name.String(), "/", "_")
		bunImg := bundle.Image{}
		bunImg.Image = name.String()
		bunImg.Digest = name.Digest().String()
		mfst.Images[n] = bunImg
	}

	err = marshallJsonFile(bundlePath, mfst)
	return err

}

func unmarshallFile(path string, str interface{}) error {
	mfstBytes, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("error reading file %s: %v", path, err)
		return err
	}
	err = yaml.Unmarshal(mfstBytes, str)
	if err != nil {
		fmt.Printf("error unmarshalling file %s: %v", path, err)
		return err
	}
	return nil
}

func marshallJsonFile(path string, str interface{}) error {
	mfstBytes, err := json.MarshalIndent(str, "", "    ")
	if err != nil {
		return err
	}
	err = writeFile(path, str, mfstBytes)
	return err
}

func marshallYamlFile(path string, str interface{}) error {
	mfstBytes, err := yaml.Marshal(str)
	if err != nil {
		return err
	}
	err = writeFile(path, str, mfstBytes)
	return err
}

func writeFile(path string, str interface{}, content []byte) error {
	err := ioutil.WriteFile(path, content, 0644)
	if err != nil {
		return err
	}
	fmt.Printf("wrote file %s\n", path)
	return nil
}

func GetImagesFromKabManifest(kabManifestPath string) ([]string, error) {
	kabMfst := &v1alpha1.Manifest{}
	err := unmarshallFile(kabManifestPath, kabMfst)
	if err != nil {
		return nil, err
	}

	images := []string{}

	err = kabMfst.VisitResources(func(res v1alpha1.KabResource) error {
		fmt.Fprintf(os.Stderr, "Scanning %s\n", res.Name)

		var err error
		var imgs []string
		if res.Content != "" {
			imgs, err = scan.ListSortedImagesFromContent([]byte(res.Content))
		} else {
			imgs, err = scan.ListSortedImagesFromKubernetesManifest(res.Path, "")
		}
		if err != nil {
			return err
		}

		images = append(images, imgs...)

		return nil
	})
	return images, nil
}

func InlineContentInKabManifest(kabManifestPath string) error {
	kabMfst := &v1alpha1.Manifest{}
	err := unmarshallFile(kabManifestPath, kabMfst)
	if err != nil {
		return err
	}

	err = kabMfst.PatchResourceContent(func(res *v1alpha1.KabResource) (string, error) {
		if res.Content != "" {
			return "", errors.New(fmt.Sprintf("content not empty for resource: %s", res.Name))
		}
		contentBytes, err := furl.Read(res.Path, "")
		if err != nil {
			return "", err
		}
		return string(contentBytes), nil
	})
	if err != nil {
		return err
	}

	err = marshallYamlFile(kabManifestPath, kabMfst)
	return err
}

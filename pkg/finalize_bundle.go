package cnab_riff

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/deislabs/cnab-go/bundle"
	"github.com/deislabs/duffle/pkg/duffle/manifest"
	"github.com/ghodss/yaml"
	"github.com/pivotal/go-ape/pkg/furl"
	"github.com/pivotal/image-relocation/pkg/image"
	"github.com/pivotal/image-relocation/pkg/registry"
	"github.com/projectriff/cnab-k8s-installer-base/pkg/apis/kab/v1alpha1"
	"github.com/projectriff/k8s-manifest-scanner/pkg/scan"
)

// this performs following tasks:
// 1. inlines the content of the resource url into the bundle
// 2. adds images to duffle.json by scanning the resource content
// 3. computes digests for images
// 4. replaces image references in kab manifest with digested references
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
	r := registry.NewRegistryClient()
	replacements := []string{}

	for _, img := range images {
		name, err := image.NewName(img)
		if err != nil {
			fmt.Printf("err %v\n", err)
		}
		n := strings.ReplaceAll(name.String(), "/", "_")
		bunImg := bundle.Image{}
		d, err := r.Digest(name)
		if err != nil {
			return err
		}
		bunImg.Digest = d.String()
		nameWithDigest, err := name.WithDigest(d)
		if err != nil {
			return err
		}
		bunImg.Image = nameWithDigest.String()
		mfst.Images[n] = bunImg

		replacements = append(replacements, img, nameWithDigest.String())
	}

	err = marshallJsonFile(bundlePath, mfst)
	if err != nil {
		return err
	}

	err = ReplaceInKabManifest(kabManifestPath, *strings.NewReplacer(replacements...))

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

func ReplaceInKabManifest(kabManifestPath string, replacer strings.Replacer) error {
	kabMfst := &v1alpha1.Manifest{}
	err := unmarshallFile(kabManifestPath, kabMfst)
	if err != nil {
		return err
	}

	err = kabMfst.PatchResourceContent(func(res *v1alpha1.KabResource) (string, error) {
		return replacer.Replace(res.Content), nil
	})
	if err != nil {
		return err
	}

	err = marshallYamlFile(kabManifestPath, kabMfst)
	return err
}

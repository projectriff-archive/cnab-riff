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
	err := unmarshalFile(bundlePath, mfst)
	if err != nil {
		return err
	}

	kabMfst := &v1alpha1.Manifest{}
	err = unmarshalFile(kabManifestPath, kabMfst)
	if err != nil {
		return err
	}

	err = kabMfst.InlineContent()
	if err != nil {
		return err
	}

	images, err := getImagesFromKabManifest(kabMfst)
	if err != nil {
		return err
	}

	mfst.Images = map[string]bundle.Image{}
	registryClient := registry.NewRegistryClient()
	replacements := []string{}

	for _, img := range images {
		name, err := image.NewName(img)
		if err != nil {
			fmt.Printf("err %v\n", err)
		}
		bundleImageKey := strings.ReplaceAll(name.String(), "/", "_")
		bundleImage := bundle.Image{}

		digest, err := registryClient.Digest(name)
		if err != nil {
			return err
		}
		bundleImage.Digest = digest.String()
		nameWithDigest, err := name.WithDigest(digest)
		if err != nil {
			return err
		}
		bundleImage.Image = nameWithDigest.String()

		mfst.Images[bundleImageKey] = bundleImage

		replacements = append(replacements, img, nameWithDigest.String())
	}

	err = marshalJsonFile(bundlePath, mfst)
	if err != nil {
		return err
	}

	err = replaceInKabManifest(kabMfst, *strings.NewReplacer(replacements...))
	if err != nil {
		return err
	}

	err = marshalYamlFile(kabManifestPath, kabMfst)
	return err

}

func unmarshalFile(path string, manifest interface{}) error {
	mfstBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("error reading file %s: %v", path, err)
	}
	err = yaml.Unmarshal(mfstBytes, manifest)
	if err != nil {
		return fmt.Errorf("error unmarshalling file %s: %v", path, err)
	}
	return nil
}

func marshalJsonFile(path string, manifest interface{}) error {
	mfstBytes, err := json.MarshalIndent(manifest, "", "    ")
	if err != nil {
		return err
	}
	return writeFile(path, mfstBytes)
}

func marshalYamlFile(path string, str interface{}) error {
	mfstBytes, err := yaml.Marshal(str)
	if err != nil {
		return err
	}
	return writeFile(path, mfstBytes)
}

func writeFile(path string, content []byte) error {
	err := ioutil.WriteFile(path, content, 0644)
	if err != nil {
		return err
	}
	fmt.Printf("wrote file %s\n", path)
	return nil
}

func getImagesFromKabManifest(kabMfst *v1alpha1.Manifest) ([]string, error) {

	images := []string{}

	err := kabMfst.VisitResources(func(res v1alpha1.KabResource) error {
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
	return images, err
}

func replaceInKabManifest(kabMfst *v1alpha1.Manifest, replacer strings.Replacer) error {

	err := kabMfst.PatchResourceContent(func(res *v1alpha1.KabResource) (string, error) {
		return replacer.Replace(res.Content), nil
	})
	if err != nil {
		return err
	}
	return err
}

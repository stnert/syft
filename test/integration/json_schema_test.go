package integration

import (
	"bytes"
	"fmt"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/anchore/syft/internal/presenter/packages"

	"github.com/anchore/syft/internal"

	"github.com/anchore/syft/syft/distro"

	"github.com/xeipuuv/gojsonschema"
)

const jsonSchemaPath = "schema/json"

func repoRoot(t *testing.T) string {
	t.Helper()
	repoRoot, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		t.Fatalf("unable to find repo root dir: %+v", err)
	}
	absRepoRoot, err := filepath.Abs(strings.TrimSpace(string(repoRoot)))
	if err != nil {
		t.Fatal("unable to get abs path to repo root:", err)
	}
	return absRepoRoot
}

func validateAgainstV1Schema(t *testing.T, json string) {
	fullSchemaPath := path.Join(repoRoot(t), jsonSchemaPath, fmt.Sprintf("schema-%s.json", internal.JSONSchemaVersion))
	schemaLoader := gojsonschema.NewReferenceLoader(fmt.Sprintf("file://%s", fullSchemaPath))
	documentLoader := gojsonschema.NewStringLoader(json)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		t.Fatal("unable to validate json schema:", err.Error())
	}

	if !result.Valid() {
		t.Errorf("failed json schema validation:")
		t.Errorf("JSON:\n%s\n", json)
		for _, desc := range result.Errors() {
			t.Errorf("  - %s\n", desc)
		}
	}
}

func TestJsonSchemaImg(t *testing.T) {

	catalog, _, src := catalogFixtureImage(t, "image-pkg-coverage")

	output := bytes.NewBufferString("")

	d, err := distro.NewDistro(distro.CentOS, "5", "rhel fedora")
	if err != nil {
		t.Fatalf("bad distro: %+v", err)
	}

	p := packages.Presenter(packages.JSONPresenterOption, src.Metadata, catalog, &d)
	if p == nil {
		t.Fatal("unable to get presenter")
	}

	err = p.Present(output)
	if err != nil {
		t.Fatalf("unable to present: %+v", err)
	}

	validateAgainstV1Schema(t, output.String())

}

func TestJsonSchemaDirs(t *testing.T) {
	catalog, _, src := catalogDirectory(t, "test-fixtures/image-pkg-coverage")

	output := bytes.NewBufferString("")

	d, err := distro.NewDistro(distro.CentOS, "5", "rhel fedora")
	if err != nil {
		t.Fatalf("bad distro: %+v", err)
	}

	p := packages.Presenter(packages.JSONPresenterOption, src.Metadata, catalog, &d)
	if p == nil {
		t.Fatal("unable to get presenter")
	}

	err = p.Present(output)
	if err != nil {
		t.Fatalf("unable to present: %+v", err)
	}

	validateAgainstV1Schema(t, output.String())
}

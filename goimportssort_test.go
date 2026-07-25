package main

import (
	"go/build"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessFile(t *testing.T) {
	asserts := assert.New(t)
	*localPrefix = "github.com/bonsai-oss/goimportssort"
	reader := strings.NewReader(`package main

// builtin
// external
// local
import (
	"fmt"
	"log"
	
	APA "bitbucket.org/example/package/name"
	APZ "bitbucket.org/example/package/name"
	"bitbucket.org/example/package/name2"
	"bitbucket.org/example/package/name3" // foopsie
	"bitbucket.org/example/package/name4"
	
	"github.com/bonsai-oss/goimportssort/package1"
	// a
	"github.com/bonsai-oss/goimportssort/package2"
	
	/*
		mijn comment
	*/
	"net/http/httptest"
	"database/sql/driver"
)
// klaslkasdko

func main() {
	fmt.Println("Hello!")
}`)
	want := `package main

import (
	"database/sql/driver"
	"fmt"
	"log"
	"net/http/httptest"

	APA "bitbucket.org/example/package/name"
	APZ "bitbucket.org/example/package/name"
	"bitbucket.org/example/package/name2"
	"bitbucket.org/example/package/name3"
	"bitbucket.org/example/package/name4"

	"github.com/bonsai-oss/goimportssort/package1"
	"github.com/bonsai-oss/goimportssort/package2"
)

func main() {
	fmt.Println("Hello!")
}
`

	output, err := processFile("", reader, os.Stdout)
	asserts.NotEqual(nil, output)
	asserts.Equal(nil, err)
	asserts.Equal(want, string(output))
}

func TestProcessFile_Order(t *testing.T) {
	asserts := assert.New(t)
	*localPrefix = "github.com/bonsai-oss/goimportssort"

	reader := strings.NewReader(
		`package main

import "fmt"

import "github.com/exampleUser/examplePackage"

import "github.com/bonsai-oss/goimportssort/package1"


func main() {
	fmt.Println("Hello!")
}`)
	*order = "lei"
	output, err := processFile("", reader, os.Stdout)
	*order = DefaultOrder // reset order for other tests
	asserts.NotEqual(nil, output)
	asserts.Equal(nil, err)
	asserts.Equal(
		`package main

import (
	"github.com/bonsai-oss/goimportssort/package1"

	"github.com/exampleUser/examplePackage"

	"fmt"
)

func main() {
	fmt.Println("Hello!")
}
`, string(output))
}

func TestProcessFile_SingleImport(t *testing.T) {
	asserts := assert.New(t)
	*localPrefix = "github.com/bonsai-oss/goimportssort"

	reader := strings.NewReader(
		`package main


import "github.com/bonsai-oss/goimportssort/package1"


func main() {
	fmt.Println("Hello!")
}`)
	output, err := processFile("", reader, os.Stdout)
	asserts.NotEqual(nil, output)
	asserts.Equal(nil, err)
	asserts.Equal(
		`package main

import (
	"github.com/bonsai-oss/goimportssort/package1"
)

func main() {
	fmt.Println("Hello!")
}
`, string(output))
}

func TestProcessFile_GenericsSupport(t *testing.T) {
	asserts := assert.New(t)
	*localPrefix = "github.com/bonsai-oss/goimportssort"

	reader := strings.NewReader(
		`package main


import "github.com/bonsai-oss/goimportssort/package1"

func filter[T any](ss []T, test func(T) bool) (ret []T) {
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
}

func main() {
	fmt.Println("Hello!")
}`)
	output, err := processFile("", reader, os.Stdout)
	asserts.NotEqual(nil, output)
	asserts.Equal(nil, err)
	asserts.Equal(
		`package main

import (
	"github.com/bonsai-oss/goimportssort/package1"
)

func filter[T any](ss []T, test func(T) bool) (ret []T) {
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
}

func main() {
	fmt.Println("Hello!")
}
`, string(output))
}

func TestProcessFile_EmptyImport(t *testing.T) {
	asserts := assert.New(t)
	*localPrefix = "github.com/bonsai-oss/goimportssort"

	reader := strings.NewReader(`package main

func main() {
	fmt.Println("Hello!")
}`)
	output, err := processFile("", reader, os.Stdout)
	asserts.NotEqual(nil, output)
	asserts.Equal(nil, err)
	asserts.Equal(`package main

func main() {
	fmt.Println("Hello!")
}`, string(output))
}

func TestProcessFile_ReadMeExample(t *testing.T) {
	asserts := assert.New(t)
	*localPrefix = "github.com/bonsai-oss/goimportssort"

	reader := strings.NewReader(`package main

import (
	"fmt"
	"log"
	APZ "bitbucket.org/example/package/name"
	APA "bitbucket.org/example/package/name"
	"github.com/bonsai-oss/goimportssort/package2"
	"github.com/bonsai-oss/goimportssort/package1"
)
import (
	"net/http/httptest"
)

import "bitbucket.org/example/package/name2"
import "bitbucket.org/example/package/name3"
import "bitbucket.org/example/package/name4"`)
	output, err := processFile("", reader, os.Stdout)
	asserts.NotEqual(nil, output)
	asserts.Equal(nil, err)
	asserts.Equal(`package main

import (
	"fmt"
	"log"
	"net/http/httptest"

	APA "bitbucket.org/example/package/name"
	APZ "bitbucket.org/example/package/name"
	"bitbucket.org/example/package/name2"
	"bitbucket.org/example/package/name3"
	"bitbucket.org/example/package/name4"

	"github.com/bonsai-oss/goimportssort/package1"
	"github.com/bonsai-oss/goimportssort/package2"
)
`, string(output))
}

func TestProcessFile_WronglyFormattedGo(t *testing.T) {
	asserts := assert.New(t)
	*localPrefix = "github.com/bonsai-oss/goimportssort"

	reader := strings.NewReader(
		`package main
import "github.com/bonsai-oss/goimportssort/package1"


func main() {
	fmt.Println("Hello!")
}`)
	output, err := processFile("", reader, os.Stdout)
	asserts.NotEqual(nil, output)
	asserts.Equal(nil, err)
	asserts.Equal(
		`package main

import (
	"github.com/bonsai-oss/goimportssort/package1"
)

func main() {
	fmt.Println("Hello!")
}
`, string(output))
}

func TestIsStandardPackage(t *testing.T) {
	for _, testCase := range []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "single element standard package",
			input:    "fmt",
			expected: true,
		},
		{
			name:     "nested standard package",
			input:    "net/http",
			expected: true,
		},
		{
			name:     "deeply nested standard package",
			input:    "database/sql/driver",
			expected: true,
		},
		{
			name:     "third party package",
			input:    "github.com/bonsai-oss/goimportssort",
			expected: false,
		},
		{
			name:     "golang.org x package is not standard",
			input:    "golang.org/x/tools/go/packages",
			expected: false,
		},
		{
			name:     "bitbucket package",
			input:    "bitbucket.org/example/package/name",
			expected: false,
		},
		{
			name:     "dotless non-standard module is not standard",
			input:    "mycustomsoftware/lib/foo",
			expected: false,
		},
		{
			name:     "dotless single element non-standard module is not standard",
			input:    "mymonorepo",
			expected: false,
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			asserts := assert.New(t)

			actual := isStandardPackage(testCase.input)

			asserts.Equal(testCase.expected, actual)
		})
	}
}

func TestIsStandardPackage_FallsBackToEmbeddedListWithoutGoroot(t *testing.T) {
	asserts := assert.New(t)

	original := build.Default.GOROOT
	build.Default.GOROOT = filepath.Join(t.TempDir(), "nonexistent")
	defer func() { build.Default.GOROOT = original }()

	asserts.True(isStandardPackage("fmt"))
	asserts.True(isStandardPackage("net/http"))
	asserts.True(isStandardPackage("unsafe"))
	asserts.False(isStandardPackage("mycustomsoftware/lib/foo"))
	asserts.False(isStandardPackage("github.com/some/external"))
}

func TestStandardPackagesGoVersion(t *testing.T) {
	asserts := assert.New(t)

	asserts.True(strings.HasPrefix(standardPackagesGoVersion, "go"),
		"embedded list must record the Go version it was generated from, got %q", standardPackagesGoVersion)
}

func TestProcessFile_PreservesGoDirectiveAboveImports(t *testing.T) {
	asserts := assert.New(t)
	*localPrefix = "github.com/bonsai-oss/goimportssort"

	reader := strings.NewReader(`package main

//go:generate go run gen_stdpackages.go

import (
	"os"
	"fmt"
)

func main() { fmt.Fprintln(os.Stdout) }
`)
	output, err := processFile("", reader, os.Stdout)
	asserts.Equal(nil, err)

	out := string(output)
	asserts.Contains(out, "//go:generate go run gen_stdpackages.go")
	asserts.Less(strings.Index(out, "//go:generate"), strings.Index(out, "import ("),
		"directive must stay above the import block")
}

func TestGetModuleName(t *testing.T) {
	asserts := assert.New(t)

	name := getModuleName()

	asserts.Equal("github.com/bonsai-oss/goimportssort", name)
}

func TestSortString(t *testing.T) {
	for _, testCase := range []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty",
			input:    "",
			expected: "",
		},
		{
			name:     "different chars",
			input:    "cab",
			expected: "abc",
		},
		{
			name:     "identical chars",
			input:    "caba",
			expected: "aabc",
		},
		{
			name:     "chars with numbers and symbols",
			input:    "caba!@#$%^&*()_+1234567890",
			expected: "!#$%&()*+0123456789@^_aabc",
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			asserts := assert.New(t)

			actual := sortString(testCase.input)

			asserts.Equal(testCase.expected, actual)
		})
	}
}

func TestIsGoFile(t *testing.T) {
	for _, testCase := range []struct {
		name          string
		inputFilePath string
		expected      bool
	}{
		{
			name:          "go file",
			inputFilePath: "goimportssort.go",
			expected:      true,
		},
		{
			name:          "non go file",
			inputFilePath: ".gitignore",
			expected:      false,
		},
		{
			name:          "directory",
			inputFilePath: ".gitlab",
			expected:      false,
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			asserts := assert.New(t)

			info, _ := os.Stat(testCase.inputFilePath)

			actual := isGoFile(info)

			asserts.Equal(testCase.expected, actual)
		})
	}
}

// TestBinary_ProcessesModuleTargetingNewerGoVersion is an end-to-end regression
// test for the class of error where the checked project targets a Go version
// higher than the toolchain available at runtime. It builds the binary and runs
// it with GOTOOLCHAIN=local against a module declaring go 1.99.0. The previous
// go-command-based std detection failed here with "go.mod requires go >= 1.99.0";
// the binary must now categorise imports purely from the local Go installation.
func TestBinary_ProcessesModuleTargetingNewerGoVersion(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping end-to-end binary test in short mode")
	}
	asserts := assert.New(t)

	binPath := filepath.Join(t.TempDir(), "goimportssort")
	buildCmd := exec.Command("go", "build", "-o", binPath, ".")
	buildCmd.Stderr = os.Stderr
	asserts.NoError(buildCmd.Run(), "building the goimportssort binary")

	moduleDir := t.TempDir()
	asserts.NoError(os.WriteFile(filepath.Join(moduleDir, "go.mod"), []byte("module testmod\n\ngo 1.99.0\n"), 0o644))
	source := `package main

import (
	"testmod/internal/thing"
	"mycustomsoftware/lib/foo"
	"github.com/some/external"
	"fmt"
)

func main() { fmt.Println(external.X, foo.Z, thing.Y) }
`
	asserts.NoError(os.WriteFile(filepath.Join(moduleDir, "main.go"), []byte(source), 0o644))

	run := exec.Command(binPath, "-l", "-local", "testmod", ".")
	run.Dir = moduleDir
	run.Env = append(os.Environ(), "GOTOOLCHAIN=local")
	output, err := run.CombinedOutput()

	asserts.NoError(err, "binary must not depend on the target module's go toolchain: %s", string(output))
	asserts.Contains(string(output), `import (
	"fmt"

	"github.com/some/external"
	"mycustomsoftware/lib/foo"

	"testmod/internal/thing"
)`)
}

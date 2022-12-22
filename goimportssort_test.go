package main

import (
	"os"
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

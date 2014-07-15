package resources_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestTracker(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Resources Suite")
}

func Fixture(filename string) string {
	path := filepath.Join("..", "fixtures", filename)
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return string(contents)
}

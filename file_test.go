package aptly

import (
	"fmt"
	"testing"

	"github.com/kr/pretty"
	"github.com/stretchr/testify/assert"
)

func TestgenerateDir(t *testing.T) {
	client := testClient()
	dir, err := client.File.generateDir("libboost-program-options-dev_1.49.0.1_i386.deb")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "libboost-program-options-dev_1.49.0.1", dir)
}

func TestUploadFile(t *testing.T) {
	client := testClient()
	file, err := client.File.UploadFile("libboost-program-options-dev_1.49.0.1_i386.deb")
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%# v\n", pretty.Formatter(file))
}

func TestListDirectories(t *testing.T) {
	client := testClient()
	dirs, err := client.File.ListDirectories()
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%# v+\n", pretty.Formatter(dirs))
}

func TestListFiles(t *testing.T) {
	client := testClient()
	files, err := client.File.ListFiles("libboost-program-options-dev_1.49.0.1")
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%# v\n", pretty.Formatter(files))
}

func TestDelete(t *testing.T) {
	client := testClient()
	err := client.File.Delete("libboost-program-options-dev_1.49.0.1")
	if err != nil {
		t.Error(err)
	}
}

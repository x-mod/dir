package dir

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDir_Open(t *testing.T) {
	dir := New(
		Root("test"),
		Folder("child"),
		Folder("child1"),
		Folder("child2"),
		Folder("child2/a"),
		Folder("child2/b"),
		Folder("child2", "c", "c", "d"),
		Folder("child3"),
	)
	assert.Nil(t, dir.Open())
	exist, err := dir.DirExists("child")
	assert.Nil(t, err)
	assert.Equal(t, true, exist)
	folders, err := dir.Folders("child2")
	log.Println("folders:", folders)
	assert.Nil(t, err)
	assert.Equal(t, 3, len(folders))

	folders3, err := dir.Folders("child3")
	log.Println("folders:", folders3)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(folders3))

	log.Println("path:", dir.Path("a", "b", "../c/index.md"))

	exist2, err := dir.Exists("child2", "a", "filename")
	assert.Nil(t, err)
	assert.Equal(t, false, exist2)

	assert.NotNil(t, dir.Remove("child2"))
	assert.Nil(t, dir.RemoveAll("child2"))
}

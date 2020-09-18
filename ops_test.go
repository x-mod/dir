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
	)
	assert.Nil(t, dir.Open())
	exist, err := dir.DirExists("child")
	assert.Nil(t, err)
	assert.Equal(t, true, exist)
	folders, err := dir.Folders("child2")
	log.Println("folders:", folders)
	assert.Nil(t, err)
	assert.Equal(t, 3, len(folders))
}

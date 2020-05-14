package compile

import (
	"testing"

	"github.com/matryer/is"
)

func TestReadAllFiles(t *testing.T) {
	is := is.New(t)
	files, err := openAllFiles("../../testdata/badSolFiles", 0)
	if err != nil {
		t.Error(err)
	}
	expected := []string{"burrito", "salad", "taco"}
	for i, f := range files {
		is.Equal(f, expected[i])
	}
}

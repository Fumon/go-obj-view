package obj

import (
	"testing"
	"io/ioutil"
)

func TestParseDump(t *testing.T) {
	file, err := ioutil.ReadFile("monkey.obj")
	if err != nil {
		t.Fatal("Could not open monkey.obj")
	}
	ob, err := Parse(string(file))
	if err != nil {
		t.Fatalf("Parse error:\n%v", err)
	}
	t.Logf("%v\n", ob)
	return
}

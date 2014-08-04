package ahimsa_test

import (
	"testing"

	"github.com/NSkelsey/protocol/ahimsa"
)

func TestFullPath(t *testing.T) {
	cmd, err := ahimsa.GetCreateSql()
	if err != nil {
		t.Fatal(err)
	}

	if len(cmd) < 20 {
		t.Fatalf("Returned create is broken: [%s]", cmd)
	}
}

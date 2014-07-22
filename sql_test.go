package protocol_test

import "testing"
import "github.com/NSkelsey/protocol"

func TestFullPath(t *testing.T) {
	cmd, err := protocol.GetCreateSql()
	if err != nil {
		t.Fatal(err)
	}

	if len(cmd) < 20 {
		t.Fatalf("Returned create is broken: [%s]", cmd)
	}
}

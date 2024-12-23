package downwardsapi_test

import (
	"testing"

	"github.com/tvanriel/cloudsdk/downwardsapi"
)

const s = `
app.kubernetes.io/com-ponent="test"
rack=22

# a comment
test=something
`

func TestMain(t *testing.T) {
	lines, err := downwardsapi.Parse(s)
	if err != nil {
		t.Errorf("error running Parse: %v", err)
	}
	if len(lines) != 3 {
		t.Errorf("len(lines) != 3: %v", len(lines))
	}
	if lines[0].String() != `app.kubernetes.io/com-ponent="test"` {
		t.Errorf("lines[0] has an unexpected value: %s", lines[0].String())
	}
	if lines[1].String() != `rack="22"` {
		t.Errorf("lines[1] has an unexpected value: %s", lines[0].String())
	}
	if lines[2].String() != `test="something"` {
		t.Errorf("lines[2] has an unexpected value: %s", lines[0].String())
	}
}

package uri

import (
	"testing"
)

func TestUri(t *testing.T) {
	u, e := Parse("local://D:/w/ds/s")
	t.Log(u.Full(), e)
}

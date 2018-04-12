package reldel

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPatch1(t *testing.T) {
	p := GetPatch("", "cat")
	assert.Equal(t, "cat", p.PatchIotas[0].Between)
	assert.Equal(t, "cat", ApplyPatch("", p))
	p = GetPatch("hungry cat", "hungry orange cats")
	bP, _ := json.MarshalIndent(p, "", " ")
	fmt.Println(string(bP))
	assert.Equal(t, "hungry orange cats", ApplyPatch("hungry cat", p))
}

package reldel

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPatch1(t *testing.T) {
	p := GetPatch("", "cat")
	fmt.Println(p)
	assert.Equal(t, "cat", p.PatchIotas[0].Between)
	assert.Equal(t, "cat", ApplyPatch("", p))
	// p = GetPatch("hungry cat", "hungry orange cats")
	// fmt.Println(p)
}

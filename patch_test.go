package reldel

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPatch1(t *testing.T) {
	p := GetPatch("", "ca-t")
	assert.Equal(t, "ca-t", p.PatchIotas[0].Between)
	assert.Equal(t, "ca-t", ApplyPatch("", p))
}

func TestGetPatch2(t *testing.T) {
	b1, _ := ioutil.ReadFile("testing/1")
	b2, _ := ioutil.ReadFile("testing/2")
	p := GetPatch(string(b1), string(b2))
	bP, _ := json.MarshalIndent(p, "", " ")
	fmt.Println(string(bP))
	assert.Equal(t, string(b2), ApplyPatch(string(b1), p))
	ioutil.WriteFile("3", []byte(ApplyPatch(string(b1), p)), 0644)
}

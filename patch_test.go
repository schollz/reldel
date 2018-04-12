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
	assert.Equal(t, "ca-t", ApplyPatch("", p))
}

func TestGetPatch2(t *testing.T) {
	b1, _ := ioutil.ReadFile("testing/3")
	b2, _ := ioutil.ReadFile("testing/4")
	p := GetPatch(string(b1), string(b2))
	bP, _ := json.MarshalIndent(p, "", " ")
	fmt.Println(string(bP))
	assert.Equal(t, string(b2), ApplyPatch(string(b1), p))
	ioutil.WriteFile("out.json", bP, 0644)
}

func TestGetPatch3(t *testing.T) {
	s1 := "The cow jumped over the moon"
	s2 := "The brown cow leaped over the moon"
	s3 := "The cow jumped over the full moon"
	fmt.Println(s1, s2, s3)
	p := GetPatch(s1, s2)
	assert.Equal(t, s2, ApplyPatch(s1, p))
	p2 := GetPatch(s1, s3)
	assert.Equal(t, s3, ApplyPatch(s1, p2))

	assert.Equal(t, "The brown cow leaped over the full moon", ApplyPatch(s2, p2))
}

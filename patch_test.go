package reldel

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"runtime/pprof"
	"testing"

	"github.com/schollz/nwalgo"
	"github.com/stretchr/testify/assert"
)

func BenchmarkFileGetPatch(b *testing.B) {
	b1, _ := ioutil.ReadFile("testing/3")
	b2, _ := ioutil.ReadFile("testing/4")
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		GetPatch(string(b1), string(b2))
	}
}

func BenchmarkAlign(b *testing.B) {
	b1, _ := ioutil.ReadFile("testing/3")
	b2, _ := ioutil.ReadFile("testing/4")
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		nwalgo.Align(string(b1), string(b2), match, mismatch, gap)
	}
}

func TestGetPatch1(t *testing.T) {
	p := GetPatch("", "ca-t")
	s, err := ApplyPatch("", p)
	assert.Nil(t, err)
	assert.Equal(t, "ca-t", s)
}

func TestGetPatchHard(t *testing.T) {
	b1, _ := ioutil.ReadFile("testing/1")
	b2, _ := ioutil.ReadFile("testing/2")
	p := GetPatch(string(b1), string(b2))
	s, err := ApplyPatch(string(b1), p)
	assert.Nil(t, err)
	assert.Equal(t, string(b2), s)
}

func TestGetPatch2(t *testing.T) {
	f, err := os.Create("cpu.profile")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()
	b1, _ := ioutil.ReadFile("testing/3")
	b2, _ := ioutil.ReadFile("testing/4")
	p := GetPatch(string(b1), string(b2))
	s, err := ApplyPatch(string(b1), p)
	assert.Nil(t, err)
	assert.Equal(t, string(b2), s)
	bP, _ := json.MarshalIndent(p, "", " ")
	ioutil.WriteFile("out.json", bP, 0644)
}

func TestGetPatch3(t *testing.T) {
	s1 := "The cow jumped over the moon"
	s2 := "The brown cow leaped over the moon"
	s3 := "The cow jumped over the full moon"
	p := GetPatch(s1, s2)
	s, err := ApplyPatch(s1, p)
	assert.Nil(t, err)
	assert.Equal(t, s2, s)
	p2 := GetPatch(s1, s3)
	s, err = ApplyPatch(s1, p2)
	assert.Nil(t, err)
	assert.Equal(t, s3, s)

	s, err = ApplyPatch(s2, p2)
	assert.Nil(t, err)
	assert.Equal(t, "The brown cow leaped over the full moon", s)
}

func TestBadPatch(t *testing.T) {
	s1 := "The cow jumped"
	s2 := "The dog jumped"
	s3 := "The cat jumped"
	// patch should fail if the exact word is changed
	p := GetPatch(s1, s2)
	s, err := ApplyPatch(s1, p)
	assert.Nil(t, err)
	assert.Equal(t, s2, s)
	p2 := GetPatch(s1, s3)
	s, err = ApplyPatch(s1, p2)
	assert.Nil(t, err)
	assert.Equal(t, s3, s)

	s, err = ApplyPatch(s2, p2)
	assert.NotNil(t, err)
}

func TestGoodPatch(t *testing.T) {
	s1 := "The cow jumped"
	s2 := "The cows jumped"
	s3 := "The cats jumped"
	// patch should fail if the exact word is changed
	p := GetPatch(s1, s2)
	s, err := ApplyPatch(s1, p)
	assert.Nil(t, err)
	assert.Equal(t, s2, s)
	p2 := GetPatch(s1, s3)
	s, err = ApplyPatch(s1, p2)
	assert.Nil(t, err)
	assert.Equal(t, s3, s)

	s, err = ApplyPatch(s2, p2)
	assert.Nil(t, err)
}

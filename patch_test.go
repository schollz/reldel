package reldel

import (
	"encoding/json"
	"fmt"
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
		GetPatch(b1, b2)
	}
}

func BenchmarkAlign(b *testing.B) {
	b1, _ := ioutil.ReadFile("testing/3")
	b2, _ := ioutil.ReadFile("testing/4")
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		nwalgo.Align((b1), (b2), match, mismatch, gap)
	}
}

func TestGetPatchEmpty(t *testing.T) {
	p := GetPatch([]byte(""), []byte(""))
	s, err := ApplyPatch([]byte(""), p)
	assert.Nil(t, err)
	assert.Equal(t, []byte(nil), s)
}

func TestGetPatch1(t *testing.T) {
	s1 := []byte("")
	s2 := []byte("ca-t")
	p := GetPatch(s1, s2)
	fmt.Println(p)
	s, err := ApplyPatch(s1, p)
	assert.Nil(t, err)
	assert.Equal(t, s2, s)
	fmt.Println(string(s))
}

func TestGetPatch2(t *testing.T) {
	b1, _ := ioutil.ReadFile("testing/1")
	b2, _ := ioutil.ReadFile("testing/2")
	p := GetPatch(b1, b2)

	s, err := ApplyPatch(b1, p)
	assert.Nil(t, err)
	assert.Equal(t, b2, s)
}

func TestGetPatch3(t *testing.T) {
	f, err := os.Create("cpu.profile")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()
	b1, _ := ioutil.ReadFile("testing/3")
	b2, _ := ioutil.ReadFile("testing/4")
	p := GetPatch((b1), (b2))
	s, err := ApplyPatch((b1), p)
	assert.Nil(t, err)
	assert.Equal(t, (b2), s)
	bP, _ := json.MarshalIndent(p, "", " ")
	ioutil.WriteFile("out.json", bP, 0644)
}

func TestGetPatch4(t *testing.T) {
	s1 := []byte("The cow jumped over the moon")
	s2 := []byte("The brown cow leaped over the moon")
	// s3 := []byte("The cow jumped over the full moon")
	p := GetPatch(s1, s2)
	s, err := ApplyPatch(s1, p)
	assert.Nil(t, err)
	assert.Equal(t, s2, s)
	// p2 := GetPatch(s1, s3)
	// s, err = ApplyPatch(s1, p2)
	// assert.Nil(t, err)
	// assert.Equal(t, s3, s)

	// s, err = ApplyPatch(s2, p2)
	// assert.Nil(t, err)
	// assert.Equal(t, []byte("The brown cow leaped over the full moon"), s)
}

// func TestGetPatch5(t *testing.T) {
// 	s1 := `Sometimes when I want a recipe to cook something new I will find several recipes for the same thing and try to use them as a guide to generate an average or "consensus" recipe. This code should make it easy to generate consensus recipes (useful!) and also show variation between recipes (interesting!).`
// 	s2 := `Sometimes when I want a recipes to cook something new I will find several recipes for the same thing and try to use them as a guide to generate an average or "consensus" recipe. This code should make it easy to generate consensus (useful!) and also show variation between recipes (interesting!).`
// 	p := GetPatch([]byte(s1), []byte(s2))
// 	bP, _ := json.MarshalIndent(p, "", " ")
// 	fmt.Println(string(bP))
// }

func TestBadPatch(t *testing.T) {
	s1 := []byte("The cow jumped")
	s2 := []byte("The dog jumped")
	s3 := []byte("The cat jumped")
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
	s1 := []byte("The cow jumped")
	s2 := []byte("The cows jumped")
	s3 := []byte("The cats jumped")
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

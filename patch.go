package reldel

import (
	"bytes"
	"fmt"
	"math/rand"
	"time"

	"github.com/schollz/nwalgo"
)

const match = 10
const mismatch = -1
const gap = -5

// Patch is the unit of patching for a string
type Patch struct {
	HeadTail   [][]byte   `json:"h"`
	PatchIotas [][][]byte `json:"p"`
	Time       time.Time  `json:"t"`
}

func (p Patch) String() string {
	s := ""
	for i, patchiota := range p.PatchIotas {
		for j, patch := range patchiota {
			s += fmt.Sprintf("%d-%d) %s\n", i, j, patch)
		}
	}
	return s
}

// GetPatch will get a patch to transform text `s1` into text `s2`
func GetPatch(s1, s2 []byte) Patch {
	headTail := [][]byte{{}, {}, {}}
	for rLength := 3; rLength < 10; rLength++ {
		isGood := true
		for i := 0; i < 100; i++ {
			headTail = [][]byte{
				[]byte(randStringBytesMaskImprSrc(rLength + 2)),
				[]byte(randStringBytesMaskImprSrc(rLength + 1)),
				[]byte(randStringBytesMaskImprSrc(rLength)),
			}
			for _, h := range headTail {
				if bytes.Contains(s1, h) || bytes.Contains(s2, h) {
					isGood = false
					break
				}
			}
			if isGood {
				break
			}
		}
		if isGood {
			break
		}
	}
	// headTail = []string{"dr", "AJ", "Ld"}
	// headTail = [][]byte{[]byte("Nk"), []byte("7S"), R[]byte("7S")}
	// for _, b3 := range headTail {
	// 	fmt.Println(string(b3))
	// }
	s1 = bytes.Replace(s1, []byte("-"), headTail[2], -1)
	s2 = bytes.Replace(s2, []byte("-"), headTail[2], -1)
	patchIotas := [][][]byte{}
	aln1B, aln2B, _ := nwalgo.Align(s1, s2, match, mismatch, gap)
	aln1 := combineThreeByteArrays(headTail[0], aln1B, headTail[1])
	aln2 := combineThreeByteArrays(headTail[0], aln2B, headTail[1])

	// fmt.Println(string(bytes.Replace(aln1, []byte("\n"), []byte("#"), -1)))
	// fmt.Println(string(bytes.Replace(aln2, []byte("\n"), []byte("#"), -1)))
	for {
		if bytes.Compare(aln1, aln2) == 0 {
			break
		}
		p, nextStart := getPatchIota(aln1, aln2, headTail)
		patchIotas = append(patchIotas, p)
		copy(aln1[0:nextStart], aln2[0:nextStart])
	}
	return Patch{
		Time:       time.Now(),
		PatchIotas: patchIotas,
		HeadTail:   headTail,
	}
}

func combineThreeByteArrays(l, m, r []byte) []byte {
	s := make([]byte, len(l)+len(m)+len(r))
	copy(s[:len(l)], l)
	copy(s[len(l):len(l)+len(m)], m)
	copy(s[len(l)+len(m):], r)
	return s
}

// ApplyPatch will transform string `s` using the supplied patch. If there is a problem
// (which can occur if there has been an edit in the same place) then an error is thrown
// and the currently patched string is returned.
func ApplyPatch(s []byte, p Patch) ([]byte, error) {
	s = combineThreeByteArrays(p.HeadTail[0], bytes.Replace(s, []byte("-"), p.HeadTail[2], -1), p.HeadTail[1])
	// fmt.Println("applying to", string(s), s)

	var err error
	for _, patchIota := range p.PatchIotas {
		s, err = applyPatchIota(s, patchIota)
		if err != nil {
			return s, err
		}
	}
	s = bytes.Replace(s, p.HeadTail[0], []byte(""), -1)
	s = bytes.Replace(s, p.HeadTail[1], []byte(""), -1)
	s = bytes.Replace(s, []byte("-"), []byte(""), -1)
	s = bytes.Replace(s, p.HeadTail[2], []byte("-"), -1)
	s = bytes.Trim(s, "\x00")
	// fmt.Println(string(s))
	return s, err
}

func count(s, sep []byte, leaveAfter ...int) (count int) {
	// special case
	if len(sep) == 0 {
		return len(s)
	}
	n := 0
	for {
		i := bytes.Index(s, sep)
		if i == -1 {
			return n
		}
		n++
		if len(leaveAfter) > 0 && n >= leaveAfter[0] {
			return n
		}
		s = s[i+len(sep):]
	}
}

// applyPatchIota will replace byte s with p[2] which is flanked on the left
// by p[0] and on the right by p[1]
func applyPatchIota(s []byte, p [][]byte) ([]byte, error) {
	pos1 := bytes.Index(s, p[0])
	if pos1 == -1 {
		return []byte(""), fmt.Errorf("left index no longer exists")
	}
	// move position up if overlapping sequence is there (go only finds
	// the non-overlapping sequences)
	for i := pos1; i < pos1+len(p[0]); i++ {
		if i+len(p[0]) > len(s)-1 {
			break
		}
		if bytes.Compare(s[i:i+len(p[0])], p[0]) == 0 {
			pos1 = i
		}
	}
	pos1 = pos1 + len(p[0])
	pos2 := bytes.Index(s, p[1])
	if pos2 == -1 {
		return []byte(""), fmt.Errorf("right index no longer exists")
	}
	// fmt.Println("1", string(s[:pos1]), string(p[2]), string(s[pos2:]))
	return combineThreeByteArrays(s[:pos1], p[2], s[pos2:]), nil
}

func getPatchIota(aln1, aln2 []byte, headTail [][]byte) ([][]byte, int) {
	aln1WithoutWhiteSpace := bytes.Replace(aln1, []byte("-"), []byte(""), -1)
	// fmt.Print("\n")
	// fmt.Println(aln1)
	// fmt.Println(aln2)
	// abcdef
	// ab-def
	//   ^
	bookends := []int{0, 0, 0, 0}
	for i := bookends[0]; i < len(aln1); i++ {
		if aln1[i] != aln2[i] {
			bookends[1] = i
			break
		}
	}
	// fmt.Printf("%+v, '%s'\n", bookends, aln1[bookends[0]:bookends[1]])

	// find unique subsequence in front
	for i := bookends[1]; i >= bookends[0]; i-- {
		if count(aln1WithoutWhiteSpace, bytes.Replace(aln1[i:bookends[1]], []byte("-"), []byte(""), -1), 2) > 1 {
			continue
		}
		bookends[0] = i
	}
	// fmt.Printf("%+v, '%s'\n", bookends, aln1[bookends[0]:bookends[1]])

	// find where next matching subsequence begins
	bookends[2] = bookends[1]
	for j := 0; j < 300; j++ {
		for i := bookends[2]; i < len(aln1); i++ {
			bookends[2] = i
			if aln1[i] == aln2[i] {
				break
			}
		}
		// fmt.Printf("1 %+v, '%s' %d\n", bookends, aln1[bookends[2]:], len(aln1))
		// find where the next matching sequence ends
		bookends[3] = bookends[2]
		for i := bookends[2]; i < len(aln1); i++ {
			if aln1[i] != aln2[i] {
				bookends[3] = i
				break
			}
		}
		if bookends[2] == bookends[3] {
			bookends[3] = len(aln1)
		}
		// fmt.Printf("2 %+v, '%s'\n", bookends, aln1[bookends[2]:bookends[3]])
		if count(aln1WithoutWhiteSpace, bytes.Replace(aln1[bookends[2]:bookends[3]], []byte("-"), []byte(""), -1), 2) == 1 {
			break
		}
		bookends[2] = bookends[3]
	}
	// now that we have a second matching sequence, try to reduce it
	for bookends[3] = bookends[2] + 1; bookends[3] < len(aln1); bookends[3]++ {
		// fmt.Println(bookends, aln1[bookends[2]:bookends[3]])
		if count(aln1WithoutWhiteSpace, bytes.Replace(aln1[bookends[2]:bookends[3]], []byte("-"), []byte(""), -1), 2) == 1 {
			break
		}
	}

	// fmt.Println(bookends)
	if bookends[3] > len(aln1) {
		bookends[3] = len(aln1)
	}
	left := aln1[bookends[0]:bookends[1]]
	right := aln1[bookends[2]:bookends[3]]
	insertion := aln2[bookends[1]:bookends[2]]
	insertion = bytes.Replace(insertion, []byte("-"), []byte(""), -1)

	// fmt.Printf("l: '%s', r: '%s', i: '%s'\n", left, right, insertion)
	return [][]byte{left, right, insertion}, bookends[2]
}

var src = rand.NewSource(time.Now().UnixNano())

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func randStringBytesMaskImprSrc(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

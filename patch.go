package reldel

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/schollz/nwalgo"
)

const match = 1
const mismatch = -1
const gap = -1

type Patch struct {
	HeadTail   []string    `json:"h"`
	PatchIotas []PatchIota `json:"p"`
	Time       time.Time   `json:"t"`
}

type PatchIota struct {
	Left    string `json:"l"`
	Right   string `json:"r"`
	Between string `json:"b"`
}

func GetPatch(s1, s2 string) Patch {
	headTail := []string{"start>>>>>>>>>", "<<<<<<<<<<end", "**dash**"}
	for rLength := 2; rLength < 10; rLength++ {
		isGood := true
		for i := 0; i < 1000; i++ {
			headTail = []string{RandStringBytesMaskImprSrc(rLength), RandStringBytesMaskImprSrc(rLength), RandStringBytesMaskImprSrc(rLength)}
			if headTail[0] == headTail[1] || headTail[1] == headTail[2] || headTail[0] == headTail[2] {
				continue
			}
			for _, h := range headTail {
				if strings.Contains(s1, h) || strings.Contains(s2, h) {
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
	s1 = strings.Replace(s1, "-", headTail[2], -1)
	s2 = strings.Replace(s2, "-", headTail[2], -1)
	patchIotas := []PatchIota{}
	aln1, aln2, _ := nwalgo.Align(s1, s2, match, mismatch, gap)
	aln1 = headTail[0] + aln1 + headTail[1]
	aln2 = headTail[0] + aln2 + headTail[1]

	// fmt.Println(strings.Replace(aln1, "\n", "#", -1))
	// fmt.Println(strings.Replace(aln2, "\n", "#", -1))
	for {
		if aln1 == aln2 {
			break
		}
		p, nextStart := getPatchIota(aln1, aln2, headTail)
		patchIotas = append(patchIotas, p)
		aln1 = aln2[0:nextStart] + aln1[nextStart:]
	}
	return Patch{
		Time:       time.Now(),
		PatchIotas: patchIotas,
		HeadTail:   headTail,
	}
}

func ApplyPatch(s string, p Patch) string {
	s = p.HeadTail[0] + strings.Replace(s, "-", p.HeadTail[2], -1) + p.HeadTail[1]
	for _, patchIota := range p.PatchIotas {
		s = applyPatchIota(s, patchIota)
	}
	s = strings.Replace(s, "-", "", -1)
	s = strings.Replace(s, p.HeadTail[2], "-", -1)
	s = strings.Replace(s, p.HeadTail[0], "", -1)
	s = strings.Replace(s, p.HeadTail[1], "", -1)

	return s
}

func count(s, substr string) int {
	return strings.Count(s, substr)
}

func applyPatchIota(s string, p PatchIota) string {
	pos1 := strings.Index(s, p.Left)
	if pos1 == -1 {
		fmt.Printf("\n\n'%s'", s)
		fmt.Printf("\n\n'%s'", p.Left)
		panic("problem")
	}
	// move position up if overlapping sequence is there (go only finds
	// the non-overlapping sequences)
	for i := pos1; i < pos1+len(p.Left); i++ {
		if i+len(p.Left) > len(s)-1 {
			break
		}
		if s[i:i+len(p.Left)] == p.Left {
			pos1 = i
		}
	}
	pos1 = pos1 + len(p.Left)
	pos2 := strings.Index(s, p.Right)
	if pos2 == -1 {
		fmt.Printf("\n\n'%s'", s)
		fmt.Printf("\n\n'%s'", p.Right)
		panic("problem")
	}
	s = s[:pos1] + p.Between + s[pos2:]
	// fmt.Println(s)
	return s
}

func getPatchIota(aln1, aln2 string, headTail []string) (PatchIota, int) {
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
	for i := bookends[0]; i < bookends[1]; i++ {
		if count(aln1, aln1[i:bookends[1]]) > 1 {
			break
		}
		bookends[0] = i
	}
	// fmt.Printf("%+v, '%s'\n", bookends, aln1[bookends[0]:bookends[1]])

	// find where next matching subsequence begins
	bookends[2] = bookends[1]
	for j := 0; j < 30; j++ {
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
		if count(aln1, aln1[bookends[2]:bookends[3]]) == 1 {
			break
		}
		bookends[2] = bookends[3]
	}
	// now that we have a second matching sequence, try to reduce it
	for bookends[3] = bookends[2] + 1; bookends[3] < len(aln1); bookends[3]++ {
		// fmt.Println(bookends, aln1[bookends[2]:bookends[3]])
		if count(aln1, aln1[bookends[2]:bookends[3]]) == 1 {
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
	// insertion = strings.Replace(insertion, "-", "", -1)

	// fmt.Printf("l: '%s', r: '%s', i: '%s'\n", left, right, insertion)
	return PatchIota{
		Left:    left,
		Right:   right,
		Between: insertion,
	}, bookends[2]
}

var src = rand.NewSource(time.Now().UnixNano())

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func RandStringBytesMaskImprSrc(n int) string {
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

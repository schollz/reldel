package reldel

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/schollz/nwalgo"
)

const HEAD = "start>>>>>"
const TAIL = "<<<<<<<end"
const match = 1
const mismatch = -1
const gap = -1

type Patch struct {
	PatchIotas []PatchIota `json:"p"`
	Time       time.Time   `json:"t"`
}

type PatchIota struct {
	Left    string `json:"l"`
	Right   string `json:"r"`
	Between string `json:"b"`
}

func GetPatch(s1, s2 string) Patch {
	s1 = HEAD + strings.Replace(s1, "-", "**dash**", -1) + TAIL
	s2 = HEAD + strings.Replace(s2, "-", "**dash**", -1) + TAIL
	patchIotas := []PatchIota{}
	aln1, aln2, _ := nwalgo.Align(s1, s2, match, mismatch, gap)
	for {
		if aln1 == aln2 {
			break
		}
		p, nextStart := getPatchIota(aln1, aln2)
		patchIotas = append(patchIotas, p)
		aln1 = aln2[0:nextStart] + aln1[nextStart:]
	}
	bP, _ := json.MarshalIndent(patchIotas, "", " ")
	fmt.Println(string(bP))
	return Patch{
		Time:       time.Now(),
		PatchIotas: patchIotas,
	}
}

func ApplyPatch(s1 string, p Patch) string {
	for _, patchIota := range p.PatchIotas {
		s1 = applyPatchIota(s1, patchIota)
	}
	return s1
}

func count(s, substr string) int {
	// need to find overlapping matches
	return strings.Count(s, substr)
}

func applyPatchIota(s string, p PatchIota) string {
	s = HEAD + s + TAIL
	locs := regexp.MustCompile("(?=("+p.Left+")").FindAllStringIndex(s, -1)
	fmt.Println(p.Left, locs)
	pos1 := locs[len(locs)-1][1]
	pos2 := regexp.MustCompile(p.Right).FindAllStringIndex(s, 1)[0][0]
	fmt.Println(s, pos1, pos2)
	return strings.TrimSuffix(strings.TrimPrefix(s[:pos1]+p.Between+s[pos2:], HEAD), TAIL)
}

func getPatchIota(aln1, aln2 string) (PatchIota, int) {

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
	fmt.Printf("%+v, '%s'\n", bookends, aln1[bookends[0]:bookends[1]])

	// find unique subsequence in front
	for i := bookends[0]; i < bookends[1]; i++ {
		if count(aln1, aln1[i:bookends[1]]) > 1 {
			break
		}
		bookends[0] = i
	}
	fmt.Printf("%+v, '%s'\n", bookends, aln1[bookends[0]:bookends[1]])

	// find where next matching subsequence begins
	bookends[2] = bookends[1]
	for j := 0; j < 30; j++ {
		for i := bookends[2]; i < len(aln1); i++ {
			bookends[2] = i
			if aln1[i] == aln2[i] {
				break
			}
		}
		fmt.Printf("1 %+v, '%s' %d\n", bookends, aln1[bookends[2]:], len(aln1))
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
		fmt.Printf("2 %+v, '%s'\n", bookends, aln1[bookends[2]:bookends[3]])
		if count(aln1, aln1[bookends[2]:bookends[3]]) == 1 {
			break
		}
		bookends[2] = bookends[3]
	}
	// now that we have a second matching sequence, try to reduce it
	for bookends[3] = bookends[2] + 1; bookends[3] < len(aln1); bookends[3]++ {
		fmt.Println(bookends, aln1[bookends[2]:bookends[3]])
		if count(aln1, aln1[bookends[2]:bookends[3]]) == 1 {
			break
		}
	}

	left := aln1[bookends[0]:bookends[1]]
	left = strings.Replace(left, "**dash**", "-", -1)
	right := aln1[bookends[2]:bookends[3]]
	right = strings.Replace(right, "**dash**", "-", -1)
	insertion := aln2[bookends[1]:bookends[2]]
	insertion = strings.Replace(insertion, "-", "", -1)
	insertion = strings.Replace(insertion, "**dash**", "-", -1)

	fmt.Printf("l: '%s', r: '%s', i: '%s'\n", left, right, insertion)
	return PatchIota{
		Left:    left,
		Right:   right,
		Between: insertion,
	}, bookends[2]
}

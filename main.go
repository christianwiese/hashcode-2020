package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type slide []Picture
type Result []slide

func main() {
	//input := ParseInput(os.Args[1])
	//
	//fmt.Println("WORKING ON: ", os.Args[1])
	//
	//
	fmt.Println("final score")
}

func Dump(out Result, file string) {
	output := fmt.Sprintf("./%s.out", strings.TrimSuffix(file, ".in"))
	f, _ := os.Create(output)
	defer f.Close()

	w := bufio.NewWriter(f)

	num := len(out)
	w.WriteString(fmt.Sprintf("%d", num))
	w.WriteString("\n")
	for _, c := range out {
		if len(c) == 1 {
			w.WriteString(fmt.Sprintf("%d", c[0].index))
		} else {
			w.WriteString(fmt.Sprintf("%d %d", c[0].index, c[1].index))
		}
		w.WriteString("\n")
	}
	w.Flush()
}

func ParseInt(v string) int {
	x, err := strconv.Atoi(v)
	if err != nil {
		panic(err)
	}
	return x
}

type Picture struct {
	index int
	vert  bool
	tags  map[string]bool
}

func ParseInput(path string) []Picture {

	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	scanner.Scan()

	numPictures := ParseInt(scanner.Text())

	pictures := []Picture{}

	for i := 0; i < numPictures; i++ {
		scanner.Scan()
		tmp := strings.Split(scanner.Text(), " ")
		var picture Picture
		if tmp[0] == "H" {
			picture.vert = false
		} else {
			picture.vert = true
		}

		picture.index = i
		slice := tmp[2:]
		s := make(map[string]bool, len(slice))
		for _, t := range slice {
			s[t] = true
		}
		picture.tags = s
		pictures = append(pictures, picture)
	}

	return pictures
}

func Score(res Result) int {
	score := 0
	for i := 0; i < len(res)-1; i++ {

		min := getSingleScore(res[i], res[i+1])
		score += min
	}
	return score
}

func getSingleScore(left slide, right slide) int {
	thisTags := getTagSet(left)
	nextTags := getTagSet(right)

	intersect := thisTags.intersect(nextTags)

	if intersect == 0 {
		return 0
	}

	l := len(thisTags) - intersect

	if l == 0 {
		return 0
	}

	r := len(nextTags) - intersect

	min := l
	if intersect < min {
		min = intersect
	}

	if r < min {
		min = r
	}

	return min
}

type tagSet map[string]bool

func getTagSet(sl slide) tagSet {
	res := map[string]bool{}
	res = sl[0].tags
	if len(sl) == 2 {
		for t := range sl[1].tags {
			res[t] = true
		}
	}
	return res
}

func (ts tagSet) intersect(right tagSet) int {
	score := 0
	for k := range ts {
		_, ok := right[k]
		if ok {
			score++
		}
	}
	return score
}

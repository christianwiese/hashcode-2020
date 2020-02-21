package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

func main() {
	input := ParseInput(os.Args[1])

	fmt.Println("WORKING ON: ", os.Args[1])

	r := rand.New(rand.NewSource(time.Now().Unix()))
	fmt.Println(r)

	result := []*Result{}

	libraries := input.libraries

	for _, l := range libraries {
		score := 0
		for _, b := range l.books {
			score += b.score
		}
		l.score = score
	}

	//Sort by signup
	sort.SliceStable(libraries, func(i, j int) bool {
		return libraries[i].signUp < libraries[j].signUp
	})
	//sort by score
	sort.SliceStable(libraries, func(i, j int) bool {
		return libraries[i].score > libraries[j].score
	})

	for _, l := range libraries {

		books := l.books
		sort.SliceStable(books, func(i, j int) bool {
			return books[i].score > books[j].score
		})
		resBooks := []int{}

		for _, b := range books {
			resBooks = append(resBooks, b.index)
		}

		resL := &Result{
			index:       l.index,
			signUp:      l.signUp,
			booksPerDay: l.booksPerDay,
			books:       resBooks,
		}
		result = append(result, resL)
	}

	finalScore := 0

	for i := 0; i < 10000000; i++ {
		r1 := r.Intn(len(result))
		r2 := r.Intn(len(result))

		result[r1], result[r2] = result[r2], result[r1]
		score := Score(input, result)
		if score > finalScore {
			finalScore = score
			fmt.Println(finalScore)
			Dump(result, os.Args[1])
		} else {
			result[r2], result[r1] = result[r1], result[r2]
		}
	}

	Dump(result, os.Args[1]) //fmt.Println("final score")
}

func Dump(o []*Result, file string) {
	scannedBooks := NewSet()
	finalRes := []*Result{}
	for _, l := range o {

		books := l.books
		resBooks := []int{}
		for _, b := range books {
			if !scannedBooks.Contains(b) {
				resBooks = append(resBooks, b)
				scannedBooks.Add(b)
			}
		}
		if len(resBooks) == 0 {
			continue
		}

		resL := &Result{
			index:       l.index,
			signUp:      l.signUp,
			booksPerDay: l.booksPerDay,
			books:       resBooks,
		}
		finalRes = append(finalRes, resL)
	}


	output := fmt.Sprintf("./%s.out", strings.TrimSuffix(file, ".in"))
	f, _ := os.Create(output)
	defer f.Close()

	w := bufio.NewWriter(f)

	num := len(finalRes)
	w.WriteString(fmt.Sprintf("%d", num))
	w.WriteString("\n")
	for _, c := range finalRes {
		w.WriteString(fmt.Sprintf("%d %d", c.index, len(c.books)))
		w.WriteString("\n")
		valuesText := []string{}
		for _, b := range c.books {
			text := strconv.Itoa(b)
			valuesText = append(valuesText, text)
		}
		w.WriteString(fmt.Sprintf("%s", strings.Join(valuesText, " ")))
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

type Input struct {
	books     map[int]int
	libraries []*Library
	days      int
}

type Book struct {
	index int
	score int
}

type Library struct {
	index       int
	signUp      int
	booksPerDay int
	books       []*Book
	score       int
}

type Result struct {
	index       int
	signUp      int
	booksPerDay int
	books       []int
}

func ParseInput(path string) *Input {

	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	scanner.Scan()
	numBooks := ParseInt(scanner.Text())
	scanner.Scan()
	numLib := ParseInt(scanner.Text())
	scanner.Scan()
	numDays := ParseInt(scanner.Text())

	books := make(map[int]int, numBooks)
	for i := 0; i < numBooks; i++ {
		scanner.Scan()
		books[i] = ParseInt(scanner.Text())
	}

	libraries := []*Library{}

	for i := 0; i < numLib; i++ {

		library := &Library{index: i}

		scanner.Scan()
		numBooks := ParseInt(scanner.Text())
		scanner.Scan()
		library.signUp = ParseInt(scanner.Text())
		scanner.Scan()
		library.booksPerDay = ParseInt(scanner.Text())

		bks := []*Book{}
		for i := 0; i < numBooks; i++ {
			scanner.Scan()
			index := ParseInt(scanner.Text())
			book := &Book{index: index, score: books[index]}
			bks = append(bks, book)
		}
		library.books = bks

		libraries = append(libraries, library)
	}

	return &Input{
		books:     books,
		libraries: libraries,
		days:      numDays,
	}
}

func Score(input *Input, res []*Result) int {
	score := 0
	daysLeft := input.days
	scannedBooks := NewSet()
	for _, l := range res {
		daysLeft -= l.signUp
		if daysLeft < 0 {
			break
		}
		actualBooks := []int{}
		for _, b := range l.books {
			if !scannedBooks.Contains(b) {
				actualBooks = append(actualBooks, b)
			}
		}
		scanDays := int(math.Ceil(float64(len(actualBooks)) / float64(l.booksPerDay)))
		if scanDays <= daysLeft {
			scannedBooks.AddAll(actualBooks)
		} else {
			possible := daysLeft * l.booksPerDay
			scannedBooks.AddAll(actualBooks[0:possible])
		}
	}

	for _, b := range scannedBooks.ToSlice() {
		score += input.books[b]
	}

	return score
}

// Set is a set of strings
type Set map[int]struct{}

// NewSet creates a new set
func NewSet() *Set {
	s := Set(make(map[int]struct{}))
	return &s
}

// Add adds a single value to the set
func (s *Set) Add(value int) *Set {
	(*s)[value] = struct{}{}
	return s
}

// AddAll adds all values to the set
func (s *Set) AddAll(values []int) *Set {
	for _, v := range values {
		s.Add(v)
	}
	return s
}

// Contains checks if the set contains the given value
func (s Set) Contains(value int) bool {
	if _, ok := s[value]; ok {
		return true
	}
	return false
}

// ToSlice transforms a set of strings into a slice of strings
func (s *Set) ToSlice() []int {
	keys := make([]int, len(*s))

	i := 0
	for k := range *s {
		keys[i] = k
		i++
	}

	return keys
}

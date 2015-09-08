package main

import (
	"flag"
	"log"
	"regexp"
	"strings"

	_ "github.com/chrislusf/netchan/example/driver"
	"github.com/chrislusf/netchan/example/flame"
)

func main() {
	flag.Parse()

	// test1()

	// test2()

	test3()

}

func test1() {
	flame.NewContext().TextFile(
		"/etc/passwd", 3,
	).Filter(func(line string) bool {
		println(line)
		return !strings.HasPrefix(line, "#")
	}).Map(func(line string, ch chan string) {
		for _, token := range strings.Split(line, ":") {
			ch <- token
		}
	}).Map(func(key string) int {
		return 1
	}).Reduce(func(x int, y int) int {
		return x + y
	}).Map(func(x int) {
		println("count:", x)
	})
}

func test2() {
	flame.NewContext().TextFile(
		"/etc/hosts", 7,
	).Partition(
		2,
	).Map(func(line string) string {
		return line
	}).Sort(func(a string, b string) bool {
		if strings.Compare(a, b) < 0 {
			return true
		}
		return false
	}).Map(func(line string) {
		println(line)
	})

}

func test3() {
	words := flame.NewContext().TextFile(
		"/etc/passwd", 3,
	).Filter(func(line string) bool {
		return !strings.HasPrefix(line, "#")
	}).Map(func(line string, ch chan string) {
		for _, token := range strings.Split(line, ":") {
			ch <- token
		}
	}).Map(func(line string) (string, string) {
		return line, line
	})

	words.Join(words).Map(func(key, left, right string) {
		println(key, ":", left, ":", right)
	})

}

func test4() {
	reg, err := regexp.Compile("[^A-Za-z0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	tokenizer := func(line string, ch chan string) {
		line = reg.ReplaceAllString(line, "-")
		for _, token := range strings.Split(line, "-") {
			ch <- strings.ToLower(token)
		}
	}
	ctx := flame.NewContext()
	leftWords := ctx.TextFile(
		"/etc/passwd", 3,
	).Map(tokenizer).Map(func(t string) (string, int) {
		return t, 1
	}).LocalSort(nil).LocalReduceByKey(func(x, y int) int {
		return x + y
	})

	rightWords := ctx.TextFile(
		"word_count.go", 3,
	).Map(tokenizer).Sort(nil)

	leftWords.Join(rightWords).Map(func(key, left, right string) {
		println(key, ":", left, ":", right)
	})

}

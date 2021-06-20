package go_bible

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"testing"
	"text/template"
	"time"
)

func TestArray(t *testing.T) {
	a := [3]int{1, 2}
	for i, v := range a {
		fmt.Printf("%d %d\n", i, v)
	}

	q := [...]int{1, 2, 3}
	fmt.Printf("%T\n", q)

	r := [...]int{99: -1}
	for i, v := range r {
		fmt.Printf("%d %d\n", i, v)
	}
}

type Currency int

const (
	USD Currency = iota // 美元
	EUR                 // 欧元
	GBP                 // 英镑
	RMB                 // 人民币
)

func TestArray2(t *testing.T) {
	symbol := [...]string{USD: "$", EUR: "€", GBP: "￡", RMB: "￥"}
	fmt.Println(RMB, symbol[RMB])
}

func TestArray3(t *testing.T) {
	c1 := sha256.Sum256([]byte("x"))
	c2 := sha256.Sum256([]byte("X"))
	fmt.Printf("%x\n%x\n%t\n%T\n", c1, c2, c1 == c2, c1)
}

func TestSlice(t *testing.T) {
	a := [...]int{0, 1, 2, 3, 4, 5}
	reverse(a[:])
	fmt.Println(a) // "[5 4 3 2 1 0]"
}
func reverse(s []int) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func TestSlice2(t *testing.T) {
	var runes []rune
	for _, r := range "Hello, 世界" {
		runes = append(runes, r)
	}
	fmt.Printf("%q\n", runes) // "['H' 'e' 'l' 'l' 'o' ',' ' ' '世' '界']"
}

func TestSlice3(t *testing.T) {
	var x, y []int
	for i := 0; i < 10; i++ {
		y = appendInt(x, i)
		fmt.Printf("%d cap=%d\t%v\n", i, cap(y), y)
		x = y
	}
}
func appendInt(x []int, y int) []int {
	var z []int
	zlen := len(x) + 1
	if zlen <= cap(x) {
		z = x[:zlen]
	} else {
		zcap := zlen
		if zcap < 2*len(x) {
			zcap = 2 * len(x)
		}
		z = make([]int, zlen, zcap)
		copy(z, x) // a built-in function; see text
	}
	z[len(x)] = y
	return z
}

func TestSlice4(t *testing.T) {
	var x []int
	x = append(x, 1)
	x = append(x, 2, 3)
	x = append(x, 4, 5, 6)
	x = append(x, x...)
	fmt.Println(x)
}

func TestSlice5(t *testing.T) {
	data := []string{"one", "", "three"}
	fmt.Printf("%q\t%p\n", nonempty(data), nonempty(data)) // `["one" "three"]`
	fmt.Printf("%q\t%p\n", data, data)                     // `["one" "three" "three"]`
}
func nonempty(strings []string) []string {
	i := 0
	for _, s := range strings {
		if s != "" {
			strings[i] = s
			i++
		}
	}
	return strings[:i]
}

func TestSlice6(t *testing.T) {
	s := []int{5, 6, 7, 8, 9}
	s2 := remove(s, 2)
	fmt.Println(s2) // "[5 6 8 9]"
	fmt.Printf("%p\t%p\n", s, s2)
}
func remove(slice []int, i int) []int {
	copy(slice[i:], slice[i+1:])
	return slice[:len(slice)-1]
}

func TestMap(t *testing.T) {
	ages := make(map[string]int)
	ages["bob"] = ages["bob"] + 1
	ages["alice"] = 31
	ages["charlie"] = 34
	// 遍历map，取key-value
	for name, age := range ages {
		fmt.Printf("%s\t%d\n", name, age)
	}

	// 排序
	names := make([]string, 0, len(ages))
	// 遍历map，只取key
	for name := range ages {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		fmt.Printf("%s\t%d\n", name, ages[name])
	}

	// 判断是否存在
	age, ok := ages["bob"]
	if !ok {
		fmt.Println(age)
	}
}

func TestMap2(t *testing.T) {
	var ages map[string]int
	fmt.Println(ages == nil)    // "true"
	fmt.Println(len(ages) == 0) // "true"
	// 下面语句会报错，需要先make
	ages["carol"] = 21
}

func TestMap3(t *testing.T) {
	var ages map[string]int
	fmt.Println(ages == nil)    // "true"
	fmt.Println(len(ages) == 0) // "true"
	// 下面语句会报错，需要先make
	ages["carol"] = 21
}

type tree struct {
	value       int
	left, right *tree
}

func TestStruct(t *testing.T) {
	seen := make(map[string]struct{})
	s := "blob"
	if _, ok := seen[s]; !ok {
		seen[s] = struct{}{}
	}
}

func Sort(values []int) {
	var root *tree
	for _, v := range values {
		root = add(root, v)
	}
	appendValues(values[:0], root)
}
func appendValues(values []int, t *tree) []int {
	if t != nil {
		values = appendValues(values, t.left)
		values = append(values, t.value)
		values = appendValues(values, t.right)
	}
	return values
}
func add(t *tree, value int) *tree {
	if t == nil {
		// Equivalent to return &tree{value: value}.
		t = new(tree)
		t.value = value
		return t
	}
	if value < t.value {
		t.left = add(t.left, value)
	} else {
		t.right = add(t.right, value)
	}
	return t
}

type Movie struct {
	Title  string
	Year   int  `json:"released"`
	Color  bool `json:"color,omitempty"`
	Actors []string
}

func TestJson(t *testing.T) {
	var movies = []Movie{
		{Title: "Casablanca", Year: 1942, Color: false,
			Actors: []string{"Humphrey Bogart", "Ingrid Bergman"}},
		{Title: "Cool Hand Luke", Year: 1967, Color: true,
			Actors: []string{"Paul Newman"}},
		{Title: "Bullitt", Year: 1968, Color: true,
			Actors: []string{"Steve McQueen", "Jacqueline Bisset"}},
	}
	// 不缩进
	data, err := json.Marshal(movies)
	if err != nil {
		log.Fatalf("JSON marshaling failed: %s", err)
	}
	fmt.Printf("%s\n", data)
	// 缩进
	data, err = json.MarshalIndent(movies, "", "    ")
	if err != nil {
		log.Fatalf("JSON marshaling failed: %s", err)
	}
	fmt.Printf("%s\n", data)

	// 反序列化
	var titles []struct{ Title string }
	if err := json.Unmarshal(data, &titles); err != nil {
		log.Fatalf("JSON unmarshaling failed: %s", err)
	}
	fmt.Println(titles)
}

const IssuesURL = "https://api.github.com/search/issues"

type IssuesSearchResult struct {
	TotalCount int `json:"total_count"`
	Items      []*Issue
}
type Issue struct {
	Number    int
	HTMLURL   string `json:"html_url"`
	Title     string
	State     string
	User      *User
	CreatedAt time.Time `json:"created_at"`
	Body      string    // in Markdown format
}
type User struct {
	Login   string
	HTMLURL string `json:"html_url"`
}

func TestJson2(t *testing.T) {
	result, err := SearchIssues(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d issues:\n", result.TotalCount)
	for _, item := range result.Items {
		fmt.Printf("#%-5d %9.9s %.55s\n",
			item.Number, item.User.Login, item.Title)
	}
}
func SearchIssues(terms []string) (*IssuesSearchResult, error) {
	q := url.QueryEscape(strings.Join(terms, " "))
	resp, err := http.Get(IssuesURL + "?q=" + q)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("search query failed: %s", resp.Status)
	}

	var result IssuesSearchResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		resp.Body.Close()
		return nil, err
	}
	resp.Body.Close()
	return &result, nil
}

const templ = `{{.TotalCount}} issues:
{{range .Items}}----------------------------------------
Number: {{.Number}}
User:   {{.User.Login}}
Title:  {{.Title | printf "%.64s"}}
Age:    {{.CreatedAt | daysAgo}} days
{{end}}`

func daysAgo(t time.Time) int {
	return int(time.Since(t).Hours() / 24)
}
func TestText(t *testing.T) {
	var report = template.Must(template.New("issuelist").
		Funcs(template.FuncMap{"daysAgo": daysAgo}).
		Parse(templ))
	result, err := SearchIssues(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	if err := report.Execute(os.Stdout, result); err != nil {
		log.Fatal(err)
	}
}

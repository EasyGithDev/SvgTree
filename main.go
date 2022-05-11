package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	svg "github.com/ajstarks/svgo"
)

const (
	canvasWidth  = 800
	canvasHeight = 600
)

var (
	fontSize  string = "15"
	textStyle string = "text-anchor:middle;font-size:" + fontSize + "px;fill:black"
	lineStyle string = "stroke:rgb(255,0,0);stroke-width:2"
	rectStyle string = "fill:rgb(255,255,255);stroke-width:2;stroke:rgb(0,0,0)"
	dx        int
	dy        int
	mx        int
	my        int

	// Default format
	format = func(t *Tree) string {
		return t.val
	}
)

// Tree position
type Point struct {
	X int
	Y int
}

func (p *Point) String() string {
	return fmt.Sprintf("(x:%d, y:%d)", p.X, p.Y)
}

func NewPoint() *Point {
	return &Point{0, 0}
}

// Tree
type Tree struct {
	*Point        // t position
	freq   int    // frequency of occurrence
	val    string // content
	left   *Tree  // left child
	right  *Tree  // right child
}

func (n *Tree) String() string {
	return fmt.Sprintf("%s %s", n.val, n.Point)
}

func NewTree() *Tree {
	return &Tree{NewPoint(), 1, "", nil, nil}
}

// Compute the tree height
func Height(t *Tree) int {
	if t == nil {
		return 0
	}

	hl := Height(t.left)
	hr := Height(t.right)
	max := hl
	if hl < hr {
		max = hr
	}

	return max + 1
}

// Display tree in pre order
func Prefixe(t *Tree) {
	if t == nil {
		return
	}
	fmt.Print(t, " ")
	Prefixe(t.left)
	Prefixe(t.right)
}

// Display tree in order
func Infixe(t *Tree) {
	if t == nil {
		return
	}
	Infixe(t.left)
	fmt.Print(t, " ")
	Infixe(t.right)
}

// Display tree in post order
func Postfixe(t *Tree) {
	if t == nil {
		return
	}
	Postfixe(t.left)
	Postfixe(t.right)
	fmt.Print(t, " ")
}

// Insert a value in tree
func Insert(t *Tree, val string) *Tree {
	if t == nil {
		t := NewTree()
		t.val = val
		return t
	}

	if t.val == val {
		t.freq++
	} else if t.val > val {
		t.left = Insert(t.left, val)
	} else {
		t.right = Insert(t.right, val)
	}
	return t
}

// Search a value in tree
func Search(t *Tree, val string) bool {

	res := false

	if t == nil {
		return res
	} else if t.val == val {
		res = true
	} else if t.val > val {
		return Search(t.left, val)
	} else {
		return Search(t.right, val)
	}

	return res
}

// Compute the position of each sub trees in the tree
// Full recursive version
func Position(t *Tree, x int, y int) int {

	if t.left != nil {
		x = Position(t.left, x, y+1)
	}

	t.X = x
	t.Y = y

	x = x + 1

	if t.right != nil {
		x = Position(t.right, x, y+1)
	}
	return x
}

// Compute the position of each sub trees in the tree
// Half recursive version
func Position2(t *Tree, x int, y int) int {

	for t != nil {
		if t.left != nil {
			x = Position(t.left, x, y+1)
		}

		t.X = x
		t.Y = y

		x = x + 1
		y = y + 1
		t = t.right
	}

	return x
}

// Apply transformation
func Transform(x int, y int) (int, int) {
	x1 := dx*x + mx
	y1 := dy*y + my

	return x1, y1
}

// Drawing the sub trees in SVG
func Draw(t *Tree, canvas *svg.SVG) {
	if t == nil {
		return
	}

	h, _ := strconv.Atoi(fontSize)
	x1, y1 := Transform(t.X, t.Y)

	if t.left != nil {
		left := t.left
		x2, y2 := Transform(left.X, left.Y)
		// log.Printf("x1:%d y1:%d x2:%d y2:%d %s\n", x1, y1, x2, y2, lineStyle)

		canvas.Line(x1, y1+(h/2), x2, y2-h, lineStyle)
		Draw(left, canvas)
	}

	if t.right != nil {
		right := t.right
		x2, y2 := Transform(right.X, right.Y)
		// log.Printf("x1:%d y1:%d x2:%d y2:%d %s\n", x1, y1, x2, y2, lineStyle)

		canvas.Line(x1, y1+(h/2), x2, y2-h, lineStyle)
		Draw(right, canvas)
	}

	canvas.Text(x1, y1, format(t), textStyle)
}

// Drawing the tree in  SVG
func Display(t *Tree, w io.Writer) {

	tWidth := Position2(t, 0, 0)
	tHeight := Height(t)

	dx = canvasWidth / tWidth
	dy = canvasHeight / tHeight
	mx = dx / 2
	my = dy / 2

	// fmt.Printf("w:%d h:%d", tWidth, tHeight)

	canvas := svg.New(w)
	canvas.Start(canvasWidth, canvasHeight)
	canvas.Rect(0, 0, canvasWidth, canvasHeight, rectStyle)

	Draw(t, canvas)
	canvas.End()
}

func main() {

	display := flag.String("d", "", "-d=[f,p] to display frequencies (f) or positions (p)")
	output := flag.String("o", "web", "-o=[web,stdout] output on webserver (web) or stdout (stdout)")

	flag.Parse()

	if len(os.Args[1+flag.NFlag():]) == 0 {
		log.Println("You must enter somme words ...")
		os.Exit(1)
	}

	if *display == "p" {
		// Position format
		format = func(t *Tree) string {
			return t.String()
		}
	} else if *display == "f" {
		// Frequence format
		format = func(t *Tree) string {
			return fmt.Sprintf("%s (f:%d)", t.val, t.freq)
		}
	}

	// Create the t
	var t *Tree

	for _, v := range os.Args[1+flag.NFlag():] {
		t = Insert(t, v)
	}

	// Send result to stdout
	if *output == "stdout" {
		Display(t, os.Stdout)

	} else {

		// Display the tree on Web browser
		s := ""
		buf := bytes.NewBufferString(s)
		Display(t, buf)

		// Send the output to the client
		http.Handle("/", http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "image/svg+xml")
				w.Write(buf.Bytes())
			}))
		err := http.ListenAndServe(":8000", nil)
		if err != nil {
			log.Fatal("ListenAndServe:", err)
		}

	}

}

package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"math"
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
	dx        float64
	dy        float64
	mx        float64
	my        float64

	// Frequence format
	format = func(t *Tree) string {
		return fmt.Sprintf("%s (f:%d)", t.val, t.freq)
	}

	// Position format
	// format = func(t *Tree) string {
	// 	return t.String()
	// }
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
func Height(t *Tree) float64 {
	if t == nil {
		return 0
	}
	return math.Max(Height(t.left), Height(t.right)) + float64(1)
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

// Drawing the sub trees in SVG
func Draw(t *Tree, canvas *svg.SVG) {
	if t == nil {
		return
	}

	h, _ := strconv.Atoi(fontSize)
	x1 := int(dx*float64(t.X) + mx)
	y1 := int(dy*float64(t.Y) + my)

	if t.left != nil {
		left := t.left
		x2 := int(dx*float64(left.X) + mx)
		y2 := int(dy*float64(left.Y) + my)
		// log.Printf("x1:%d y1:%d x2:%d y2:%d %s\n", x1, y1, x2, y2, lineStyle)

		canvas.Line(x1, y1+(h/2), x2, y2-h, lineStyle)
		Draw(left, canvas)
	}

	if t.right != nil {
		right := t.right
		x2 := int(dx*float64(right.X) + mx)
		y2 := int(dy*float64(right.Y) + my)
		// log.Printf("x1:%d y1:%d x2:%d y2:%d %s\n", x1, y1, x2, y2, lineStyle)

		canvas.Line(x1, y1+(h/2), x2, y2-h, lineStyle)
		Draw(right, canvas)
	}

	canvas.Text(x1, y1, format(t), textStyle)
}

// Drawing the tree in  SVG
func Display(t *Tree, w io.Writer) {

	tWidth := Position(t, 0, 0)
	tHeight := Height(t)

	canvas := svg.New(w)
	canvas.Start(canvasWidth, canvasHeight)
	canvas.Rect(0, 0, canvasWidth, canvasHeight, rectStyle)

	dx = float64(canvasWidth) / float64(tWidth)
	dy = float64(canvasHeight) / float64(tHeight)
	mx = dx / 2
	my = dy / 2

	Draw(t, canvas)
	canvas.End()
}

func main() {

	if len(os.Args[1:]) == 0 {
		log.Println("You must enter somme words ...")
		os.Exit(1)
	}

	// Create the t
	var t *Tree

	for _, v := range os.Args[1:] {
		t = Insert(t, v)
	}

	// Display the t
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

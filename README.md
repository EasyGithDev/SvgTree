# SvgTree
Display a binary tree with SVG

The program calculates the frequency of appearance of the words contained in a text.
A binary search tree containing the words of the text is made.
Then, the position of the nodes is calculated in a reference.
The tree is drawn using SVG format.
Finally, the SVG is sent to the browser to display.

## Install

Select a folder

```sh
cd myfloder
```

Clone the project
```sh
git clone git@github.com:EasyGithDev/SvgTree.git svgtree
```

Install the depencies

```sh
cd svgtree
go mod init
go get github.com/ajstarks/svgo
```

## Run

```sh
go run main.go this is time to said hello world one more time
```

## View

Open a web browser and enter the URL :

http://localhost:8000/
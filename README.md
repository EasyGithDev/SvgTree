# SvgTree
Display a binary tree with SVG

Here is a complete GO program which reads a document and produces an alphabetized list of words found therein together with the number of occurrences of each word. 

The method keeps a binary tree of words such that the left descendant tree for each word has all the words lexicographically smaller than the given word, and the right descendant has all the larger words.

Both the insertion and the drawing routine are recursive. 

Finally, the program send the SVG associated with the tree to the browser to display.

## Install

Select or create a folder :

```sh
cd myfloder
```

Clone the project into your selected folder :

```sh
git clone git@github.com:EasyGithDev/SvgTree.git svgtree
```

Install the depencies to work with SVG :

```sh
cd svgtree
go mod init
go get github.com/ajstarks/svgo
```

## Run

You may execute the program with a short text as parameter :

```sh
go run main.go this is time to said hello world one more time
```

![alt text](../assets/tree.svg?raw=true)

If you want display the words frequencies use the -d option like this : 

```sh
go run main.go -d=f this is time to said hello world one more time
```

![alt text](../assets/tree-f.svg?raw=true)


If you want display the node position use the -d option like this : 

```sh
go run main.go -d=p this is time to said hello world one more time
```

![alt text](../assets/tree-p.svg?raw=true)

## Display the result

Open a web browser and enter the URL :

http://localhost:8000/

## Write the result

You can choose to generate a SGV file to save the result.

You must change the output like this :

```sh
go run main.go -o=stdout this is time to said hello world one more time > tree.svg
```

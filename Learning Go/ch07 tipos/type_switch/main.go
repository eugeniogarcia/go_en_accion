package main

import (
	"errors"
	"fmt"
)

// define un tipo que representa un nodo en el árbol de análisis
type treeNode struct {
	val    treeVal
	lchild *treeNode
	rchild *treeNode
}

// interface...
type treeVal interface {
	isToken()
}

// dos tipos que implementan la interfaz treeVal
type number int

func (number) isToken() {}

type operator func(int, int) int

func (operator) isToken() {}
func (o operator) process(n1, n2 int) int {
	return o(n1, n2)
}

// mapa de operadores. Cada key es una funcion anonima que implementa el tipo operator
var operators = map[string]operator{
	"+": operator(func(n1, n2 int) int { //es redundante poner operator, pero lo hago para que se vea claro que es del tipo operator; En los otros casos no lo pongo
		return n1 + n2
	}),
	"-": func(n1, n2 int) int {
		return n1 - n2
	},
	"*": func(n1, n2 int) int {
		return n1 * n2
	},
	"/": func(n1, n2 int) int {
		return n1 / n2
	},
}

func walkTree(t *treeNode) (int, error) {
	//t.val es un interface. Obtenemos el tipo concreto de t.val usando type switch. Podemos contrastarlo con cualquier tipo, o interface; En este ejemplo usamos solo tipos
	switch val := t.val.(type) {
	case nil:
		return 0, errors.New("invalid expression")
	case number:
		return int(val), nil
	case operator:
		left, err := walkTree(t.lchild)
		if err != nil {
			return 0, err
		}
		right, err := walkTree(t.rchild)
		if err != nil {
			return 0, err
		}
		return val.process(left, right), nil
	default:
		return 0, errors.New("unknown node type")
	}
}

func parse(s string) (*treeNode, error) {
	//devuelve la direccion - puntero - de un nodo. Tiene harcodeado el parser del string "5*10+20"
	return &treeNode{
		val: operators["+"],
		lchild: &treeNode{
			val:    operators["*"],
			lchild: &treeNode{val: number(5)},
			rchild: &treeNode{val: number(10)},
		},
		rchild: &treeNode{val: number(20)},
	}, nil
}

func main() {
	parseTree, err := parse("5*10+20")
	if err != nil {
		panic(err)
	}
	result, err := walkTree(parseTree)
	fmt.Println(result, err)
}

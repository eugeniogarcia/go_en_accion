package main

import "fmt"

type ImpossiblePrintableInt interface {
	int | ~int8
	String() string
}

type ImpossibleStruct[T ImpossiblePrintableInt] struct {
	val T
}

type MyInt int

func (mi MyInt) String() string {
	return fmt.Sprint(mi)
}

type MyInt8 int8

func (mi MyInt8) String() string {
	return fmt.Sprint(mi)
}

func main() {
	s1 := ImpossibleStruct[int]{10}    //no sirve porque aunque implementa int, no implementa String(). No se cumple la restricción del generic
	s2 := ImpossibleStruct[MyInt]{10}  //no sirve porque aunque String(), MyInt no esta entre los tipos que hemos incluido en la definición del interface . No se cumple la restricción del generic
	s3 := ImpossibleStruct[MyInt8]{10} //sirve porque al poner ~ en la definición del tipo (~int8), lo que decimos es que admitimos que el parametro sea int8 o cualquier tipo que use int8 como derivado

	fmt.Println(s1, s2, s3)
}

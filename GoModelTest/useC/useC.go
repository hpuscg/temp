/*
#Time      :  2020/9/3 3:13 下午
#Author    :  chuangangshen@deepglint.com
#File      :  useC.go
#Software  :  GoLand
*/
package main

//#include "useC.c"
/*
#include <stdio.h>

void sayYes(){
    printf("yes!\n");
}
*/
import "C"

func main() {
	C.sayHello()
	C.sayYes()
}

func MapLock() {
	// sync.Map{}
}

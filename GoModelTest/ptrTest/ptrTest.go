package main

func main() {
	type b struct {

	}
	a := make([]*b, 20)
	s := new(b)
	a[1] = s
}

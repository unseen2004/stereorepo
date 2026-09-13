package main

import "fmt"

func main() {
	// I AM NOT DONE
	x := 1
	{
		x := 2
		fmt.Println(x)
	}
	fmt.Println(x)
}

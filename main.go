package main

import "os"

func main() {

	a := App{}

	a.Initialize()

	a.Run(":8080")
}
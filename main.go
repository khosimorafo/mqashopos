package main

import "os"

func main() {

	port := os.Getenv("8080")

	a := App{}

	a.Initialize()

	a.Run(port)
}
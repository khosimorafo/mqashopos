package main

import "os"

func main() {

	port := os.Getenv("PORT")

	a := App{}

	a.Initialize()

	a.Run(port)
}
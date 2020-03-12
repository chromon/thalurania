package main

import "chalurania/comet"

func main() {
	s := comet.NewServer("s1")
	s.Serve()
}
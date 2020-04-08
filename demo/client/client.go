package main

import (
	"bufio"
	"chalurania/demo/client/arguments"
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Println("welcome to thalurania im.")

	scanner := bufio.NewScanner(os.Stdin)
	var args string

	for {

		c := arguments.NewCommand("tim")
		c.CommandInit()

		fmt.Print("~ ")
		scanner.Scan()
		if err := scanner.Err(); err != nil {
			_, err = fmt.Fprintln(os.Stderr, "error:", err)
		}
		fmt.Printf("text: %q\r\n", scanner.Text())
		args = scanner.Text()

		c.ParseCommand(strings.Split(args, " ")[1:])
		c.VisitCommand()

		fmt.Println(c.FlagSet.Lookup("p"))
		fmt.Println(c.FlagSet.Lookup("s"))
		fmt.Println(c.FlagSet.Lookup("b"))

		fmt.Println(c.CommandMap["p"])
	}
}

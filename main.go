package main

import (
	"fmt"
	"os"

	"dial/cmd"
)

func main() {
	if len(os.Args) > 1 {
		cmd.Execute()
		return
	}
	launchGUI()
}

func launchGUI() {
	fmt.Println("GUI not implemented yet; run dial with a subcommand to use the CLI.")
}

package main

import (
    "os"

    "dail/cmd"
)

func main() {
    if len(os.Args) > 1 {
        cmd.Execute()
	return
    }
    launchGUI()
}

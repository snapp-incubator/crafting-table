package assets

import "fmt"

const asciiArt = `
 ██████ ██████   █████  ███████ ████████ ██ ███    ██  ██████      ████████  █████  ██████  ██      ███████ 
██      ██   ██ ██   ██ ██         ██    ██ ████   ██ ██              ██    ██   ██ ██   ██ ██      ██      
██      ██████  ███████ █████      ██    ██ ██ ██  ██ ██   ███        ██    ███████ ██████  ██      █████   
██      ██   ██ ██   ██ ██         ██    ██ ██  ██ ██ ██    ██        ██    ██   ██ ██   ██ ██      ██      
 ██████ ██   ██ ██   ██ ██         ██    ██ ██   ████  ██████         ██    ██   ██ ██████  ███████ ███████ 
`

func PrintAsciiArt() {
	fmt.Print(asciiArt)
}

package banner

import "fmt"

const asciiArt = `
██╗   ██╗██╗██████╗ ███████╗ ██████╗██╗  ██╗███████╗ ██████╗██╗  ██╗
██║   ██║██║██╔══██╗██╔════╝██╔════╝██║  ██║██╔════╝██╔════╝██║ ██╔╝
██║   ██║██║██████╔╝█████╗  ██║     ███████║█████╗  ██║     █████╔╝ 
╚██╗ ██╔╝██║██╔══██╗██╔══╝  ██║     ██╔══██║██╔══╝  ██║     ██╔═██╗ 
 ╚████╔╝ ██║██████╔╝███████╗╚██████╗██║  ██║███████╗╚██████╗██║  ██╗
  ╚═══╝  ╚═╝╚═════╝ ╚══════╝ ╚═════╝╚═╝  ╚═╝╚══════╝ ╚═════╝╚═╝  ╚═╝`

const tagline = "vibecheck is a lightweight, cross-platform command line AI-tool which\nautomatically generates meaningful and consistent Git Commit Messages\nby analyzing your code changes — ship faster with vibecheck"

// Print writes the vibecheck ASCII art to stdout.
func Print() {
	fmt.Println(asciiArt)
	fmt.Println()
	fmt.Println(tagline)
}

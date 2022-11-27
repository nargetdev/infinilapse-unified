package color

import "fmt"

func PrintYellow(say string) {
	yellow := fmt.Sprintf("\x1b[%dm%s\x1b[0m", 33, say)
	fmt.Println(yellow)
}

func PrintMagenta(say string) {
	cyan := fmt.Sprintf("\x1b[%dm%s\x1b[0m", 35, say)
	fmt.Println(cyan)
}

func PrintCyanBold(say string) {
	cyan := fmt.Sprintf("\x1b[1m\x1b[%dm%s\x1b[0m\x1b[22m", 36, say)
	fmt.Println(cyan)
}
func PrintCyan(say string) {
	cyan := fmt.Sprintf("\x1b[%dm%s\x1b[0m", 36, say)
	fmt.Println(cyan)
}

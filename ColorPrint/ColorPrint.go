package colorprint

import (
	"fmt"
	"github.com/fatih/color"
	"time"
)

func PrintYellow(username string, message string) {
	now := time.Now()
	color.Set(color.FgYellow)
	fmt.Println("["+now.Format("2006-01-02 15:04:05")+"]", username+":", message)
}
func PrintPurple(username string, message string) {
	now := time.Now()
	color.Set(color.FgMagenta)
	fmt.Println("["+now.Format("2006-01-02 15:04:05")+"]", username+":", message)
}
func PrintRed(username string, message string) {
	now := time.Now()
	color.Set(color.FgRed)
	fmt.Println("["+now.Format("2006-01-02 15:04:05")+"]", username+":", message)
}
func PrintGreen(username string, message string) {
	now := time.Now()
	color.Set(color.FgGreen)
	fmt.Println("["+now.Format("2006-01-02 15:04:05")+"]", username+":", message)
}
func PrintBlue(username string, message string) {
	now := time.Now()
	color.Set(color.FgBlue)
	fmt.Println("["+now.Format("2006-01-02 15:04:05")+"]", username+":", message)
}
func PrintWhite(username string, message string) {
	now := time.Now()
	color.Set(color.FgWhite)
	fmt.Println("["+now.Format("2006-01-02 15:04:05")+"]", username+":", message)
}
func PrintGold(username string, message string) {
	now := time.Now()
	color.Set(color.FgHiYellow)
	fmt.Println("["+now.Format("2006-01-02 15:04:05")+"]", username+":", message)
}
func PrintCyan(username string, message string) {
	now := time.Now()
	color.Set(color.FgCyan)
	fmt.Println("["+now.Format("2006-01-02 15:04:05")+"]", username+":", message)
}

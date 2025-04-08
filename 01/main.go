package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

func main() {
	fmt.Print(ansi.RequestCursorPositionReport)

	style := lipgloss.NewStyle().Foreground(lipgloss.Yellow)
	lipgloss.Println(style.Render("Hello Gophercon!"))

	fmt.Print(ansi.KittyKeyboard(ansi.KittyAllFlags, 1))
	fmt.Print(ansi.SetWindowTitle("Hello Gophercon!"))
	fmt.Print(ansi.RequestTermcap("cols"))
}

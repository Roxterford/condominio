package logger

import (
	"fmt"
)

// Función para imprimir un mensaje de error
func PrintError(err error) {
	fmt.Printf("\x1b[31mERROR: %s\x1b[0m\n", err.Error())
}

// Función para imprimir un mensaje de advertencia
func PrintWarning(message string) {
	fmt.Printf("\x1b[33mWARNING: %s\x1b[0m\n", message)
}

// Función para imprimir un mensaje informativo
func PrintInfo(message string) {
	fmt.Printf("\x1b[32mINFO: %s\x1b[0m\n", message)
}

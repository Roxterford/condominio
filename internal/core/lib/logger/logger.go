package logger

import (
	"fmt"

	"github.com/Sanaruca/condominio/internal/core/envirotment"
)

var debug bool

// Función para imprimir un mensaje de error
func Error(err error) {
	fmt.Printf("\x1b[31mERROR: %s\x1b[0m\n", err.Error())
}

// Función para imprimir un mensaje de advertencia
func Warning(message string, msgArgs ...any) {
	coloredFormat := "\x1b[33mWARNING: " + message + "\x1b[0m\n"
	fmt.Printf(coloredFormat, msgArgs...)
}

// Función para imprimir un mensaje informativo
func Info(message string, msgArgs ...any) {
	coloredFormat := "\x1b[32mINFO: " + message + "\x1b[0m\n"
	fmt.Printf(coloredFormat, msgArgs...)
}

// Debug imprime un mensaje formateado si el modo debug está activo.
func Debug(message string, args ...any) {
	if !debug {
		return
	}
	coloredFormat := "\x1b[34mDEBUG: " + message + "\x1b[0m\n"
	fmt.Printf(coloredFormat, args...)
}

func init() {
	debug = envirotment.GetAppEnv() == envirotment.Dev
}

package utils

// Helper de soporte para limpiar la sintaxis de punteros
func SafeStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

package graph

import "github.com/Sanaruca/condominio/internal/core/common/mes"

func mesDeref(m *mes.Mes) mes.Mes {
	if m == nil {
		return 0
	}
	return *m
}

func stringDeref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

package session

import "github.com/Sanaruca/condominio/internal/usuarios"

type Session interface {
	Usuario() *usuarios.UsuarioPayload
}

type session struct {
	usuario *usuarios.UsuarioPayload
}

func New(usuario *usuarios.UsuarioPayload) Session {
	return &session{
		usuario: usuario,
	}
}

func (s *session) Usuario() *usuarios.UsuarioPayload {
	return s.usuario
}

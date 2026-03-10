package session

type CredencialDeUsuario struct {
	ID    string
	Email string
}

type Session interface {
	Usuario() *CredencialDeUsuario
}
type SessionWithUser interface {
	Usuario() CredencialDeUsuario
}

type session struct {
	usuario *CredencialDeUsuario
}

func New(usuario *CredencialDeUsuario) Session {
	return &session{
		usuario: usuario,
	}
}

func (s *session) Usuario() *CredencialDeUsuario {
	return s.usuario
}

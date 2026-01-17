package context

import (
	"context"
	"time"

	// "time"

	"github.com/Sanaruca/condominio/internal/core/session"
	"gorm.io/gorm"
)

const (
	BASE_CONTEXT_KEY = "base_context"
)

type BaseContext struct {
	context context.Context
	DB      *gorm.DB
	session session.Session
}
type AppContext = BaseContext

func New(ctx context.Context, db *gorm.DB) BaseContext {
	return BaseContext{
		context: ctx,
		DB:      db,
	}
}

func (b BaseContext) WithSession(session session.Session) BaseContext {
	return BaseContext{
		context: b.context,
		DB:      b.DB,
		session: session,
	}
}

///////////////////////////////////////////////////////////

func (b BaseContext) Context() context.Context {
	return b.context
}

func (b BaseContext) Session() session.Session {
	return b.session
}

func (b BaseContext) Deadline() (time.Time, bool) {
	return b.context.Deadline()
}
func (b BaseContext) Value(key any) any {
	return b.context.Value(key)
}

func (b BaseContext) Done() <-chan struct{} {
	return b.context.Done()
}

func (b BaseContext) Err() error {
	return b.context.Err()
}

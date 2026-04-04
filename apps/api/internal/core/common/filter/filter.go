package filter

// TODO: Evaluar retorno de core.Error en lugar de error cuando se normalice la estructura de errores del proyecto
// para mantener consistencia con el resto del codebase

const DEFAULT_MAX_DEPTH = 10 // Límite por defecto para evitar estructuras maliciosas y asegurar performance

type Filter[T Filterable] struct {
	raw    map[string]any
	clause Clause
}

func NewFilter[T Filterable](input map[string]any) *Filter[T] {
	return &Filter[T]{raw: input}
}

func (f *Filter[T]) Build() (Clause, error) {
	if f == nil || len(f.raw) == 0 {
		return nil, nil
	}
	return Build[T](f.raw)
}

func (f *Filter[T]) Clause() Clause {
	return f.clause
}

func (f *Filter[T]) Validate() error {
	clause, err := f.Build()
	if err != nil {
		return err
	}
	f.clause = clause
	return nil
}

// Config mantiene la configuración de la construcción del filtro
type Config struct {
	MaxDepth int
}

// Option define una función que modifica la configuración
type Option func(*Config)

// WithMaxDepth permite cambiar el límite de recursividad
func WithMaxDepth(depth int) Option {
	return func(c *Config) {
		c.MaxDepth = depth
	}
}

type Filterable interface {
	FilterSpec() Spec
}

// Build construye y valida un árbol de filtros
func Build[T Filterable](input any, opts ...Option) (Clause, error) {

	// Configuración por defecto
	config := &Config{
		MaxDepth: DEFAULT_MAX_DEPTH,
	}
	// Aplicar opciones del usuario
	for _, opt := range opts {
		opt(config)
	}

	// 1. Obtener Spec del genérico
	var entity T
	spec := entity.FilterSpec()

	// 2. Parsear (Raw -> AST)
	rootClause, err := parseWithDepth(input, 0, config.MaxDepth)
	if err != nil {
		return nil, err
	}

	// 3. Validar (AST + Spec -> Checked AST)
	validator := &Validator{FilterSpec: spec}
	if err := validator.Validate(rootClause); err != nil {
		return nil, err
	}

	return rootClause, nil
}

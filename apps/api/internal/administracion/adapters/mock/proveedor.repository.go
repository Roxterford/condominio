package mock

import (
	"context"
	"sync"

	"github.com/Sanaruca/condominio/internal/administracion/models/proveedor"
	"github.com/Sanaruca/condominio/internal/administracion/models/proveedor/tipodeproveedor"
	"github.com/Sanaruca/condominio/internal/core"
	"github.com/lucsky/cuid"
)

type MockProveedorRepository struct {
	mu          sync.Mutex
	proveedores map[string]*proveedor.Proveedor
	secuencia   int
	factory     *proveedor.ProveedorFactory
}

func NewMockProveedorRepository(factory *proveedor.ProveedorFactory) *MockProveedorRepository {
	return &MockProveedorRepository{
		proveedores: make(map[string]*proveedor.Proveedor),
		factory:     factory,
	}
}

func (m *MockProveedorRepository) Guardar(
	ctx context.Context,
	proveedor *proveedor.Proveedor,
) core.Error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.proveedores[proveedor.ID()] = proveedor
	return nil
}

func (m *MockProveedorRepository) ObtenerPorID(
	ctx context.Context,
	id string,
) (*proveedor.Proveedor, core.Error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if p, ok := m.proveedores[id]; ok {
		return p, nil
	}
	return nil, proveedor.ErrProveedorNoEncontrado
}

func (m *MockProveedorRepository) ObtenerTodos(
	ctx context.Context,
) ([]*proveedor.Proveedor, core.Error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	result := make([]*proveedor.Proveedor, 0, len(m.proveedores))
	for _, p := range m.proveedores {
		result = append(result, p)
	}
	return result, nil
}

func (m *MockProveedorRepository) Actualizar(
	ctx context.Context,
	p *proveedor.Proveedor,
) core.Error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.proveedores[p.ID()]; !ok {
		return proveedor.ErrProveedorNoEncontrado
	}
	m.proveedores[p.ID()] = p
	return nil
}

func (m *MockProveedorRepository) ExistePorRif(ctx context.Context, rif string) (bool, core.Error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, p := range m.proveedores {
		if p.Rif().String() == rif {
			return true, nil
		}
	}
	return false, nil
}

func (m *MockProveedorRepository) ExistePorID(ctx context.Context, id string) (bool, core.Error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	_, ok := m.proveedores[id]
	return ok, nil
}

type ProveedorBuilder struct {
	id        string
	rif       string
	nombre    string
	tipo      tipodeproveedor.TipoDeProveedor
	email     string
	telefono  string
	direccion *string
}

func NewProveedorBuilder() *ProveedorBuilder {
	return &ProveedorBuilder{
		id:       cuid.New(),
		rif:      "J-123456789",
		nombre:   "Proveedor Test",
		tipo:     tipodeproveedor.TipoDeProveedor(1),
		email:    "test@proveedor.com",
		telefono: "+584141234567",
	}
}

func (b *ProveedorBuilder) WithID(id string) *ProveedorBuilder {
	b.id = id
	return b
}

func (b *ProveedorBuilder) WithRif(rif string) *ProveedorBuilder {
	b.rif = rif
	return b
}

func (b *ProveedorBuilder) WithNombre(nombre string) *ProveedorBuilder {
	b.nombre = nombre
	return b
}

func (b *ProveedorBuilder) WithTipo(tipo tipodeproveedor.TipoDeProveedor) *ProveedorBuilder {
	b.tipo = tipo
	return b
}

func (b *ProveedorBuilder) WithEmail(email string) *ProveedorBuilder {
	b.email = email
	return b
}

func (b *ProveedorBuilder) WithTelefono(telefono string) *ProveedorBuilder {
	b.telefono = telefono
	return b
}

func (b *ProveedorBuilder) WithDireccion(direccion string) *ProveedorBuilder {
	b.direccion = &direccion
	return b
}

func (b *ProveedorBuilder) Build(factory *proveedor.ProveedorFactory) *proveedor.Proveedor {

	proveedor, err := factory.Nuevo(
		b.rif,
		b.nombre,
		b.tipo,
		b.email,
		b.telefono,
		b.direccion,
	)
	if err != nil {
		panic("error building proveedor for testing: " + err.Error())
	}
	return proveedor
}

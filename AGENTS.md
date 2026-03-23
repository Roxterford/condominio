## Convenciones de Commits

- Formato: [Conventional Commits](https://www.conventionalcommits.org/) con emojis
- Estructura: `tipo: 🔖 descripción`
- Tipos permitidos: `feat`, `fix`, `build`, `refactor`, `docs`, `test`, `chore`

Ejemplos:

```
feat: 🎨 agregar registro de cuotas
fix: 🐛 corregir validación de formulario
build: 🧩 configurar moonrepo
```

## Convenciones de desarrollo

### Nombramiento General

- **Variables**: `snake_case`
- **Constantes**: `UPPER_SNAKE_CASE`
- **Archivos**: `snake_case`
- **Directorios**: `snake_case`

### Testing y Estructura (DDD & Clean Architecture)

Para mantener la integridad de las capas y asegurar una suite de pruebas robusta, se aplican las siguientes reglas:

1. **Cercanía del Test (Locality of Testing)**:

- Los archivos de prueba deben residir en el mismo directorio que el código que auditan. Esto facilita el mantenimiento en arquitecturas multicapa.
- **Nomenclatura**: El archivo de test debe heredar el nombre del archivo original, seguido del sufijo correspondiente al ecosistema:
- **Go**: `nombre_archivo_test.go`
- **Otros (Node/Python/etc)**: `nombre_archivo.test.ext` o `nombre_archivo_test.ext` según la convención estándar del framework.

2. **Pruebas por Capa (Clean Architecture)**:

- **Domain Tests**: Deben ser pruebas unitarias puras, sin dependencias externas (moteando interfaces de infraestructura).
- **Application/Use Case Tests**: Enfocados en el flujo de la lógica de negocio.
- **Infrastructure Tests**: Pruebas de integración para validar implementaciones reales (ej. Repositorios de base de datos o clientes API).

3. **Encapsulamiento de Tests en Go**:

- Se recomienda el uso del sufijo `_test` en el nombre del paquete (ej. `package service_test`) para forzar el testing de "caja negra". Esto asegura que solo se pruebe la API pública, respetando el encapsulamiento definido en DDD.

## 🏗️ Especificación de Factorías: Construcción Dual

Para garantizar la estabilidad del sistema a largo plazo, todas las factorías de entidades deben implementar una separación clara entre la **Creación de Negocio** y la **Rehidratación de Persistencia**.

### 1. El Método `Nuevo`/`Nueva` (Creación / Escritura)

Es el guardián de la integridad del sistema. Se utiliza cuando un usuario o proceso intenta introducir datos nuevos.

- **Entrada:** Tipos primitivos (`string`, `int`, etc.).
- **Responsabilidad:**
- Validar los datos contra las políticas y reglas de negocio **vigentes**.
- Transformar los primitivos en **Value Objects** usando constructores estrictos (que retornan error).
- Generar Identificadores (`ID`) y marcas de tiempo actuales.

- **Resultado:** `(Entidad, error de dominio)`. Si el dato es inválido, la operación se rechaza.

### 2. El Método `Assemble` (Rehidratación / Lectura)

Es el método utilizado por los Repositorios para reconstruir entidades desde la base de datos.

- **Entrada:** Tipos primitivos recuperados de la persistencia.
- **Responsabilidad:**
- Reconstruir la entidad tal cual existe en la base de datos.
- Respetar los IDs y fechas originales guardados en el disco.

- **Resultado:** `Entidad`. Este método es **infalible**; no debe retornar errores de validación de negocio para asegurar que el sistema siempre pueda leer su propia historia.

### 3. Lógica de Actualización (Update)

En los flujos de edición, el sistema protege la integridad sin bloquear datos antiguos:

- Si un campo **no ha sido modificado** por el usuario, se mantiene el valor rehidratado por `Assemble`.
- Si un campo **es modificado**, el Caso de Uso debe invocar la validación estricta (`NewX`) para asegurar que el cambio cumpla la ley actual.

## Stack Tecnológico

### Arquitectura General

- **Monorepo**: [moonrepo](https://moonrepo.dev/) para gestión de proyectos
- **Gestor de Paquetes**: Bun como principal, con soporte para npm
- **Control de Versiones**: Git con Lefthook para hooks
- **Base de Datos**: SQLite con Prisma ORM y Drizzle ORM (dual)

### Frontend (apps/panel)

- **Framework**: SvelteKit 2.50.2
- **Lenguaje**: TypeScript 5.9.3
- **Estilos**: TailwindCSS 4.1.18 + Flowbite Svelte
- **Componentes**: Storybook 10.2.17
- **Testing**: Vitest + Playwright
- **Autenticación**: Better Auth 1.4.21
- **Formularios**: SvelteKit Superforms + Valibot
- **Build Tool**: Vite 7.3.1

### Backend (apps/api)

- **Lenguaje**: Go 1.25.5
- **GraphQL**: gqlgen 0.17.86
- **Base de Datos**: GORM v1.31.1 con SQLite
- **Autenticación**: JWT golang-jwt/jwt/v5
- **Testing**: Testify + Go SQL Mock

### Base de Datos

- **Motor**: SQLite
- **ORM Principal**: Prisma 7.2.0
- **ORM Secundario**: Drizzle ORM 0.45.1
- **Migraciones**: Prisma Migrate + Drizzle Kit
- **Seeders**: tsx para ejecución de seeds

### Herramientas de Desarrollo

- **Linting**: ESLint + Prettier
- **Formato**: Prettier con plugins Svelte y Tailwind
- **Git Hooks**: Lefthook
- **Validación de Commits**: Commitlint + Commitizen
- **IDE**: VS Code con configuración TypeScript

### Testing

- **Unitarios**: Vitest (frontend) + Go testing (backend)
- **Integración**: Storybook Test addon
- **E2E**: Playwright
- **Coverage**: Vitest Coverage v8

## Estructura del Proyecto

```
condominio/
├── apps/
│   ├── panel/          # Frontend SvelteKit
│   └── api/            # Backend GraphQL Go
├── prisma/             # Esquema y seeds de Prisma
├── generated/          # Cliente Prisma generado
├── .moon/              # Configuración moonrepo
├── .github/            # Templates de issues
└── apps/panel/
    ├── src/
    │   ├── lib/
    │   │   ├── server/
    │   │   │   ├── auth.ts        # Better Auth config
    │   │   │   └── db/            # Drizzle schemas
    │   │   └── components/        # Componentes UI
    │   └── routes/                # Páginas SvelteKit
    ├── .storybook/               # Config Storybook
    ├── e2e/                      # Tests Playwright
    └── stories/                  # Historias Storybook
```

## Scripts Principales

### Root Level

- `dev`: Inicia desarrollo de todos los proyectos
- `build`: Construye todos los proyectos
- `test`: Ejecuta todas las pruebas
- `lint`: Linting en todo el monorepo
- `format`: Formateo con Prettier
- `db`: Operaciones de base de datos

### Frontend (apps/panel)

- `dev`: Servidor de desarrollo SvelteKit
- `build`: Build de producción
- `storybook`: Servidor Storybook
- `test:unit`: Tests unitarios Vitest
- `test:e2e`: Tests E2E Playwright
- `db:*`: Operaciones Drizzle (push, generate, migrate, studio)
- `auth:schema`: Generar schema Better Auth

### Backend (apps/api)

- `cov`: Generar reporte de cobertura
- `check-unix`: Verificar compatibilidad Unix
- `check-golines`: Verificar instalación de golines

## Configuración Clave

### Variables de Entorno

- `SECRET_KEY`: Clave para autenticación
- `DATABASE_URL`: URL de base de datos SQLite
- `ORIGIN`: URL base para Better Auth
- `BETTER_AUTH_SECRET`: Secret para Better Auth

### Convenciones de Código

- **Indentación**: Tabs (Prettier)
- **Comillas**: Simple quotes
- **Ancho**: 100 caracteres
- **Semicolones**: No trailing commas
- **TypeScript**: Strict mode activado

## Flujo de Trabajo

### 1. Desarrollo

```bash
# Iniciar todos los servicios
bun dev

# O individualmente
moon run panel:dev    # Frontend
moon run api:dev      # Backend (si existe)
```

### 2. Commits

```bash
# Los commits usan Commitizen automáticamente
git commit  # Abrirá interfaz interactiva

# Formato esperado:
feat: 🎨 agregar nueva funcionalidad
fix: 🐛 corregir error
docs: 📝 actualizar documentación
```

### 3. Testing

```bash
# Todas las pruebas
moon run :test-all

# Frontend específico
moon run panel:test-unit
moon run panel:test-e2e
```

### 4. Base de Datos

```bash
# Prisma
bun prisma generate
bun prisma db push

# Drizzle
moon run panel:db-generate
moon run panel:db-push
moon run panel:db-studio
```

## Consideraciones Especiales

### Base de Datos Dual

- **Prisma**: Usado para esquema principal y seeds
- **Drizzle**: Usado para autenticación y operaciones específicas
- Ambos conectan a la misma base de datos SQLite

### Autenticación

- **Better Auth**: Configurado con adapter Drizzle
- **SvelteKit Integration**: Cookies automáticas
- **Email/Password**: Habilitado por defecto

### Internacionalización

- El proyecto está configurado principalmente en español
- Templates de GitHub en español
- Comentarios y documentación en español

### Monorepo Management

- **moonrepo**: Gestiona dependencias y tareas
- **Workspaces**: Configurado para apps/panel
- **Affected**: Testing solo en archivos modificados

## Recursos Adicionales

### Documentación

- [SvelteKit Docs](https://kit.svelte.dev/)
- [Prisma Docs](https://www.prisma.io/docs/)
- [moonrepo Docs](https://moonrepo.dev/docs)
- [Better Auth Docs](https://better-auth.com/docs)

### Templates GitHub

- Reporte de bugs: `.github/ISSUE_TEMPLATE/reporte-de-bug.md`
- Solicitudes: `.github/ISSUE_TEMPLATE/solicutud-de-funcionalidad.md`

### Configuración IDE

- VS Code con TypeScript SDK
- Extensión recomendada para Svelte
- Configuración de Prettier y ESLint integradas

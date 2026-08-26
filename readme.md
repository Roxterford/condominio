# 🏠 Condominio

gestión de condominios con monorepo moderno.

![Next.js](https://img.shields.io/badge/next%20js-000000?style=for-the-badge&logo=nextdotjs&logoColor=white)
![TypeScript](https://img.shields.io/badge/TypeScript-007ACC?style=for-the-badge&logo=typescript&logoColor=white)
![Go](https://img.shields.io/badge/GoLand-000000?style=for-the-badge&logo=goland&logoColor=white)
![GraphQL](https://img.shields.io/badge/GraphQl-E10098?style=for-the-badge&logo=graphql&logoColor=white)
![Tailwind](https://img.shields.io/badge/Tailwind_CSS-030712?style=for-the-badge&logo=tailwind-css&logoColor=1ac3ff)

<p align="center">
  <img src="https://skillicons.dev/icons?i=react,prisma,sqlite,bun,vitest" />
</p>

---

## ✨ Stack Tecnológico

### 🖥️ Frontend

| Tecnología      | Versión |
| --------------- | ------- |
| Next.js         | 16.2.0  |
| TypeScript      | 5.9.3   |
| TailwindCSS     | 4.x     |
| Storybook       | 10.2.17 |
| Vitest          | -       |
| Playwright      | -       |
| Better Auth     | 1.4.21  |
| React Hook Form | -       |
| Valibot         | -       |

### ⚙️ Backend

| Tecnología     | Versión |
| -------------- | ------- |
| Go             | 1.25.5  |
| gqlgen         | 0.17.86 |
| GORM           | 1.31.1  |
| golang-jwt/jwt | v5      |
| Testify        | -       |

### 🗄️ Base de Datos

| Tecnología  | Versión |
| ----------- | ------- |
| SQLite      | -       |
| Prisma      | 7.2.0   |
| Drizzle ORM | 0.45.1  |

### 🛠️ Herramientas

- **Monorepo**: moonrepo
- **Paquete**: Bun
- **Linting**: ESLint + Prettier
- **Git Hooks**: Lefthook
- **Commits**: Commitlint + Commitizen

---

## 📁 Estructura del Proyecto

```
condominio/
├── apps/
│   ├── panel/          # Frontend Next.js
│   └── api/            # Backend GraphQL Go
├── prisma/             # Esquema y seeds
├── generated/          # Cliente Prisma
├── .moon/              # Config moonrepo
├── .github/            # Templates issues
└── apps/panel/
    ├── src/
    │   ├── lib/
    │   │   ├── server/
    │   │   │   ├── auth.ts    # Better Auth
    │   │   │   └── db/        # Drizzle schemas
    │   │   └── components/   # Componentes UI
    │   └── app/               # Rutas Next.js
    ├── .storybook/
    ├── e2e/
    └── stories/
```

---

## 🚀 Scripts

### Root

| Comando      | Descripción         |
| ------------ | ------------------- |
| `bun dev`    | Iniciar desarrollo  |
| `bun build`  | Construir proyectos |
| `bun test`   | Ejecutar pruebas    |
| `bun lint`   | Linting total       |
| `bun format` | Formatear código    |
| `bun db`     | Operaciones BD      |

### Frontend

```bash
moon run panel:dev         # Desarrollo Next.js
moon run panel:build       # Build producción
moon run panel:storybook   # Servir Storybook
moon run panel:test-unit   # Tests Vitest
moon run panel:test-e2e    # Tests Playwright
moon run panel:db-push     # Drizzle push
moon run panel:db-studio   # Drizzle studio
```

### Backend

```bash
moon run api:dev           # Desarrollo Go
moon run api:cov           # Coverage
```

---

## ⚙️ Variables de Entorno

```env
SECRET_KEY=<clave>
DATABASE_URL=<url sqlite>
ORIGIN=<url base>
BETTER_AUTH_SECRET=<secret>
```

---

## 🔄 Flujo de Trabajo

### 1️⃣ Desarrollo

```bash
bun dev                    # Todos los servicios
moon run panel:dev        # Solo frontend
moon run api:dev          # Solo backend
```

### 2️⃣ Commits

```bash
git commit  # Interfaz interactiva Commitizen
```

### 3️⃣ Testing

```bash
moon run :test-all        # Todas las pruebas
moon run panel:test-unit  # Unitarios frontend
moon run panel:test-e2e   # E2E frontend
```

### 4️⃣ Base de Datos

```bash
# Prisma
bun prisma generate && bun prisma db push

# Drizzle
moon run panel:db-generate
moon run panel:db-push
moon run panel:db-studio
```

---

## 💡 Consideraciones Especiales

### 🗃️ Base de Datos Dual

- **Prisma**: Esquema principal + seeds
- **Drizzle**: Autenticación + operaciones específicas
- Ambos comparten la misma SQLite

### 🔐 Autenticación

- Better Auth con adapter Drizzle
- Cookies automáticas en Next.js
- Email/Password habilitado

### 🌍 Internacionalización

- Proyecto en español
- Templates GitHub en español

---

## 📚 Recursos

| Recurso     | Enlace                               |
| ----------- | ------------------------------------ |
| Next.js     | [docs](https://nextjs.org/docs)      |
| Prisma      | [docs](https://www.prisma.io/docs/)  |
| moonrepo    | [docs](https://moonrepo.dev/docs)    |
| Better Auth | [docs](https://better-auth.com/docs) |

### Templates GitHub

- `.github/ISSUE_TEMPLATE/reporte-de-bug.md`
- `.github/ISSUE_TEMPLATE/solicutud-de-funcionalidad.md`

---

## 📋 TODOS

- [x] Evaluar mover `./apps/api/sql` a `./prisma/sql`
- [x] Renombrar `Villas` → `Unidades` (sistema agnóstico)
- [ ] Arreglar la condicion in en los filtros dinamicos
- [ ] Limitar la cantidad de gastos que puede contener una Cuota, esto por motivos de simplicidad
- [ ] Detectar y solucionar problemas n + 1

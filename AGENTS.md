# AGENTS.md

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

## Stack

- Monorepo con [moonrepo](https://moonrepo.dev/)
- Frontend: SvelteKit (apps/panel)
- Backend: Supabase

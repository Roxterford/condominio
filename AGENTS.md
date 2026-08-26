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

### Estructura de Componentes Next.js

Para mantener una organización limpia y separable, los componentes se clasifican en:

1. **Componentes específicos de página** → `app/<ruta>/components/`
   - Solo son usados por esa página
   - Ej: `app/cuotas/registrar/components/registrar-cuota-form.tsx`
   - Importar con ruta relativa: `import { X } from "./components/x"`

2. **Componentes reutilizables** → `components/<nombre>/`
   - Usados en múltiples páginas
   - Ej: `components/reusable/data-table.tsx`
   - Importar con alias: `import { DataTable } from "@/components/reusable/data-table"`

3. **Componentes UI (shadcn)** → `components/ui/`
   - Exclusivo para componentes de shadcn/ui
   - Ej: `components/ui/button.tsx`
   - Importar con alias: `import { Button } from "@/components/ui/button"`

## Skills Disponibles

| Skill | Ubicación |
|-------|-----------|
| `conventional-commit` | `.agents/skills/conventional-commit/SKILL.md` |
| `apollo-client` | `.agents/skills/apollo-client/SKILL.md` |
| `find-skills` | `.agents/skills/find-skills/SKILL.md` |
| `nextjs-page-structure` | `.agents/skills/nextjs-page-structure/SKILL.md` |
| `shadcn` | `.agents/skills/shadcn/SKILL.md` |

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

## Convenciones de Código

- **Indentación**: Tabs (Prettier)
- **Comillas**: Simple quotes
- **Ancho**: 100 caracteres
- **Semicolones**: No trailing commas
- **TypeScript**: Strict mode activado
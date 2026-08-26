---
name: nextjs-page-structure
description: >
  Defines component organization patterns for Next.js App Router pages. Use when creating, refactoring, or organizing page.tsx files and their related components. Implements the "page-specific vs reusable" component separation rule.
user-invocable: false
allowed-tools: Read Write Edit Glob Grep Bash
---

# Next.js Page Component Structure

This skill defines the organizational pattern for components in Next.js App Router pages. It ensures separation between page-specific components and reusable components.

## When to Apply

Use this skill when:
- Creating a new `page.tsx` in `app/` directory
- Refactoring existing pages
- Adding new components to a page
- Deciding where to place a new component file

## Directory Structure Pattern

```
app/
└── cuotas/
    └── registrar/
        ├── page.tsx           # Page component (Server or Client)
        └── components/        # Page-specific components ONLY
            ├── registrar-cuota-form.tsx
            └── gastos-table.tsx

components/
├── ui/                       # shadcn/ui components ONLY
│   ├── button.tsx
│   └── dialog.tsx
├── reusable/                 # Shared components across pages
│   ├── data-table.tsx
│   └── filter-bar.tsx
└── overlays/                 # Project-specific overlays
    └── agregar-gasto-overlay.tsx
```

## Component Classification

### 1. Page-Specific Components (Co-located)

Located in `page/components/` directory. Used only by that specific page.

**Examples:**
- Form components unique to the page
- Data display components tied to page logic
- Custom table renderers for that page's data

**When to use:** Component is tightly coupled to the page's data, form logic, or user flow.

### 2. Reusable Components (Centralized)

Located in `components/` (excluding `ui/`). Used across multiple pages.

**Examples:**
- Generic data tables
- Filter/search components
- Navigation components
- Custom card layouts

**When to use:** Component could be useful in other pages, has generic props, no page-specific business logic.

### 3. UI Components (shadcn)

Located in `components/ui/`. All shadcn/ui components must stay here.

**Rule:** `@/components/ui` is **exclusive** to shadcn/ui components. Never place custom components here.

## Import Aliases

Use these import patterns:

```tsx
// ✅ Correct: Page component imports
import { RegistrarCuotaForm } from "./components/registrar-cuota-form";

// ✅ Correct: Reusable components
import { DataTable } from "@/components/reusable/data-table";
import { AgregarGastoOverlay } from "@/components/overlays/agregar-gasto-overlay";

// ✅ Correct: shadcn/ui components
import { Button } from "@/components/ui/button";
import { Dialog } from "@/components/ui/dialog";

// ❌ Wrong: Custom components in ui folder
import { MyCustomComponent } from "@/components/ui/my-custom-component"; // BAD

// ❌ Wrong: Page component imported from /components root
import { RegistrarCuotaForm } from "@/components/registrar-cuota-form"; // BAD if it's page-specific
```

## Page + Component Separation Pattern

When a page needs client-side interactivity (hooks, event handlers), separate into:

1. **Server Component (page.tsx)** - Data fetching, SSR
2. **Client Component (page/components/*.tsx)** - Interactive UI

```tsx
// app/cuotas/registrar/page.tsx (Server Component)
import { graphql } from "@/graphql";
import { execute } from "@/graphql/execute";
import { RegistrarCuotaForm } from "./components/registrar-cuota-form";

const ProveedoresQuery = graphql(`
  query Proveedores {
    obtenerProveedores { id nombre }
  }
`);

export default async function RegistrarCuota() {
  const { obtenerProveedores: proveedores } = await execute(ProveedoresQuery);
  return <RegistrarCuotaForm proveedores={proveedores} />;
}
```

```tsx
// app/cuotas/registrar/components/registrar-cuota-form.tsx (Client Component)
"use client";

import { useOverlay } from "@/hooks/useOverlay";
// ... page-specific component code
```

## Decision Flowchart

When creating a new component:

```
Is it a shadcn/ui component?
├── YES → Use shadcn CLI: npx shadcn@latest add <component>
└── NO
    │
    ├── Will other pages use it?
    │   ├── YES → @/components/reusable/<component-name>.tsx
    │   └── NO
    │       │
    │       └── Is it part of a page.tsx?
    │           ├── YES → page-directory/components/<component-name>.tsx
    │           └── NO → Consider if it should exist separately
```

## Examples

### Adding a new page component

User wants to create "Lista de proveedores" page:

1. Create `app/proveedores/lista/page.tsx`
2. If it needs client components, create `app/proveedores/lista/components/`
3. Place only page-specific components there
4. Any reusable components go to `@/components/`

### Moving an existing component

Existing file: `@/components/registrar-cuota-form.tsx`

This component is specific to `/cuotas/registrar` page, so move to:
```
app/cuotas/registrar/components/registrar-cuota-form.tsx
```

### Creating a shared table component

If building a table that will be used in multiple pages:
```
components/reusable/data-table.tsx
```

Import as: `import { DataTable } from "@/components/reusable/data-table"`

## Enforcement

When reviewing code, check:
- No custom components in `@/components/ui/`
- Page-specific components in `./components/` relative to page
- Reusable components in `@/components/` (not `ui/`)
- shadcn components properly in `@/components/ui/`
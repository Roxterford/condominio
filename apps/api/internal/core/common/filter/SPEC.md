# Especificación de Filtros Dinámicos

## Descripción General

Los filtros dinámicos son un lenguaje de consulta declarativo que permite construir consultas SQL complejas a partir de estructuras de datos simples (mapas/JSON). El sistema transforma un mapa de clave-valor en un Abstract Syntax Tree (AST) que luego se convierte a cláusulas WHERE/O WHERE.

## Estructura de Datos

Un filtro es un `map[string]any` donde las claves pueden ser:

- **Campos de entidad**: `status`, `age`, `id`, etc.
- **Operadores lógicos**: `and`, `or`, `not`

## Operadores de Comparación

| Clave | Descripción | Ejemplo |
|-------|-------------|---------|
| `eq` | Igual a | `{"status": {"eq": "active"}}` |
| `neq` | Diferente a | `{"age": {"neq": 0}}` |
| `gt` | Mayor que | `{"age": {"gt": 18}}` |
| `gte` | Mayor o igual que | `{"monto": {"gte": 1000}}` |
| `lt` | Menor que | `{"stock": {"lt": 10}}` |
| `lte` | Menor o igual que | `{"prioridad": {"lte": 5}}` |
| `like` | Coincidencia de patrón | `{"nombre": {"like": "%juan%"}}` |
| `in` | Está en lista | `{"status": {"in": ["active", "pending"]}}` |

### Equality Implícita

Cuando el valor es un valor primitivo (string, number, boolean), se asume equality:

```json
{"status": "active"}
```

Es equivalente a:

```json
{"status": {"eq": "active"}}
```

## Operadores Lógicos

### AND

Combina múltiples condiciones con lógica AND. Puede ser explícito (array) o implícito (múltiples claves).

**Explícito (array)**:
```json
{"and": [{"status": "active"}, {"age": {"gt": 18}}]}
```

**Implícito (múltiples claves)**:
```json
{"status": "active", "age": {"gt": 18}}
```

Genera: `status = 'active' AND age > 18`

### OR

Combina múltiples condiciones con lógica OR. **Importante: debe usar array**.

```json
{"or": [{"status": "pending"}, {"status": "failed"}]}
```

Genera: `status = 'pending' OR status = 'failed'`

> **Nota**: La sintaxis `{"or": {"id": "X"}}` (objeto en lugar de array) NO es válida. Siempre use array.

### NOT

Invierte una condición (aplica De Morgan). Puede usar objeto o array:

```json
{"not": {"status": "active"}}
```

O con array para múltiples condiciones:
```json
{"not": [{"status": "active"}, {"deleted": true}]}
```

Genera: `NOT (status = 'active')`

## Combinación de Operadores

### Regla Principal

Los operadores lógicos (`or`, `and`) deben usar **array** para combinar múltiples condiciones. El operador lógico explícito debe envolver el resto de condiciones.

**Caso válido con OR**:
```json
{"or": [{"id": {"eq": "abc"}}, {"id": {"eq": "xyz"}}]}
```

Genera: `id = 'abc' OR id = 'xyz'`

**Caso válido con AND explícito**:
```json
{"and": [{"status": "active"}, {"age": {"gt": 18}}]}
```

Genera: `status = 'active' AND age > 18`

### Múltiples Predicates (AND Implícito)

Cuando hay múltiples campos sin operador lógico explícito, se asume AND:

```json
{"status": "active", "age": {"gte": 18}}
```

Genera: `status = 'active' AND age >= 18`

### Condiciones Compuestas en un Campo

Un mismo campo puede tener múltiples condiciones (siempre se combinan con AND):

```json
{"age": {"gt": 18, "lt": 65}}
```

Genera: `age > 18 AND age < 65`

## Anidamiento

Los operadores lógicos pueden anidarse arbitrariamente:

```json
{
  "or": [
    {"status": "active", "age": {"gte": 18}},
    {"status": "pending", "verified": true}
  ]
}
```

Genera: `(status = 'active' AND age >= 18) OR (status = 'pending' AND verified = true)`

## Transformación a SQL

El parser transforma el mapa a un AST (Abstract Syntax Tree):

```
Input JSON           →    AST              →    SQL
{"id": "X"}         →  Predicate         →  id = 'X'
{"or": [...]}       →  LogicalClause(OR) →  (...) OR (...)
{"a": "x", "b": "y"} →  LogicalClause(AND) → a = 'x' AND b = 'y'
```

## Casos de Borde

### Empty Filter

Un mapa vacío retorna nil (sin WHERE clause).

### Single Predicate

Un solo predicate retorna directamente el nodo (sin wrap).

### Empty Array en Operador Lógico

Un operador con array vacío genera un nodo válido sin hijos (equivale a tautología o contradicción según el operador).

### Sintaxis Inválida

**OR con objeto (NO válido)**:
```json
{"or": {"id": "X"}}
```
Este formato no es válido. Use siempre array:
```json
{"or": [{"id": "X"}]}
```

### Valor Nulo

Valores `null` son ignorados (el campo no se incluye en el filtro).

## Ejemplos Completos

### Ejemplo 1: Filtro simple con OR

**Input**:
```json
{
  "or": [
    {"id": {"eq": "cmnjvw4xo00011tcmumziguhz"}},
    {"id": {"eq": "cmnki3r880005fucmudtiht09"}}
  ]
}
```

**SQL**: `id = 'cmnjvw4xo00011tcmumziguhz' OR id = 'cmnki3r880005fucmudtiht09'`

### Ejemplo 2: Filtro con predicate + OR

**Input**:
```json
{
  "status": "active",
  "or": [
    {"monto": {"gte": 1000}},
    {"prioridad": {"eq": "high"}}
  ]
}
```

**SQL**: `(status = 'active') AND (monto >= 1000 OR prioridad = 'high')`

> **Nota**: Cuando hay predicates + operador lógico, el parser genera AND implícito. Use solo `or` (sin campos adicionales) si necesita solo OR.

### Ejemplo 3: Filtro con predicate + NOT

**Input**:
```json
{
  "status": "active",
  "not": {"deleted": true}
}
```

**SQL**: `(status = 'active') AND NOT (deleted = true)`

## Guía de Referencia Rápida

### Casos Comunes

| Caso de Uso | Input JSON | SQL Resultado |
|------------|------------|---------------|
| Buscar por ID específico | `{"id": "abc123"}` | `id = 'abc123'` |
| Buscar por múltiples IDs | `{"or": [{"id": {"eq": "abc"}}, {"id": {"eq": "xyz"}}]}` | `id = 'abc' OR id = 'xyz'` |
| Filtrar por estado activo | `{"status": "active"}` | `status = 'active'` |
| Filtrar por estados múltiples | `{"or": [{"status": "active"}, {"status": "pending"}]}` | `status = 'active' OR status = 'pending'` |
| Rango de valores | `{"monto": {"gte": 1000, "lt": 5000}}` | `monto >= 1000 AND monto < 5000` |
| Excluir registros | `{"not": {"deleted": true}}` | `NOT (deleted = true)` |
| Filtrar por lista | `{"status": {"in": ["active", "pending", "failed"]}}` | `status IN ('active', 'pending', 'failed')` |
| Búsqueda por texto | `{"nombre": {"like": "%juan%"}}` | `nombre LIKE '%juan%'` |
| Múltiples filtros AND | `{"status": "active", "age": {"gte": 18}}` | `status = 'active' AND age >= 18` |

### Casos Complejos

| Caso de Uso | Input JSON | SQL Resultado |
|------------|------------|---------------|
| OR anidado con campos | `{"or": [{"status": "active", "monto": {"gte": 1000}}, {"status": "pending"}]}` | `(status = 'active' AND monto >= 1000) OR status = 'pending'` |
| Predicate + OR | `{"status": "active", "or": [{"verified": true}, {"priority": "high"}]}` | `(status = 'active') AND (verified = true OR priority = 'high')` |
| Triple OR | `{"or": [{"status": "a"}, {"status": "b"}, {"status": "c"}]}` | `status = 'a' OR status = 'b' OR status = 'c'` |
| NOT con objeto | `{"not": {"status": "deleted"}}` | `NOT (status = 'deleted')` |
| NOT con array | `{"not": [{"status": "deleted"}, {"archived": true}]}` | `NOT (status = 'deleted' AND archived = true)` |

### Errores Comunes

| Error | Input Incorrecto | Input Correcto |
|-------|------------------|----------------|
| OR sin array | `{"or": {"id": "X"}}` | `{"or": [{"id": {"eq": "X"}}]}` |
| AND sin array | `{"and": {"status": "a"}}` | `{"and": [{"status": "a"}]}` |
| Olvidar array | `{"status": {"eq": "active"}, "or": {"id": {"eq": "X"}}}` | `{"or": [{"status": {"eq": "active"}}, {"id": {"eq": "X"}}]}` |

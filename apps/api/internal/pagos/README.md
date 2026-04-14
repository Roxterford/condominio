# Sistema de Eventos de Pagos

## Arquitectura

El sistema utiliza una arquitectura event-driven para procesar pagos automáticamente:

### Componentes

1. **Event Bus (Redis)** - Publica eventos cuando se registran pagos
2. **Event Worker** - Procesa eventos pendientes de Redis
3. **Dispatcher** - Conecta eventos con sus handlers
4. **Event Handlers** - Ejecutan la lógica de negocio

### Flujo Completo

```
1. Registro de Pago
   ↓
2. Generación de evento "pago.registrado"
   ↓  
3. Publicación a Redis Streams
   ↓
4. Worker procesa evento (via HTTP endpoint)
   ↓
5. Dispatcher ejecuta handler
   ↓
6. AplicarPago use case aplica pago a deuda
```

## Endpoints HTTP

### Procesar Eventos
- **URL**: `POST /api/worker/process`
- **Descripción**: Procesa todos los eventos pendientes en Redis
- **Uso**: Ideal para cron jobs externos

### Health Check
- **URL**: `GET /api/worker/health`
- **Descripción**: Verifica el estado del worker y conexión a Redis

## Configuración

### Handlers de Eventos

Los handlers se configuran en `internal/pagos/config/event_handlers.go`:

```go
eventHandlersConfig := pagoConfig.NewEventHandlersConfig(aplicarPago)
```

### Eventos Disponibles

- **pago.registrado**: Se dispara cuando se crea un nuevo pago
  - Handler: Aplica automáticamente el pago a la deuda más antigua

## Uso con Cron Job

Para procesar eventos automáticamente, configurar un cron job:

```bash
# Cada 5 minutos
*/5 * * * * curl -X POST http://localhost:8081/api/worker/process
```

## Configuración de Redis

- **Stream**: `pagos`
- **Consumer Group**: `pagos-group`
- **Consumer**: `cron-worker`

## Notas Importantes

1. Los eventos se publican incluso si falla el worker (comentario en código)
2. El worker está diseñado para ejecución bajo demanda (no corre continuamente)
3. El sistema es tolerante a fallos: los eventos permanecen en Redis hasta ser procesados

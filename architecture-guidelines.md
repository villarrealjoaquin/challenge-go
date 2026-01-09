# Arquitectura y Lineamientos del Proyecto

Este documento define los principios arquitectónicos y de calidad
que deben respetarse en este proyecto Go.

El objetivo es priorizar claridad, simplicidad, testabilidad
y buenas prácticas idiomáticas de Go.

---

## 🏗️ Estructura de Capas

La aplicación se organiza en capas claras y simples:

### Responsabilidades

#### Handlers

- Manejan HTTP (request / response)
- Obtienen `context.Context` desde `http.Request`
- No contienen lógica de negocio
- Delegan todo el trabajo al service

#### Services

- Contienen la lógica de negocio
- Orquestan providers
- No conocen detalles de HTTP
- Son fácilmente testeables usando mocks

#### Providers

- Encapsulan dependencias externas (HTTP, APIs, etc.)
- Implementan interfaces
- Manejan errores de I/O
- Respetan cancelación vía `context.Context`

---

## 🔄 Uso de `context.Context`

### Principios

- El `context` nace en el handler (`r.Context()`)
- Se propaga hacia abajo (handler → service → provider)
- Nunca crear `context.Background()` dentro del flujo
- Nunca usar `context` para pasar datos de negocio

### Cuándo usar `context.Context`

- Operaciones de I/O
- Llamadas HTTP
- Funciones que pueden bloquear o tardar

### Excepciones: Uso de `context.Background()` en tests

**En tests es aceptable usar `context.Background()`** porque:

- Los tests no forman parte del flujo de la aplicación
- Los tests necesitan crear context para probar el comportamiento
- Están fuera del flujo handler → service → provider

**Ejemplo aceptable en tests:**

```go
// ✅ Correcto en tests
func TestService_GetMetrics(t *testing.T) {
    ctx := context.Background() // Aceptable en tests
    service := NewMetricsService(mockRepo)
    response, err := service.GetMetrics(ctx, "Author")
    // ... assertions
}
```

**Importante:** Esta excepción aplica **únicamente** en tests. En código de producción, siempre usar el context que se recibe como parámetro.

---

## 🔌 Dependencias e Inyección

### Principios

- Todas las dependencias deben inyectarse por constructor
- Evitar dependencias hardcodeadas
- Usar interfaces para permitir mocking en tests

### Ejemplo de inyección correcta

```go
// ✅ Correcto: Inyección por constructor
func NewGetMetrics(metricsService MetricsService) GetMetrics {
    return GetMetrics{metricsService}
}

// ❌ Incorrecto: Dependencia hardcodeada
func NewGetMetrics() GetMetrics {
    service := NewMetricsService(...) // Hardcodeado
    return GetMetrics{service}
}
```

---

## 📋 Flujo de Datos

```
HTTP Request
    ↓
Handler (parsea request, extrae context)
    ↓
Service (lógica de negocio)
    ↓
Repository (acceso a datos)
    ↓
Provider (llamadas externas HTTP/APIs)
    ↓
Response (vuelve por las capas)
```

---

## 🧪 Testing

### Principios

- En tests, usar mocks para aislar services y handlers
- Cada capa debe ser testeable de forma independiente
- Los mocks deben implementar las mismas interfaces que las implementaciones reales
- En tests es aceptable usar `context.Background()` o `context.WithTimeout()` para crear context de prueba

### Ejemplo de test con mocks

```go
type mockMetricsService struct {
    response *MetricsResponse
    err      error
}

func (m *mockMetricsService) GetMetrics(ctx context.Context, author string) (*MetricsResponse, error) {
    return m.response, m.err
}

func TestGetMetrics_OK(t *testing.T) {
    mockService := &mockMetricsService{
        response: &MetricsResponse{...},
    }
    handler := NewGetMetrics(mockService)
    // ... test implementation
}
```

---

## ✅ Checklist de Cumplimiento

Antes de escribir código, verificar:

- [ ] ¿Está en la capa correcta?
- [ ] ¿Usa `context.Context` correctamente?
- [ ] ¿Las dependencias se inyectan por constructor?
- [ ] ¿No hay dependencias hardcodeadas?
- [ ] ¿Se pueden crear mocks fácilmente?
- [ ] ¿Respeta la separación de responsabilidades?

---

## ⚠️ Violaciones Comunes

### ❌ Crear `context.Background()` en el flujo

```go
// ❌ Incorrecto
func (s *service) DoSomething() error {
    ctx := context.Background() // ¡NO!
    return s.repo.GetAll(ctx)
}

// ✅ Correcto
func (s *service) DoSomething(ctx context.Context) error {
    return s.repo.GetAll(ctx) // Usa el context recibido
}
```

### ❌ Lógica de negocio en handlers

```go
// ❌ Incorrecto
func (h *Handler) Handle() gin.HandlerFunc {
    return func(ctx *gin.Context) {
        books := h.repo.GetAll(ctx.Request.Context())
        // Lógica de negocio aquí (calcular promedio, etc.)
        avg := calculateAverage(books) // ¡NO! Esto va en el service
    }
}

// ✅ Correcto
func (h *Handler) Handle() gin.HandlerFunc {
    return func(ctx *gin.Context) {
        metrics, err := h.service.GetMetrics(ctx.Request.Context(), author)
        // Solo formatea la respuesta HTTP
    }
}
```

### ❌ Usar `context` para pasar datos de negocio

```go
// ❌ Incorrecto
ctx := context.WithValue(ctx, "userID", userID)

// ✅ Correcto
func (s *service) GetUserData(ctx context.Context, userID string) {
    // userID como parámetro explícito
}
```

---

## 📚 Convenciones Adicionales

### Nombres de archivos

- `handlers.go`, `services.go`, `providers.go` para implementaciones
- `*_test.go` para tests
- Interfaces y tipos en el mismo paquete que sus implementaciones

### Estructura de directorios

```
.
├── handlers/
│   ├── handlers.go
│   └── handlers_test.go
├── services/
│   ├── metrics.go
│   └── metrics_test.go
├── repositories/
│   ├── books.go
│   ├── books_test.go
│   └── mockImpls/
├── providers/
│   ├── books.go
│   └── books_test.go
└── models/
    └── books.go
```

---

## 🎯 Objetivos de Calidad

- **Claridad**: El código debe ser fácil de entender
- **Simplicidad**: Evitar sobre-ingeniería
- **Testabilidad**: Cada componente debe ser fácilmente testeable
- **Mantenibilidad**: Fácil de modificar y extender

---

## 📝 Notas Importantes

- Si el código propuesto viola alguna guideline, debe explicarse y corregirse
- Preferir soluciones simples e idiomáticas de Go
- Indicar claramente en qué capa vive cada responsabilidad
- En tests, usar mocks para aislar services y handlers

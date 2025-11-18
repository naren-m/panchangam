# Project Architecture

This document describes the architecture, structure, and design patterns of the Panchangam project.

## System Overview

The Panchangam project is a full-stack application that provides Hindu astronomical calendar calculations. It consists of:

- **Backend**: Go-based microservices architecture with gRPC
- **Frontend**: React-based web application
- **Data Layer**: Redis cache, Swiss Ephemeris astronomical data
- **Observability**: OpenTelemetry instrumentation

## Architecture Diagram

```
┌─────────────────┐
│   Web Browser   │
└────────┬────────┘
         │ HTTP/REST
         ▼
┌─────────────────┐
│  React Frontend │ (TypeScript, Vite, TailwindCSS)
│   (UI Layer)    │
└────────┬────────┘
         │ HTTP/REST
         ▼
┌─────────────────┐
│  Gateway Server │ (gRPC-Gateway, CORS, REST→gRPC)
│   (API Layer)   │
└────────┬────────┘
         │ gRPC
         ▼
┌─────────────────────────────────────┐
│        gRPC Server                  │
│  ┌─────────────┐  ┌──────────────┐ │
│  │ Panchangam  │  │  Sky View    │ │
│  │  Service    │  │   Service    │ │
│  └──────┬──────┘  └──────┬───────┘ │
│         │                 │         │
│         └────────┬────────┘         │
│                  │                  │
│         ┌────────▼────────┐         │
│         │   Astronomy     │         │
│         │   Core Logic    │         │
│         └────────┬────────┘         │
│                  │                  │
└──────────────────┼──────────────────┘
                   │
         ┌─────────┴─────────┐
         ▼                   ▼
┌──────────────┐    ┌────────────────┐
│ Redis Cache  │    │Swiss Ephemeris │
│   (KV Store) │    │   (Data Files) │
└──────────────┘    └────────────────┘
         │
         ▼
┌──────────────────┐
│  OpenTelemetry   │
│   (Observability)│
└──────────────────┘
```

## Project Structure

```
panchangam/
├── main.go                     # Main entry point (currently minimal)
├── go.mod/go.sum              # Go dependencies
├── Makefile                   # Build and test automation
│
├── cmd/                       # Command-line applications
│   ├── gateway/              # REST API gateway (gRPC-Gateway)
│   ├── grpc-server/          # Main gRPC server
│   ├── panchangam-cli/       # CLI tool for Panchangam calculations
│   ├── sunrise-demo/         # Demo application for sunrise calculations
│   ├── sunrise-simple/       # Simple sunrise calculator
│   └── test-service/         # Test utilities
│
├── services/                 # Business logic services
│   ├── panchangam/          # Panchangam calculation service
│   └── skyview/             # Sky visualization service
│
├── astronomy/               # Core astronomical calculations
│   ├── ephemeris/          # Ephemeris data handling (Swiss Ephemeris)
│   ├── validation/         # Historical validation tests
│   ├── sunrise.go          # Sunrise/sunset calculations
│   ├── lunar.go            # Lunar calculations (Tithi)
│   ├── festivals.go        # Festival date calculations
│   └── *_test.go           # Unit tests
│
├── gateway/                # Gateway server implementation
│   └── handlers/          # HTTP handlers and middleware
│
├── server/                # gRPC server implementation
│   └── interceptors/     # gRPC interceptors (auth, logging)
│
├── proto/                # Protocol Buffer definitions
│   └── *.proto          # gRPC service definitions
│
├── api/                  # API related code
│   ├── examples/        # API usage examples
│   └── implementations/ # API implementations
│
├── client/              # Client libraries
│
├── observability/       # Observability utilities
│   ├── tracing.go      # OpenTelemetry tracing
│   └── metrics.go      # Metrics collection
│
├── cache/              # Caching layer (Redis)
│
├── ui/                 # React frontend
│   ├── src/
│   │   ├── components/     # React components
│   │   ├── hooks/         # Custom React hooks
│   │   ├── services/      # API client services
│   │   ├── types/         # TypeScript type definitions
│   │   └── utils/         # Utility functions
│   ├── public/            # Static assets
│   ├── package.json       # Node dependencies
│   └── vite.config.ts    # Vite configuration
│
├── docs/               # Documentation
│   ├── algorithms/    # Algorithm documentation
│   ├── api/          # API documentation
│   ├── design/       # Design documents
│   ├── regional/     # Regional variations
│   └── validation/   # Validation documentation
│
├── llm/               # LLM/AI agent context documentation
│   ├── README.md
│   ├── coding-standards.md
│   ├── testing-guidelines.md
│   ├── project-architecture.md (this file)
│   └── domain-context.md
│
├── scripts/          # Utility scripts
├── docker/          # Docker configurations
└── test/           # Integration tests
```

## Backend Architecture

### Layer Separation

The backend follows a clean architecture pattern with clear layer separation:

```
┌────────────────────────────────────────┐
│         API Layer (gRPC/REST)          │  ← External interface
├────────────────────────────────────────┤
│        Service Layer (Business)        │  ← Business logic
├────────────────────────────────────────┤
│      Domain Layer (Astronomy Core)     │  ← Core calculations
├────────────────────────────────────────┤
│   Data Layer (Cache, Ephemeris)       │  ← Data access
└────────────────────────────────────────┘
```

### Core Components

#### 1. API Gateway (`cmd/gateway`)

**Purpose**: Provides HTTP/REST interface to gRPC services

**Responsibilities**:
- HTTP to gRPC translation (gRPC-Gateway)
- CORS handling
- Request/response transformation
- Rate limiting (future)

**Key Files**:
- `main.go`: Gateway server initialization
- `handler.go`: HTTP handlers

#### 2. gRPC Server (`cmd/grpc-server`)

**Purpose**: Main application server hosting gRPC services

**Responsibilities**:
- Service registration
- gRPC interceptors (logging, auth, metrics)
- Connection management
- Health checks

**Key Files**:
- `main.go`: Server initialization and startup
- `interceptors.go`: Custom gRPC interceptors

#### 3. Services Layer (`services/`)

**Purpose**: Business logic implementation

**Panchangam Service** (`services/panchangam/`):
- Calculate daily Panchangam
- Muhurta calculations
- Festival date calculations
- Ayanamsa conversions

**Sky View Service** (`services/skyview/`):
- 3D sky visualization data
- Planetary positions
- Constellation mapping

#### 4. Astronomy Core (`astronomy/`)

**Purpose**: Core astronomical calculations

**Key Modules**:
- **Ephemeris** (`astronomy/ephemeris/`): Swiss Ephemeris wrapper
- **Sunrise/Sunset**: Solar event calculations
- **Lunar**: Tithi, Nakshatra calculations
- **Yoga/Karana**: Combined calculations
- **Festivals**: Hindu festival date calculation

**Design Pattern**: Pure functions for calculations
```go
// Example: Pure calculation function
func CalculateTithi(sunLongitude, moonLongitude float64) (int, float64, error) {
    // No side effects, deterministic output
}
```

#### 5. Cache Layer (`cache/`)

**Purpose**: Performance optimization via Redis

**Cached Data**:
- Planetary positions (TTL: 1 hour)
- Daily Panchangam (TTL: 24 hours)
- Ephemeris lookups (TTL: varies)

**Cache Strategy**: Cache-aside pattern
```go
// 1. Check cache
if data, found := cache.Get(key); found {
    return data
}

// 2. Calculate/fetch
data := calculate()

// 3. Store in cache
cache.Set(key, data, ttl)

return data
```

#### 6. Observability (`observability/`)

**Purpose**: Monitoring and debugging

**Components**:
- **Tracing**: OpenTelemetry spans for request tracking
- **Metrics**: Service metrics (latency, errors, requests)
- **Logging**: Structured logging

## Frontend Architecture

### Component Structure

```
ui/src/
├── components/
│   ├── Panchangam/
│   │   ├── PanchangamDisplay.tsx      # Main container
│   │   ├── TithiCard.tsx              # Tithi display
│   │   ├── NakshatraCard.tsx          # Nakshatra display
│   │   └── YogaKaranaCard.tsx         # Yoga/Karana display
│   │
│   ├── SkyView/
│   │   ├── SkyVisualization.tsx       # 3D sky viewer
│   │   └── PlanetaryPositions.tsx     # Planet positions
│   │
│   ├── Common/
│   │   ├── DatePicker.tsx
│   │   ├── LocationPicker.tsx
│   │   └── LoadingSpinner.tsx
│   │
│   └── Layout/
│       ├── Header.tsx
│       ├── Footer.tsx
│       └── Navigation.tsx
│
├── hooks/
│   ├── usePanchangamData.ts          # Fetch Panchangam data
│   ├── useSkyViewData.ts             # Fetch sky view data
│   └── useLocation.ts                # Geolocation hook
│
├── services/
│   ├── api.ts                        # Base API client
│   ├── panchangamService.ts          # Panchangam API calls
│   └── skyviewService.ts             # Sky view API calls
│
├── types/
│   ├── panchangam.ts                 # Panchangam types
│   ├── astronomy.ts                  # Astronomy types
│   └── api.ts                        # API response types
│
└── utils/
    ├── formatters.ts                 # Date/time formatters
    ├── validators.ts                 # Input validation
    └── constants.ts                  # App constants
```

### State Management

**Current Approach**: React hooks and Context API

**State Categories**:
- **Local State**: Component-specific (useState)
- **Shared State**: App-wide settings (Context)
- **Server State**: API data (custom hooks with caching)

Example:
```typescript
// Custom hook for server state
function usePanchangamData(date: Date, location: Location) {
    const [data, setData] = useState<PanchangamData | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<Error | null>(null);

    useEffect(() => {
        // Fetch data
    }, [date, location]);

    return { data, loading, error };
}
```

### Routing

Currently single-page application. Future routes:
- `/` - Home/Daily Panchangam
- `/calendar` - Monthly calendar view
- `/sky` - Sky visualization
- `/muhurta` - Muhurta finder

## Data Flow

### Typical Request Flow

1. **User Action** (Frontend)
   ```typescript
   // User selects date
   const handleDateChange = (date: Date) => {
       setSelectedDate(date);
   };
   ```

2. **API Request** (Frontend)
   ```typescript
   // Hook fetches data
   const { data } = usePanchangamData(date, location);
   ```

3. **Gateway Receives Request** (Backend)
   ```
   GET /api/v1/panchangam?date=2024-01-01&lat=19.076&lon=72.877
   ```

4. **gRPC Translation** (Gateway)
   ```
   HTTP REST → gRPC CalculatePanchangam()
   ```

5. **Service Processing** (Service Layer)
   ```go
   func (s *PanchangamService) CalculatePanchangam(ctx context.Context, req *pb.PanchangamRequest) (*pb.PanchangamResponse, error) {
       // Check cache
       // Calculate if not cached
       // Return response
   }
   ```

6. **Core Calculation** (Astronomy Core)
   ```go
   // Get ephemeris data
   sunPos := ephemeris.GetSunPosition(date)
   moonPos := ephemeris.GetMoonPosition(date)

   // Calculate Panchangam elements
   tithi := CalculateTithi(sunPos, moonPos)
   nakshatra := CalculateNakshatra(moonPos)
   // ...
   ```

7. **Response Flow** (Backend → Frontend)
   ```
   gRPC Response → REST JSON → Frontend State Update
   ```

## Communication Protocols

### gRPC Service Definitions

```protobuf
service PanchangamService {
    rpc CalculatePanchangam(PanchangamRequest) returns (PanchangamResponse);
    rpc CalculateMuhurta(MuhurtaRequest) returns (MuhurtaResponse);
    rpc GetFestivals(FestivalRequest) returns (FestivalResponse);
}

message PanchangamRequest {
    string date = 1;        // ISO 8601 format
    double latitude = 2;
    double longitude = 3;
    string timezone = 4;
}

message PanchangamResponse {
    Tithi tithi = 1;
    Nakshatra nakshatra = 2;
    Yoga yoga = 3;
    Karana karana = 4;
    SolarEvents solar_events = 5;
}
```

### REST API Endpoints

```
GET  /api/v1/panchangam          # Get daily Panchangam
GET  /api/v1/muhurta             # Get muhurta timings
GET  /api/v1/festivals           # Get festival dates
GET  /api/v1/skyview             # Get sky view data
GET  /health                     # Health check
```

## Design Patterns

### 1. Repository Pattern (Future)

```go
type PanchangamRepository interface {
    GetDaily(ctx context.Context, date time.Time, loc Location) (*Panchangam, error)
    GetRange(ctx context.Context, start, end time.Time, loc Location) ([]*Panchangam, error)
}
```

### 2. Factory Pattern

```go
func NewPanchangamService(ephemeris EphemerisProvider, cache CacheProvider) *PanchangamService {
    return &PanchangamService{
        ephemeris: ephemeris,
        cache:     cache,
    }
}
```

### 3. Strategy Pattern

```go
type AyanamsaStrategy interface {
    Calculate(date time.Time) float64
}

type LahiriAyanamsa struct{}
type RamanAyanamsa struct{}
```

### 4. Dependency Injection

```go
type Service struct {
    ephemeris EphemerisProvider  // Interface, not concrete type
    cache     CacheProvider      // Interface, not concrete type
}
```

## Performance Considerations

### Caching Strategy

1. **Ephemeris Data**: Cache planetary positions (expensive to calculate)
2. **Daily Panchangam**: Cache complete day's data
3. **Sunrise/Sunset**: Cache per location and date

### Optimization Techniques

1. **Connection Pooling**: gRPC connection reuse
2. **Batch Calculations**: Calculate multiple days in single request
3. **Lazy Loading**: Load data on demand
4. **Code Splitting**: Split frontend bundles by route

## Security

### Authentication (Future)
- JWT tokens for API authentication
- API key validation

### Authorization
- Role-based access control (RBAC)
- Rate limiting per user/API key

### Data Validation
- Input validation at all layers
- Coordinate bounds checking
- Date range validation

## Deployment

### Docker Containers
- `gateway`: REST API gateway
- `grpc-server`: gRPC services
- `ui`: Frontend static files (Nginx)
- `redis`: Cache layer

### Environment Configuration
- Development: `.env.development`
- Production: `.env.production`

## Extension Points

### Adding New Calculations
1. Add function to `astronomy/` package
2. Add tests with 90% coverage
3. Expose via service layer
4. Add gRPC endpoint
5. Update gateway routing
6. Update frontend

### Adding New Services
1. Create service in `services/` folder
2. Define Protocol Buffer interface
3. Register with gRPC server
4. Add gateway routes
5. Implement business logic
6. Write comprehensive tests

## Technology Stack Summary

| Layer | Technology | Purpose |
|-------|-----------|---------|
| Frontend | React 18 + TypeScript | UI framework |
| UI Build | Vite | Build tool and dev server |
| Styling | TailwindCSS | CSS framework |
| 3D Graphics | Three.js | Sky visualization |
| Backend | Go 1.23 | Server-side logic |
| RPC | gRPC | Service communication |
| API Gateway | gRPC-Gateway | REST to gRPC |
| Cache | Redis | Performance optimization |
| Astronomy | Swiss Ephemeris | Planetary calculations |
| Observability | OpenTelemetry | Tracing and metrics |
| Testing (Backend) | Go testing, testify | Unit/integration tests |
| Testing (Frontend) | Vitest, Testing Library | Component tests |
| E2E Testing | Playwright | End-to-end tests |

## Maintenance and Evolution

### Adding Features
1. Define requirements
2. Design API changes
3. Update Protocol Buffers if needed
4. Implement backend logic with tests
5. Update frontend with tests
6. Update documentation
7. Create PR with issue reference

### Refactoring
1. Maintain backward compatibility
2. Use feature flags for gradual rollout
3. Update tests first
4. Refactor in small steps
5. Verify coverage remains ≥90%

### Breaking Changes
1. Version API endpoints
2. Deprecate old endpoints
3. Provide migration guide
4. Support old version for transition period

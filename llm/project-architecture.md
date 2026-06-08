# Project Architecture

This document describes the current architecture and source structure of the Panchangam project.

## System Overview

The Panchangam project is a full-stack application that provides Hindu astronomical calendar calculations. It consists of:

- **Backend**: Go gRPC server with a REST gateway
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
│  Gateway Server │ (REST handlers, CORS, validation)
│   (API Layer)   │
└────────┬────────┘
         │ gRPC
         ▼
┌─────────────────────────────────────┐
│        gRPC Server                  │
│        Panchangam Service           │
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
├── go.mod/go.sum              # Go dependencies
├── Makefile                   # Build and test automation
│
├── cmd/                       # Command-line applications
│   ├── gateway/              # REST gateway command
│   ├── server/               # Main gRPC server
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
├── gateway/                # REST gateway implementation
│   ├── server.go          # Server setup and route registration
│   ├── panchangam_handler.go
│   ├── current_tithi_handler.go
│   ├── sky_view_handler.go
│   ├── cache_handlers.go
│   ├── middleware.go
│   └── errors.go
│
├── proto/                # Protocol Buffer definitions
│   └── *.proto          # gRPC service definitions
│
├── api/                  # API related code
│   ├── examples/        # API usage examples
│   └── implementations/ # API implementations
│
├── observability/       # Observability utilities
│   ├── observability.go # OpenTelemetry setup and interceptors
│   └── errors.go        # Error recording and correlation helpers
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

The backend is split by job so each package has a clear owner:

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

#### 1. API Gateway (`cmd/gateway` and `gateway/`)

**Purpose**: Provides the HTTP/REST interface.

**Responsibilities**:
- Start the HTTP server
- Parse and validate REST requests
- Call `Panchangam.Get` through the gRPC client
- Serve sky view data from `services/skyview`
- Return structured JSON errors with request IDs

**Key Files**:
- `cmd/gateway/main.go`: Gateway process startup, flags, cache setup, and shutdown
- `gateway/server.go`: Route registration and HTTP server setup
- `gateway/panchangam_handler.go`: Daily Panchangam REST handler
- `gateway/current_tithi_handler.go`: Compact current tithi REST handler
- `gateway/sky_view_handler.go`: Sky visualization REST handler
- `gateway/middleware.go`: Request logging and `/api/v1/health`
- `gateway/errors.go`: Shared API error response shape

#### 2. gRPC Server (`cmd/server`)

**Purpose**: Hosts the Panchangam gRPC service.

**Responsibilities**:
- Start the gRPC server
- Register `Panchangam.Get`
- Add the observability unary interceptor
- Register the gRPC health service and reflection
- Handle graceful shutdown

**Key Files**:
- `cmd/server/main.go`: gRPC process startup, service registration, health checks, and shutdown
- `services/panchangam/server.go`: Panchangam server construction and dependencies
- `services/panchangam/service.go`: `Panchangam.Get` validation, tracing, and response building

#### 3. Services Layer (`services/`)

**Purpose**: Business logic used by the gateway and gRPC server.

**Panchangam Service** (`services/panchangam/`):
- Validates `GetPanchangamRequest`
- Records request spans and validation failures
- Fetches ephemeris-backed Panchangam data
- Builds the gRPC response

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
│   ├── Calendar/                    # Month navigation and calendar grid
│   ├── DayDetail/                   # Day detail modal and pancha anga views
│   ├── DataPresentation/            # Tabular data views
│   ├── TableView/                   # Table-focused calendar view
│   ├── GraphView/                   # Graph-focused calendar view
│   ├── ViewSwitcher/                # View mode controls
│   ├── CelestialChart/              # Celestial chart rendering
│   ├── EclipticBeltVisualization/   # Shared panchangam element calculations
│   ├── SkyVisualization/            # Sky sphere and time controls
│   ├── LocationPicker/              # Location selection
│   ├── Settings/                    # Settings and API health UI
│   └── common/                      # Shared loading and error components
│
├── hooks/
│   ├── usePanchangam.ts             # Panchangam data state
│   ├── useProgressivePanchangam.ts  # Progressive loading support
│   ├── useDayDetail.ts              # Day detail state
│   └── useOffline.ts                # Offline state
│
├── services/
│   ├── api/                         # Shared API client modules
│   ├── skyViewApi.ts                # Sky view API helpers
│   └── locationService.ts           # Location lookup helpers
│
├── types/
│   ├── panchangam.ts                # Panchangam types
│   └── skyVisualization.ts          # Sky visualization types
│
└── utils/
    ├── astronomy/                   # Coordinate, zodiac, nakshatra helpers
    ├── dateHelpers.ts               # Date helpers
    └── exportHelpers.ts             # Export helpers
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
function usePanchangam(date: Date, settings: Settings) {
    const [data, setData] = useState<PanchangamData | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<Error | null>(null);

    useEffect(() => {
        // Fetch data
    }, [date, settings]);

    return { data, loading, error };
}
```

### Routing

The frontend currently runs as a single-page application. Calendar, day detail, sky, graph, and settings views are handled by React component state.

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
   const { data } = usePanchangam(date, settings);
   ```

3. **Gateway Receives Request** (Backend)
   ```
   GET /api/v1/panchangam?date=2024-01-01&lat=19.076&lon=72.877
   ```

4. **gRPC Translation** (Gateway)
   ```
   HTTP REST → gRPC Panchangam.Get()
   ```

5. **Service Processing** (Service Layer)
   ```go
   func (s *PanchangamServer) Get(ctx context.Context, req *pb.GetPanchangamRequest) (*pb.GetPanchangamResponse, error) {
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
service Panchangam {
    rpc Get(GetPanchangamRequest) returns (GetPanchangamResponse);
}

message GetPanchangamRequest {
    string date = 1;          // YYYY-MM-DD
    double latitude = 2;
    double longitude = 3;
    string timezone = 4;
    string region = 5;
    string calculation_method = 6;
    string locale = 7;
}

message GetPanchangamResponse {
    PanchangamData panchangam_data = 1;
}
```

### REST API Endpoints

```
GET  /api/v1/panchangam          # Get daily Panchangam
GET  /api/v1/tithi/current       # Get compact current tithi summary
GET  /api/v1/sky-view            # Get sky visualization data
GET  /api/v1/health              # Health check
```

## Current Structure

Keep code paths direct:

1. `cmd/*` starts processes and handles flags, environment variables, health checks, and shutdown.
2. `gateway/*` owns HTTP parsing, validation, route registration, structured REST errors, and service calls.
3. `services/panchangam/*` owns the `Panchangam.Get` gRPC behavior.
4. `services/skyview/*` owns sky visualization calculations used by the REST gateway.
5. `astronomy/*` owns astronomical calculations.
6. `observability/*` owns tracing, metrics, and request spans.

Add a new interface or helper only when a real caller needs it. Do not add pattern-only wrappers.

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

## Input Safety

- Input validation at all layers
- Coordinate bounds checking
- Date range validation
- Structured JSON errors with request IDs
- gRPC validation errors recorded on spans

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
| API Gateway | Custom Go HTTP handlers | REST API and service calls |
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

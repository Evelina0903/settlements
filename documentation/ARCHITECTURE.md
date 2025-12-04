# Settlements Application Architecture

## Overview

The Settlements application is a Go-based web service that manages and visualizes data about Russian settlements (cities, villages, etc.). It provides APIs for querying settlement data and generates statistical visualizations based on settlement types and geographical distribution.

**Technology Stack:**
- Language: Go 1.24.0
- ORM: GORM (PostgreSQL)
- Database: PostgreSQL
- Custom HTTP Router with middleware support
- Template-based HTML rendering with embedded JSON data

---

## Core Architecture Layers

### 1. **Domain Models Layer** (`internal/models/`)

Represents the core business entities with database mappings using GORM tags.

#### Key Classes:

**`City`**
- Represents a settlement/city entity
- **Attributes:**
  - `ID` (uint): Primary key
  - `Name` (string): Settlement name
  - `Population` (int): Total population
  - `Childrens` (int): Number of children
  - `Latitude`, `Longitude` (float64): Geographical coordinates
  - `TypeID`, `DistrictID` (uint): Foreign keys
  - `Type`, `District` (relations): Associated domain objects
- **Purpose:** Maps to `cities` table in PostgreSQL

**`Type`**
- Represents settlement type classifications (e.g., "город", "деревня")
- **Attributes:**
  - `ID` (uint): Primary key
  - `Name` (string): Type name
  - `Citys` ([]City): One-to-many relationship
- **Purpose:** Maps to `types` table; provides settlement classification

**`District`**
- Represents Russian administrative regions (districts)
- **Attributes:**
  - `ID` (uint): Primary key
  - `Name` (string): District/region name
  - `Citys` ([]City): One-to-many relationship
- **Purpose:** Maps to `districts` table; enables geographical grouping

---

### 2. **Data Transfer Objects Layer** (`internal/dto/`)

Defines API response contracts and decouples domain models from presentation.

**`CityDTO`**
- Flat data transfer object for city information
- **Attributes:** ID, Name, Type (string), District (string), Population, Childrens, Latitude, Longitude
- **Purpose:** Used for API responses; resolves foreign keys into readable strings
- **Usage:** Repository transforms `City` model → `CityDTO` for external consumption

---

### 3. **Repository Layer** (`internal/repo/`)

Implements data access patterns and abstracts database queries.

**`CityRepo`**
- **Constructor:** `New(db *gorm.DB) *CityRepo`
- **Methods:**
  - `All() *[]CityDTO`: Fetches all cities with preloaded Type and District relations
  - `MinLongitude() float64`: Queries minimum longitude value for geographical range
  - `MaxLongitude() float64`: Queries maximum longitude value for geographical range
  - `GetCitiesInLongitudeGap(lMin, lMax float64) *[]CityDTO`: Range query for longitude-based data slicing
- **Dependencies:** GORM database instance
- **Pattern:** Encapsulates GORM queries and model-to-DTO conversion logic
- **Error Handling:** Logs fatal errors on query failures

---

### 4. **Service Layer** (`internal/service/`)

Implements business logic, data aggregation, and statistical computations.

**`Service`**
- **Constructor:** `New(cityRepo *repo.CityRepo) *Service`
- **Dependencies:** CityRepo instance
- **Key Methods:**

  **`GetAllSettelmetTypeData() *[]SettlementTypeData`**
  - Aggregates population statistics by settlement type
  - Computes: average population, average children, min/max population per type
  - Returns results sorted by average population (descending)
  - **Process:** Group by Type → Accumulate metrics → Calculate averages

  **`GetLongitudePopulationData() *[]GraphData`**
  - Distributes settlements into 100 longitude buckets
  - Calculates total population per bucket
  - Returns sorted by longitude for visualization
  - **Process:** Determine longitude range → Create 100 divisions → Sum population per bucket

  **`GetDistrictPopulationData() *[]GraphData`**
  - Aggregates total population by district
  - Returns results sorted by population (descending)
  - **Process:** Group by District → Sum population → Sort

**`SettlementTypeData`** (DTO-style struct)
- JSON-serializable output for settlement type statistics
- Fields: Type, AvgPopulation, AvgChildrens, MinPopulation, MaxPopulation

**`GraphData`** (Generic data structure)
- Flexible container for chart data
- Fields: X (any type), Y (int value)
- Used for both longitude-based and district-based visualizations

---

### 5. **Transport Layer** (`internal/transport/http/`)

Handles HTTP request routing and response generation.

#### Router Component (`router/`)

**`Router`** (Custom implementation)
- **Architecture:** Trie-based path router with segment-based matching
- **Key Structures:**
  - `node`: Tree node storing segment, children, parameter placeholders, handlers
  - `Params`: Map for route parameters (e.g., `:id` values)
- **Features:**
  - Static route matching (highest priority)
  - Dynamic parameter matching (`:param`)
  - Catch-all wildcard support (`*path`)
  - Middleware pipeline support (wrap-around pattern)
  - HTTP status handling (404, 405)
- **Methods:**
  - `New() *Router`: Creates router instance
  - `Handle(method, path string, h HandlerFunc)`: Registers route handler
  - `GET(path, handler)`, `POST(path, handler)`: HTTP method helpers
  - `Use(m MiddlewareFunc)`: Registers middleware (applied in reverse order during request)
  - `ServeHTTP(w, r)`: Makes router compatible with `net/http`
  - `match(method, path)`: Internal matching algorithm with parameter extraction
- **Middleware Pattern:** Functions wrap handlers in a chain for cross-cutting concerns

#### Controller Component (`controller/`)

**`MainController`**
- **Constructor:** `New(service *service.Service) *MainController`
- **Dependencies:** Service instance
- **Method:**
  - `GetMainPage(w, r, params)`: Renders dashboard with embedded JSON data
    - Fetches 3 datasets from service:
      1. Settlement type statistics (bar chart data)
      2. Longitude-based population distribution (line chart data)
      3. District-based population distribution (bar chart data)
    - Marshals data to JSON
    - Injects into HTML template as JavaScript variables
    - Returns rendered HTML with embedded visualization data

**Template Structure:**
- Embedded template with injected JavaScript variables
- Loads `web/templates/index.html` as base template
- Injects `tableData`, `chartData1`, `chartData2` into frontend context

---

### 6. **Configuration Layer** (`internal/config/`)

Manages environment-based configuration.

**`Config`**
- Aggregates server and database configuration
- **Attributes:** Server (ServerConfig), Database (DatabaseConfig)

**`ServerConfig`**
- **Attributes:** Port (string, default "3000")

**`DatabaseConfig`**
- **Attributes:** Host, Port (int), User, Password, Name
- **Method:** `DSN() string` - Generates PostgreSQL connection string
- **Pattern:** Builder-style DSN generation for GORM compatibility

**`Load()` Function**
- Reads environment variables via `godotenv`
- Provides sensible defaults (postgres:5432, user: postgres)
- Returns fully initialized Config struct or error
- **Fallback:** Graceful handling when `.env` file missing

---

### 7. **Database Layer** (`internal/db/`)

Handles database connection lifecycle.

**`Connect(cfg *DatabaseConfig) (*gorm.DB, error)`**
- Establishes PostgreSQL connection using GORM
- Takes configuration with DSN details
- Returns GORM database instance for all data access operations
- Error propagation for connection failures

---

### 8. **Data Loading Service** (`internal/service/data_loader/`)

Specialized service for bulk CSV data import.

**`DataLoader`**
- **Constructor:** `New(db *gorm.DB) *DataLoader`
- **Method:** `LoadCityData(filePath string) error`
  - Reads CSV file with settlement data
  - Processes each row via `processRow()`
  - Skips rows with insufficient columns
  - Continues on individual row errors
  
**`processRow(record []string)` Logic**
- Parses CSV columns: Region, Settlement, Type, Population, Children, Latitude, Longitude
- Maps abbreviated settlement types to full Russian names (e.g., "г" → "город")
- Handles coordinate transformation (converts negative longitudes)
- **Database Operations:**
  - `FirstOrCreate` Type record
  - `FirstOrCreate` District record
  - Aggregates population if settlement matches region (duplicate handling)
  - Creates or updates City record with relations
- **Deduplication:** Merges records when settlement name equals region name

**Settlement Type Mapping**
- 37+ type classifications (cities, villages, railway stations, etc.)
- Maps abbreviations to full Russian names
- Fallback to original abbreviation if not in mapping

---

### 9. **Utilities Layer** (`internal/util/`)

**`DateConverter`** (mentioned in glob results)
- Purpose: Date/time conversion utilities
- Usage context: Data transformation during loading

---

## Application Flow Diagrams

### Request Handling Flow
```
HTTP Request
    ↓
Router.ServeHTTP()
    ↓
Route Matching (Trie traversal)
    ↓
Middleware Pipeline (wrapped execution)
    ↓
MainController.GetMainPage()
    ↓
Service.Get*Data() methods (parallel)
    ↓
CityRepo.All() / Range queries
    ↓
PostgreSQL (GORM)
    ↓
JSON Marshal → Template Injection
    ↓
HTML Response
```

### Data Loading Flow
```
CSV File
    ↓
DataLoader.LoadCityData()
    ↓
Process Each Row
    ↓
Create/Find Type & District
    ↓
Create/Update City with aggregation
    ↓
PostgreSQL (GORM)
```

---

## Dependency Injection Pattern

The application uses **constructor-based dependency injection:**

```
main() 
  → config.Load() 
  → db.Connect(cfg) 
  → CityRepo.New(db)
  → Service.New(repo)
  → MainController.New(service)
  → Router registration
  → HTTP server startup
```

Each layer receives dependencies through constructors, enabling testability and loose coupling.

---

## Key Design Patterns

1. **Repository Pattern:** `CityRepo` abstracts database access
2. **Service Layer:** `Service` encapsulates business logic separate from HTTP concerns
3. **DTO Pattern:** `CityDTO` decouples API contracts from domain models
4. **Middleware Pipeline:** Router supports cross-cutting concerns without modifying handlers
5. **Trie Router:** Custom segment-based routing for efficient path matching
6. **Configuration Externalization:** Environment variables via `godotenv`
7. **Error Handling:** Log fatals in repository, graceful skipping in data loader

---

## Database Schema (Inferred)

```
types
├── id (PK)
└── name

districts
├── id (PK)
└── name

cities
├── id (PK)
├── name
├── type_id (FK → types)
├── district_id (FK → districts)
├── population
├── childrens
├── latitude
├── longitude
```

---

## Entry Points

1. **Web Server:** Main application serving HTTP requests for dashboard
2. **Data Loader:** CLI tool (`cmd/loader/main.go`) for importing CSV settlement data
   - Usage: `./loader -file datasets/dataset.csv`

---

## Security & Error Handling

- **Security:** Parameterized queries via GORM prevent SQL injection
- **Error Handling:** 
  - Repository: Fatal logging on database errors
  - Data Loader: Non-fatal row skipping with error reporting
  - HTTP: Standard 404/405 responses
- **Configuration:** Sensitive data from environment variables, not hardcoded

---

## Class Hierarchy & Relationships

### Overall Class Hierarchy Tree

```
Application
│
├─── Configuration Layer
│    └─── Config
│         ├─── ServerConfig
│         └─── DatabaseConfig
│
├─── Database Layer
│    └─── gorm.DB (external)
│
├─── Domain Models Layer
│    ├─── City
│    │    ├─── Type (relation)
│    │    └─── District (relation)
│    ├─── Type
│    │    └─── []City (one-to-many)
│    └─── District
│         └─── []City (one-to-many)
│
├─── DTO Layer
│    └─── CityDTO
│
├─── Repository Layer
│    └─── CityRepo
│         └─── gorm.DB (dependency)
│
├─── Service Layer
│    ├─── Service
│    │    └─── CityRepo (dependency)
│    ├─── SettlementTypeData (output DTO)
│    └─── GraphData (generic output)
│
├─── Data Loading Layer
│    └─── DataLoader
│         └─── gorm.DB (dependency)
│         └─── settlementsTypes map (configuration)
│
└─── HTTP Transport Layer
     ├─── Router
     │    ├─── node (internal tree)
     │    ├─── Params (map type alias)
     │    ├─── HandlerFunc (function type)
     │    └─── MiddlewareFunc (function type)
     │
     └─── MainController
          └─── Service (dependency)
               └─── template (embedded)

```

### Detailed Composition Relationships

#### Domain Model Relationships

```
Type ──────────┐
               │ 1:N
               ▼
             City  ◄──────────┐
          (belongs to both)   │
               ▲               │ 1:N
               │               │
            District ──────────┘
```

**Explanation:**
- `City` has a many-to-one relationship with `Type` via `TypeID` and `Type` field
- `City` has a many-to-one relationship with `District` via `DistrictID` and `District` field
- `Type` and `District` expose one-to-many relationships back to `City` via `Citys` slice

#### Dependency Injection Chain

```
                           ┌─────────────────────┐
                           │   main() entry      │
                           └──────────┬──────────┘
                                      │
                    ┌─────────────────┼─────────────────┐
                    │                 │                 │
                    ▼                 ▼                 ▼
            ┌────────────────┐ ┌────────────┐ ┌──────────────────┐
            │ config.Load()  │ │ db.Connect │ │ DataLoader.New() │
            │      ↓         │ │     ↓      │ │       ↓          │
            │   Config ◄─────┤ │   gorm.DB  │ │   DataLoader     │
            └────────────────┘ │     ▲      │ └──────────────────┘
                               └─────┼──────┘
                                     │
                               ┌─────┴────────────┐
                               │                  │
                               ▼                  ▼
                         ┌─────────────┐  ┌────────────────┐
                         │CityRepo.New │  │ Service.New    │
                         │     ↓       │  │     ↓          │
                         │  CityRepo◄──┤  │   Service      │
                         └─────────────┘  │     ▲          │
                                          │     │          │
                                          └─────┼──────────┘
                                                │
                               ┌────────────────┘
                               │
                               ▼
                    ┌──────────────────────────┐
                    │MainController.New(srv)   │
                    │       ↓                  │
                    │   MainController◄────────┤
                    │   + service: Service     │
                    └──────────────────────────┘
                               │
                               ▼
                    ┌──────────────────────────┐
                    │   Router.Handle(...)     │
                    │   Router.GET("/", ctrl)  │
                    │       ↓                  │
                    │    Router (registered)   │
                    └──────────────────────────┘
```

### Class Reference Tables

#### Domain Models

| Class | Fields | Relationships | GORM Table |
|-------|--------|---------------|-----------|
| **City** | ID, Name, Population, Childrens, Latitude, Longitude, TypeID, DistrictID, Type, District | Belongs-To Type, Belongs-To District | cities |
| **Type** | ID, Name, Citys | Has-Many City | types |
| **District** | ID, Name, Citys | Has-Many City | districts |

#### Repository Layer

| Class | Constructor | Key Methods | Dependencies |
|-------|-------------|------------|--------------|
| **CityRepo** | `New(db *gorm.DB)` | `All()`, `MinLongitude()`, `MaxLongitude()`, `GetCitiesInLongitudeGap()` | `*gorm.DB` |

#### Service Layer

| Class | Constructor | Key Methods | Dependencies | Output Types |
|-------|-------------|------------|--------------|--------------|
| **Service** | `New(repo *CityRepo)` | `GetAllSettelmetTypeData()`, `GetLongitudePopulationData()`, `GetDistrictPopulationData()` | `*CityRepo` | `[]SettlementTypeData`, `[]GraphData` |

#### HTTP Transport Layer

| Class | Constructor | Key Methods | Dependencies | Purpose |
|-------|-------------|------------|--------------|---------|
| **Router** | `New()` | `Handle()`, `GET()`, `POST()`, `Use()`, `ServeHTTP()`, `match()` | None | URL routing & middleware pipeline |
| **MainController** | `New(service)` | `GetMainPage(w, r, params)` | `*Service` | HTTP request handling |

#### Supporting Types

| Type | Kind | Fields/Values | Purpose |
|------|------|---------------|---------|
| **CityDTO** | struct | ID, Name, Type, District, Population, Childrens, Latitude, Longitude | API response contract |
| **SettlementTypeData** | struct | Type, AvgPopulation, AvgChildrens, MinPopulation, MaxPopulation | Chart data for settlement types |
| **GraphData** | struct | X (any), Y (int) | Generic chart data container |
| **Params** | map[string]string | Route parameters | Path variable capture |
| **HandlerFunc** | func type | `(w, r, params) → void` | HTTP handler signature |
| **MiddlewareFunc** | func type | `(h) → h` | Middleware wrapper signature |

#### Configuration Classes

| Class | Fields | Methods | Purpose |
|-------|--------|---------|---------|
| **Config** | Server, Database | - | Root configuration aggregator |
| **ServerConfig** | Port (string) | - | HTTP server settings |
| **DatabaseConfig** | Host, Port, User, Password, Name | `DSN() string` | Database connection settings |

#### Data Loading

| Class | Constructor | Key Methods | Dependencies |
|-------|-------------|------------|--------------|
| **DataLoader** | `New(db *gorm.DB)` | `LoadCityData(path)`, `processRow(record)` | `*gorm.DB` |

### Visibility & Scope

```
Application Public API:
├─ config.Load() ──────────────► Config
├─ db.Connect(cfg) ────────────► *gorm.DB
├─ repo.New(db) ───────────────► *CityRepo
├─ service.New(repo) ──────────► *Service
├─ controller.New(service) ────► *MainController
├─ router.New() ───────────────► *Router
├─ data_loader.New(db) ────────► *DataLoader
│
Internal Structures:
├─ router.node (unexported)
├─ service.SettlementTypeData (exported DTO)
├─ service.GraphData (exported DTO)
└─ data_loader.settlementsTypes (unexported map)
```

### Package Organization

```
settlements/
│
├── internal/
│   ├── config/          → Config, ServerConfig, DatabaseConfig
│   ├── db/              → Connect() function
│   ├── models/          → City, Type, District (domain entities)
│   ├── dto/             → CityDTO (data transfer object)
│   ├── repo/            → CityRepo (data access)
│   ├── service/         → Service, SettlementTypeData, GraphData
│   │   └── data_loader/ → DataLoader
│   ├── transport/http/  
│   │   ├── router/      → Router, node, Params, HandlerFunc, MiddlewareFunc
│   │   └── controller/  → MainController
│   └── util/            → DateConverter
│
└── cmd/
    └── loader/          → main() entry point for data loading
```

---

## Summary

The Settlements application demonstrates a **clean, layered architecture** with clear separation of concerns:

| Layer | Purpose | Key Classes |
|-------|---------|------------|
| Models | Domain representation | City, Type, District |
| DTO | API contracts | CityDTO |
| Repository | Data access | CityRepo |
| Service | Business logic | Service |
| HTTP Transport | Request handling | Router, MainController |
| Configuration | Settings management | Config |
| Database | Persistence | Connect |
| Data Loading | Bulk import | DataLoader |

This architecture supports maintainability, testability, and scalability while keeping the codebase focused on domain logic.

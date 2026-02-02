# Go Language Lifecycle

এই document এ Go programming language এর complete lifecycle explain করা হয়েছে - compilation থেকে execution, runtime, goroutines, garbage collection পর্যন্ত।

## 📋 Table of Contents

1. [Program Lifecycle](#program-lifecycle)
2. [Compilation Process](#compilation-process)
3. [Runtime Lifecycle](#runtime-lifecycle)
4. [Application Lifecycle](#application-lifecycle)
5. [Goroutine Lifecycle](#goroutine-lifecycle)
6. [Memory Lifecycle](#memory-lifecycle)
7. [Garbage Collection](#garbage-collection)
8. [Our E-Commerce Server Lifecycle](#our-e-commerce-server-lifecycle)

---

## 🚀 Program Lifecycle

### Complete Flow

```
Source Code (.go files)
    ↓
Go Compiler (go build)
    ↓
Object Files (.o)
    ↓
Linker
    ↓
Executable Binary
    ↓
OS Loads Binary
    ↓
Go Runtime Initializes
    ↓
main() Function Executes
    ↓
Program Runs
    ↓
Program Terminates
```

---

## 🔨 Compilation Process

### Step 1: Source Code
```go
// main.go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
```

### Step 2: Compilation (`go build`)
```bash
go build main.go
```

**What happens:**
1. **Lexical Analysis**: Source code → Tokens
2. **Parsing**: Tokens → Abstract Syntax Tree (AST)
3. **Type Checking**: Type validation
4. **Code Generation**: AST → Intermediate Representation (IR)
5. **Optimization**: Code optimization
6. **Assembly Generation**: IR → Assembly code
7. **Object File Creation**: Assembly → Object file (.o)
8. **Linking**: Object files → Executable binary

### Step 3: Binary Creation
```
main.go → main (executable)
```

**Binary contains:**
- Compiled code
- Go runtime
- Type information
- Reflection data

---

## ⚙️ Runtime Lifecycle

### Go Runtime Components

```
Go Runtime
├── Scheduler (Goroutine scheduler)
├── Garbage Collector (GC)
├── Memory Allocator
├── Network Poller
└── System Calls Handler
```

### Runtime Initialization Sequence

```
1. OS loads binary
   ↓
2. Runtime initialization
   ├── Memory allocator setup
   ├── Scheduler initialization
   ├── GC initialization
   └── Network poller setup
   ↓
3. Package initialization
   ├── Import packages
   ├── Run init() functions
   └── Initialize global variables
   ↓
4. main() function execution
   ↓
5. Program runs
   ↓
6. Cleanup and exit
```

---

## 📱 Application Lifecycle

### Our E-Commerce Server Example

```go
// cmd/server/main.go

func main() {
    // ============================================
    // PHASE 1: INITIALIZATION
    // ============================================
    
    // 1.1 Load configuration
    cfg := config.Load()
    // Reads .env file, sets up config struct
    
    // 1.2 Initialize utilities
    utils.InitJWT(cfg.JWTSecret)
    // Sets up JWT secret for token generation
    
    // 1.3 Connect to database
    if err := database.Connect(); err != nil {
        log.Fatalf("Failed to connect: %v", err)
    }
    // Establishes PostgreSQL connection (pgx + sqlc store)
    
    // 1.4 Setup router
    r := router.NewRouter()
    r.RegisterRoutes()
    // Registers all HTTP routes
    
    // 1.5 Apply middleware
    handler := middleware.CORS(middleware.Logging(r))
    // Chains middleware: CORS → Logging → Router
    
    // 1.6 Create HTTP server
    srv := &http.Server{
        Addr:         ":" + cfg.Port,
        Handler:      handler,
        ReadTimeout:  15 * time.Second,
        WriteTimeout: 15 * time.Second,
        IdleTimeout:  60 * time.Second,
    }
    
    // ============================================
    // PHASE 2: STARTUP
    // ============================================
    
    // 2.1 Start server in goroutine
    go func() {
        log.Printf("Server starting on port %s", cfg.Port)
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Server failed: %v", err)
        }
    }()
    // Server starts listening on port 8080
    
    // ============================================
    // PHASE 3: RUNNING
    // ============================================
    
    // 3.1 Server is now running
    // - Accepting HTTP connections
    // - Processing requests
    // - Handling multiple clients concurrently
    
    // 3.2 Wait for shutdown signal
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    // Blocks until SIGINT (Ctrl+C) or SIGTERM received
    
    // ============================================
    // PHASE 4: SHUTDOWN
    // ============================================
    
    // 4.1 Graceful shutdown
    log.Println("Shutting down server...")
    
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    // 4.2 Stop accepting new connections
    // 4.3 Wait for existing requests to complete (max 30 seconds)
    if err := srv.Shutdown(ctx); err != nil {
        log.Fatalf("Server forced to shutdown: %v", err)
    }
    
    // 4.4 Cleanup resources
    database.Disconnect()
    // Closes database connections
    
    // ============================================
    // PHASE 5: TERMINATION
    // ============================================
    
    log.Println("Server exited")
    // Program terminates
}
```

### Lifecycle Phases

#### Phase 1: Initialization
- Configuration loading
- Database connection
- Router setup
- Middleware chain creation

#### Phase 2: Startup
- Server starts listening
- Ready to accept connections
- Background goroutines start

#### Phase 3: Running
- Processing requests
- Handling concurrent connections
- Background tasks running

#### Phase 4: Shutdown
- Signal received (SIGINT/SIGTERM)
- Stop accepting new connections
- Wait for active requests to complete
- Cleanup resources

#### Phase 5: Termination
- All resources released
- Program exits

---

## 🔄 Goroutine Lifecycle

### What is a Goroutine?

Goroutine = Lightweight thread managed by Go runtime

### Goroutine States

```
New → Runnable → Running → Blocked → Runnable → Running → Dead
```

### State Transitions

```
1. New (Created)
   ↓ go func() { ... }()
   
2. Runnable (Ready to run)
   ↓ Scheduler picks it up
   
3. Running (Executing)
   ↓ Blocks (I/O, channel, sleep)
   
4. Blocked (Waiting)
   ↓ Event occurs
   
5. Runnable (Ready again)
   ↓ Scheduler picks it up
   
6. Running (Executing)
   ↓ Function completes
   
7. Dead (Terminated)
```

### Example: HTTP Request Handling

```go
// When a request arrives:

// 1. Main goroutine accepts connection
conn, err := listener.Accept()

// 2. New goroutine created for request
go handleRequest(conn)

// 3. Goroutine lifecycle:
func handleRequest(conn net.Conn) {
    // State: New → Runnable → Running
    
    // Read request (may block)
    data := make([]byte, 1024)
    n, err := conn.Read(data)  // Blocked (waiting for data)
    // State: Running → Blocked → Runnable → Running
    
    // Process request
    response := processRequest(data)  // Running
    
    // Write response (may block)
    conn.Write(response)  // Blocked → Runnable → Running
    
    // Close connection
    conn.Close()  // Running
    
    // Function ends
    // State: Running → Dead
}
```

### Goroutine Scheduling

```
Go Scheduler (M:N Model)
├── M OS Threads (Machine threads)
├── N Goroutines
└── Work Stealing Algorithm

Scheduler:
1. Runs on OS threads
2. Manages goroutines
3. Preemptive scheduling
4. Work-stealing for load balancing
```

### Our Server's Goroutines

```go
// Main goroutine
func main() {
    // Server setup...
    
    // New goroutine for server
    go func() {
        srv.ListenAndServe()  // Blocks, accepts connections
    }()
    
    // Main goroutine continues
    <-quit  // Blocks waiting for signal
    
    // Shutdown...
}
```

**Goroutines in our server:**
1. **Main goroutine**: Server setup, signal handling
2. **Server goroutine**: Accepts HTTP connections
3. **Request goroutines**: Each HTTP request handled in separate goroutine (automatic by net/http)

---

## 💾 Memory Lifecycle

### Memory Allocation

```
Stack Memory (Fast, Automatic)
├── Local variables
├── Function parameters
└── Return addresses

Heap Memory (Slower, GC managed)
├── Pointers
├── Slices (when capacity > length)
├── Maps
├── Channels
└── Interfaces
```

### Variable Lifecycle

```go
func example() {
    // 1. Stack allocation (automatic)
    x := 10  // Allocated on stack
    
    // 2. Heap allocation (via escape analysis)
    y := make([]int, 1000)  // May escape to heap if large
    
    // 3. Pointer (points to heap)
    ptr := new(int)  // Allocated on heap
    *ptr = 20
    
    // 4. Function returns
    // Stack frame destroyed
    // Heap objects remain (until GC)
}
```

### Escape Analysis

Go compiler decides: Stack or Heap?

```go
// Stays on stack
func stackExample() int {
    x := 10  // Stack
    return x
}

// Escapes to heap
func heapExample() *int {
    x := 10  // Heap (returned pointer)
    return &x
}
```

---

## 🗑️ Garbage Collection

### GC Lifecycle

```
Program Running
    ↓
Memory Allocated
    ↓
Objects Created
    ↓
Objects Become Unreachable
    ↓
GC Triggered (automatically)
    ↓
Mark Phase (find reachable objects)
    ↓
Sweep Phase (free unreachable objects)
    ↓
Memory Freed
    ↓
Program Continues
```

### GC Algorithm: Tri-Color Mark and Sweep

```
1. White: Unmarked (candidate for collection)
2. Gray: Marked, but children not checked
3. Black: Marked, all children checked (keep)
```

### GC Process

```
Initial State:
All objects = White

Step 1: Mark Roots (Gray)
- Global variables
- Stack variables
- Registers

Step 2: Mark Reachable (Gray → Black)
- Follow pointers from gray objects
- Mark referenced objects as gray
- Mark processed objects as black

Step 3: Sweep (White → Free)
- All white objects are unreachable
- Free their memory

Result:
- Black objects = Kept
- White objects = Freed
```

### GC Triggers

1. **Heap size threshold**: When heap grows to certain size
2. **Time-based**: Periodic GC (every 2 minutes by default)
3. **Manual**: `runtime.GC()` (not recommended)

### GC in Our Server

```go
// Database connection
Client = db.NewClient()  // Allocated on heap

// Request handling
func handler(w http.ResponseWriter, r *http.Request) {
    var req models.CreateUserRequest  // Stack
    utils.DecodeJSON(r, &req)  // Heap allocation for JSON
    
    user, err := database.Queries.CreateUser(...)  // Heap
    // User object allocated on heap
    
    response := models.ToUserResponse(user)  // Stack (struct copy)
    
    // After function returns:
    // - req: Stack frame destroyed
    // - user: Heap, will be GC'd when unreachable
    // - response: Stack frame destroyed
}
```

---

## 🔍 Our E-Commerce Server: Complete Lifecycle

### Startup Sequence

```
1. Binary Execution Starts
   ├── Go runtime initializes
   ├── Packages imported
   └── init() functions run
   
2. main() Function
   ├── config.Load()
   │   └── Reads .env, creates Config struct
   │
   ├── utils.InitJWT()
   │   └── Sets JWT secret
   │
   ├── database.Connect()
   │   ├── Opens PostgreSQL (pgx)
   │   ├── Creates Queries (sqlc Store)
   │   └── Tests connection
   │
   ├── router.NewRouter()
   │   └── Creates router, registers routes
   │
   ├── middleware.CORS(...)
   │   └── Wraps router with CORS
   │
   ├── middleware.Logging(...)
   │   └── Wraps with logging
   │
   └── http.Server{...}
       └── Creates server struct
   
3. Server Startup
   ├── go func() { srv.ListenAndServe() }()
   │   └── New goroutine starts
   │       └── Blocks on Accept()
   │
   └── signal.Notify(...)
       └── Main goroutine waits for signal
   
4. Server Running
   ├── Accepting connections
   ├── Each request → new goroutine
   └── Processing requests concurrently
```

### Request Processing Lifecycle

```
HTTP Request Arrives
    ↓
Server Accepts Connection
    ↓
New Goroutine Created (automatic by net/http)
    ↓
Goroutine State: New → Runnable → Running
    ↓
Request Processing:
    ├── CORS Middleware
    ├── Logging Middleware
    ├── Router Matching
    ├── Auth Middleware (if protected)
    ├── Handler Function
    │   ├── Decode JSON
    │   ├── Database Query
    │   │   └── May block (I/O)
    │   │       └── Goroutine: Running → Blocked → Runnable → Running
    │   ├── Process Data
    │   └── Send Response
    └── Logging Middleware (response)
    ↓
Response Sent
    ↓
Connection Closed
    ↓
Goroutine State: Running → Dead
    ↓
Goroutine Terminated
```

### Shutdown Sequence

```
SIGINT/SIGTERM Received
    ↓
Signal Channel Unblocks
    ↓
Shutdown Process Starts
    ↓
srv.Shutdown(ctx)
    ├── Stop accepting new connections
    ├── Wait for active requests (max 30s)
    └── Close idle connections
    ↓
database.Disconnect()
    └── Close database connections
    ↓
All Goroutines Terminate
    ↓
GC Runs (cleanup)
    ↓
Program Exits
```

---

## 📊 Memory Management in Our Server

### Memory Allocation Examples

```go
// 1. Configuration (heap, long-lived)
cfg := config.Load()  // Struct on heap

// 2. Database Client (heap, long-lived)
Queries = sqlc.NewStore(DB)  // Store on heap

// 3. Request Handler (stack + heap)
func handler(w http.ResponseWriter, r *http.Request) {
    // Stack allocations
    var req models.CreateUserRequest
    
    // Heap allocation (JSON decode)
    utils.DecodeJSON(r, &req)
    
    // Heap allocation (database query result)
    user, _ := database.Queries.GetUserByID(ctx, userID)
    
    // Stack allocation (response struct)
    response := models.ToUserResponse(user)
    
    // Heap allocation (JSON encoding)
    utils.RespondWithJSON(w, 200, response)
    
    // Function returns:
    // - Stack frame destroyed
    // - Heap objects remain until GC
}
```

### GC Impact

- **Frequent allocations**: Each request allocates memory
- **GC runs periodically**: Frees unreachable objects
- **Minimal pause**: Go's GC is concurrent (mostly)
- **Automatic**: No manual memory management needed

---

## 🎯 Key Takeaways

### 1. Compilation
- Go compiles to native binary
- Single binary deployment
- Fast compilation

### 2. Runtime
- Automatic memory management
- Goroutine scheduler
- Built-in GC

### 3. Concurrency
- Goroutines: Lightweight threads
- Channels: Communication
- Select: Multiplexing

### 4. Memory
- Stack: Fast, automatic
- Heap: GC managed
- Escape analysis: Compiler decides

### 5. Lifecycle
- Initialization → Running → Shutdown
- Graceful shutdown support
- Resource cleanup

### 6. Our Server
- Single binary
- Concurrent request handling
- Automatic memory management
- Graceful shutdown

---

## 🔗 Related Concepts

### Channels Lifecycle
```
Created → Send/Receive → Closed → Dead
```

### Context Lifecycle
```
Created → WithValue/WithTimeout → Cancel → Done
```

### HTTP Server Lifecycle
```
Created → Listen → Accept → Handle → Close
```

---

এই lifecycle understanding করে efficient, scalable Go applications build করা যায়! 🚀

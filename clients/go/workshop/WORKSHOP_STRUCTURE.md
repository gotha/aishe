# Workshop File Structure

This document describes the complete file structure of the AISHE Go Client Workshop.

## 📁 Directory Structure

```
clients/go/workshop/
├── README.md                          # Main workshop guide for students
├── INSTRUCTOR_GUIDE.md                # Detailed guide for instructors
├── WORKSHOP_STRUCTURE.md              # This file - structure overview
│
├── session1-basic/                    # Session 1: Basic HTTP Client
│   ├── starter/                       # Starting point for students
│   │   ├── go.mod                     # Go module file
│   │   ├── client.go                  # Client implementation (with TODOs)
│   │   └── main.go                    # Main program (with TODOs)
│   └── solution/                      # Reference solution
│       ├── go.mod                     # Go module file
│       ├── client.go                  # Complete client implementation
│       └── main.go                    # Complete main program
│
├── session2-redis/                    # Session 2: Redis Cache Client
│   ├── starter/                       # Starting point for students
│   │   ├── go.mod                     # Go module file (includes go-redis)
│   │   ├── client.go                  # Client implementation (with TODOs)
│   │   └── main.go                    # Main program (with TODOs)
│   └── solution/                      # Reference solution
│       ├── go.mod                     # Go module file (includes go-redis)
│       ├── client.go                  # Complete client implementation
│       └── main.go                    # Complete main program with benchmarking
│
└── session3-langcache/                # Session 3: LangCache Client
    ├── starter/                       # Starting point for students
    │   ├── go.mod                     # Go module file
    │   ├── client.go                  # Client implementation (with TODOs)
    │   └── main.go                    # Main program (with TODOs)
    └── solution/                      # Reference solution
        ├── go.mod                     # Go module file (includes go-redis for benchmark)
        ├── client.go                  # Complete client implementation
        ├── main.go                    # Complete main program
        └── benchmark.go               # Benchmark comparing all three approaches
```

## 📝 File Descriptions

### Root Level Files

#### README.md
- **Purpose:** Main student-facing documentation
- **Contents:**
  - Workshop overview and objectives
  - Prerequisites and setup instructions
  - Detailed guide for each session
  - Learning objectives and key concepts
  - Troubleshooting guide
  - Performance comparison table

#### INSTRUCTOR_GUIDE.md
- **Purpose:** Comprehensive guide for workshop instructors
- **Contents:**
  - Session breakdown with timing
  - Teaching tips and common pitfalls
  - Setup instructions
  - Assessment criteria
  - Troubleshooting guide
  - Post-workshop follow-up

#### WORKSHOP_STRUCTURE.md
- **Purpose:** Documentation of file structure (this file)
- **Contents:**
  - Directory tree
  - File descriptions
  - Design decisions

### Session 1: Basic HTTP Client

**Learning Focus:** HTTP client fundamentals, JSON handling, error management

#### starter/client.go
- Struct definitions for API types (Source, AnswerResponse, etc.)
- Client struct with baseURL and httpClient
- Function signatures with TODO comments
- Hints for implementation

#### starter/main.go
- Basic program structure
- TODO comments for student implementation
- Expected output examples

#### solution/client.go
- Complete implementation of all functions
- Proper error handling
- Configuration via constants
- Resource cleanup

#### solution/main.go
- Complete working example
- Health check demonstration
- Question answering with output formatting

### Session 2: Redis Cache Client

**Learning Focus:** Caching patterns, Redis integration, performance measurement

#### starter/client.go
- All Session 1 types plus Redis-specific types
- Client struct with Redis client
- ClientOptions for configuration
- Function signatures with detailed TODO comments
- Hints about cache-aside pattern

#### starter/main.go
- Program structure for demonstrating caching
- TODO comments for cache hit/miss demonstration
- Performance measurement placeholders

#### solution/client.go
- Complete implementation with Redis integration
- Cache key generation with SHA-256
- Cache-aside pattern implementation
- TTL configuration
- Cache clearing functionality

#### solution/main.go
- Complete demonstration of caching
- Performance comparison (first vs second request)
- Multiple question examples
- Cache clearing demonstration
- Performance summary output

### Session 3: LangCache Client

**Learning Focus:** Semantic caching, REST API integration, authentication, comparison

#### starter/client.go
- All basic types plus LangCache-specific types
- Client struct with LangCache configuration
- Function signatures with detailed TODO comments
- Hints about semantic similarity

#### starter/main.go
- Program structure for semantic caching demo
- TODO comments for similar question testing
- Semantic matching explanation

#### solution/client.go
- Complete implementation with LangCache REST API
- Bearer token authentication
- Semantic search implementation
- Cache storage with JSON serialization
- Flush functionality

#### solution/main.go
- Complete demonstration of semantic caching
- Multiple similar questions to show semantic matching
- Performance measurement
- Detailed output with explanations
- Comparison with Session 2 approach
- Can run in benchmark mode: `go run . benchmark`

#### solution/benchmark.go
- **Special file for comprehensive comparison**
- Implements all three client types (Basic, Redis, LangCache)
- Runs same questions through all three approaches
- Measures and compares performance
- Generates detailed comparison table
- Shows cache hit rates for each approach
- Demonstrates semantic matching advantage

## 🎯 Design Decisions

### Why Separate Starter and Solution Directories?

1. **Clear separation** - Students can't accidentally see solutions
2. **Easy comparison** - Students can compare their work with solutions
3. **Flexible teaching** - Instructors can choose to show solutions or not
4. **Self-paced learning** - Students can check solutions when stuck

### Why Package Main Instead of Separate Packages?

1. **Simplicity** - Easier for students to run with `go run .`
2. **No import complexity** - No need to understand Go module paths yet
3. **Self-contained** - Each session is independent
4. **Quick iteration** - Students can test immediately

### Why TODOs Instead of Empty Functions?

1. **Guided learning** - TODOs provide structure and hints
2. **Incremental progress** - Students can implement one function at a time
3. **Clear expectations** - Students know what to implement
4. **Reduced frustration** - Hints help when students are stuck

### Why Include go.mod Files?

1. **Dependency management** - Ensures correct versions
2. **Reproducibility** - Same environment for all students
3. **No setup friction** - Students don't need to run `go mod init`
4. **Best practices** - Shows proper Go project structure

### Why Benchmark in Session 3?

1. **Culmination** - Brings together all three approaches
2. **Comparison** - Shows trade-offs clearly
3. **Real-world relevance** - Benchmarking is important in production
4. **Motivation** - Shows why semantic caching matters

## 🔄 Workshop Flow

### Session 1 Flow
1. Student opens `session1-basic/starter/`
2. Reads TODOs in `client.go`
3. Implements `NewClient()`, `CheckHealth()`, `AskQuestion()`
4. Tests with `go run .`
5. Compares with `session1-basic/solution/` if needed

### Session 2 Flow
1. Student opens `session2-redis/starter/`
2. Builds on Session 1 knowledge
3. Implements Redis caching logic
4. Tests cache hit/miss scenarios
5. Measures performance improvement
6. Compares with solution

### Session 3 Flow
1. Student opens `session3-langcache/starter/`
2. Implements LangCache integration
3. Tests semantic matching
4. Compares with Session 2 approach
5. Runs benchmark: `cd solution && go run . benchmark`
6. Understands trade-offs between all three approaches

## 📊 File Sizes and Complexity

| File | Lines | Complexity | Purpose |
|------|-------|------------|---------|
| session1/starter/client.go | ~95 | Low | Learning HTTP basics |
| session1/solution/client.go | ~144 | Low | Reference implementation |
| session2/starter/client.go | ~150 | Medium | Learning caching |
| session2/solution/client.go | ~238 | Medium | Complete cache implementation |
| session3/starter/client.go | ~185 | Medium-High | Learning semantic caching |
| session3/solution/client.go | ~353 | High | Complete LangCache integration |
| session3/solution/benchmark.go | ~300 | High | Comprehensive comparison |

## 🎓 Progressive Complexity

The workshop is designed with progressive complexity:

1. **Session 1:** Foundation - HTTP, JSON, errors
2. **Session 2:** Intermediate - Caching, Redis, performance
3. **Session 3:** Advanced - Semantic caching, REST APIs, comparison

Each session builds on the previous one, reinforcing concepts while introducing new ones.

## 🔧 Customization Points

Instructors can customize:

1. **Questions** - Change the example questions in main.go files
2. **Timeouts** - Adjust client timeouts for different environments
3. **Cache TTL** - Change cache expiration times
4. **Similarity threshold** - Adjust LangCache threshold (default: 0.9)
5. **Benchmark questions** - Modify benchmark.go to test different scenarios

## 📚 Additional Notes

### Why No Tests?

- Workshop focuses on implementation, not testing
- Testing could be added as an advanced exercise
- Keeps workshop duration manageable

### Why Constants Instead of Environment Variables?

- Simplicity - constants are visible and easy to modify
- No setup required - works out-of-the-box for Sessions 1 & 2
- Clear errors - easy to see if configuration is missing
- Educational - students see configuration in the code
- Consistency - same pattern across all sessions
- Focus - workshop is about HTTP clients, not configuration management

### Future Enhancements

Potential additions for future versions:

1. **Session 4:** Advanced patterns (retry, circuit breaker, rate limiting)
2. **Testing module:** Unit and integration tests
3. **CLI module:** Building a command-line tool
4. **Metrics module:** Adding Prometheus metrics
5. **Tracing module:** Adding OpenTelemetry tracing

---

**This structure supports a comprehensive, hands-on learning experience! 🚀**


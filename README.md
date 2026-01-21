# 1. What is `net/http` in Go?

`net/http` is Go’s **standard, production-grade HTTP implementation**.

It provides:

* An **HTTP server**
* An **HTTP client**
* A **routing + handler model**
* Middleware-like composition (without calling it middleware)
* TLS, cookies, headers, streaming, HTTP/2, etc.

> In Go, **HTTP is not a framework**.
> It’s a **library + philosophy**: small interfaces, explicit wiring, no magic.

---

# 2. Core Philosophy (Very Important)

Before APIs, understand **why it looks “simple” but feels different from Express**.

### Go’s principles behind `net/http`

1. **Interfaces over classes**
2. **Functions are first-class**
3. **Composition over inheritance**
4. **Concurrency is built-in**
5. **Explicit is better than implicit**

That’s why:

* No controllers
* No middleware keyword
* No request lifecycle hooks
* No router by default

You **build everything explicitly**, but in a clean, testable way.

---

# 3. The Absolute Core: `http.Handler`

Everything in `net/http` revolves around **one interface**:

```go
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}
```

### This is the heart of the entire package.

If something can:

* Receive a request
* Write a response

…it is an **HTTP handler**.

---

## 3.1 `http.HandlerFunc`

Go lets **functions implement interfaces**.

```go
type HandlerFunc func(ResponseWriter, *Request)

func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) {
    f(w, r)
}
```

So this works:

```go
func hello(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Hello World"))
}
```

And Go treats it as a `Handler`.

> Express equivalent:
>
> ```js
> (req, res) => res.send("Hello")
> ```

But in Go, **this is an interface implementation**, not magic.

---

# 4. `http.ResponseWriter`

`ResponseWriter` is how we send responses.

It’s an **interface**, not a struct.

```go
type ResponseWriter interface {
    Header() Header
    Write([]byte) (int, error)
    WriteHeader(statusCode int)
}
```

### Important rules

1. **Headers must be set before writing**
2. `Write()` implicitly sends status `200 OK`
3. Once headers are sent, they are locked

Example:

```go
w.Header().Set("Content-Type", "application/json")
w.WriteHeader(http.StatusCreated)
w.Write([]byte(`{"ok": true}`))
```

If we call `Write()` first → status becomes `200`.

---

# 5. `http.Request`

`*http.Request` represents **the entire incoming request**.

It is immutable *in practice* (you shouldn’t mutate it casually).

Key fields:

```go
type Request struct {
    Method string
    URL    *url.URL
    Header Header
    Body   io.ReadCloser
    Context context.Context
}
```

### Important parts

#### Method

```go
r.Method // "GET", "POST", etc.
```

#### URL

```go
r.URL.Path
r.URL.Query().Get("id")
```

#### Headers

```go
r.Header.Get("Authorization")
```

#### Body (VERY IMPORTANT)

```go
bodyBytes, _ := io.ReadAll(r.Body)
```

* Body is a **stream**
* Can be read **only once**
* Large bodies are streamed (memory-safe)

---

# 6. Starting a Server

The simplest server:

```go
http.ListenAndServe(":8080", nil)
```

### What does `nil` mean?

It means:

> “Use the **DefaultServeMux**”

---

# 7. ServeMux (Go’s Router)

`http.ServeMux` is Go’s **default router**.

```go
mux := http.NewServeMux()
```

Routes are registered like this:

```go
mux.HandleFunc("/", homeHandler)
mux.Handle("/api", apiHandler)
```

### Path matching rules

* `/` → matches everything
* `/api/` → prefix match
* `/api/users` → exact or deeper

⚠️ No route parameters (`:id`) built-in.

That’s why people use:

* `chi`
* `gorilla/mux`
* `httprouter`

But **ServeMux is extremely fast and simple**.

---

# 8. DefaultServeMux (Global Router)

When we do:

```go
http.HandleFunc("/", handler)
```

We are registering routes on a **global router**.

Internally:

```go
var DefaultServeMux = NewServeMux()
```

This is okay for:

* Small apps
* Learning
* Examples

For real apps → **always create your own mux**.

---

# 9. `http.Server` (Advanced Control)

For production, we usually create a server explicitly:

```go
server := &http.Server{
    Addr:    ":8080",
    Handler: mux,
}
server.ListenAndServe()
```

This allows:

* Timeouts
* TLS config
* Graceful shutdown

Example:

```go
&http.Server{
    ReadTimeout:  5 * time.Second,
    WriteTimeout: 10 * time.Second,
    IdleTimeout:  120 * time.Second,
}
```

This is **mandatory for production**.

---

# 10. Middleware in Go (Important Concept)

Go does not have “middleware” as a keyword.

Instead:

> Middleware = **a function that wraps a handler and returns a handler**

### Middleware signature

```go
func middleware(next http.Handler) http.Handler
```

Example:

```go
func logging(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Println(r.Method, r.URL.Path)
        next.ServeHTTP(w, r)
    })
}
```

Usage:

```go
handler := logging(finalHandler)
```

### Multiple middleware = function composition

```go
handler := auth(logging(finalHandler))
```

This is **functional programming**, not framework magic.

---

# 11. Context (`r.Context()`)

Every request has a context:

```go
ctx := r.Context()
```

Used for:

* Cancellation
* Deadlines
* Request-scoped values

When client disconnects → context is canceled.

Example:

```go
select {
case <-ctx.Done():
    return
case result := <-dbQuery:
    // respond
}
```

This is **huge for scalable systems**.

---

# 12. HTTP Client (`http.Client`)

Go also provides an HTTP client.

Basic request:

```go
resp, err := http.Get("https://api.example.com")
```

Better way:

```go
client := &http.Client{
    Timeout: 10 * time.Second,
}

req, _ := http.NewRequest("GET", url, nil)
req.Header.Set("Authorization", "Bearer token")

resp, err := client.Do(req)
```

⚠️ Always:

```go
defer resp.Body.Close()
```

---

# 13. Cookies

### Reading

```go
cookie, err := r.Cookie("session")
```

### Writing

```go
http.SetCookie(w, &http.Cookie{
    Name:     "session",
    Value:    "abc123",
    HttpOnly: true,
})
```

---

# 14. File Uploads & Forms

### Parse form

```go
r.ParseForm()
r.FormValue("email")
```

### Multipart (file upload)

```go
r.ParseMultipartForm(10 << 20) // 10MB
file, header, err := r.FormFile("avatar")
```

---

# 15. Streaming Responses

Go excels at streaming.

```go
w.Write([]byte("chunk 1"))
w.(http.Flusher).Flush()
```

Used for:

* SSE
* Large files
* Real-time responses

---

# 16. Concurrency Model (CRITICAL)

Every incoming request is handled in its **own goroutine**.

> You do NOT create threads manually.

This means:

* Handlers must be **thread-safe**
* Shared state must be protected (`sync.Mutex`)
* Blocking = expensive

Compared to Node:

* Node → event loop
* Go → goroutine per request

---

# 17. Why Go HTTP Feels “Low-Level”

Because it is.

But the tradeoff:

* 🔥 Extremely fast
* 🔥 Explicit control
* 🔥 Easy to test
* 🔥 Zero framework lock-in

Frameworks like Gin, Fiber, Echo are **thin layers on top of `net/http`**.

---

# 18. How We Should Learn This

**Step 1**
Master:

* `Handler`
* `ServeMux`
* Middleware pattern
* Context

**Step 2**
Build:

* Auth middleware
* JSON API
* Graceful shutdown

**Step 3**
Then use a router like `chi`—not before.

---

# Go Batcher

A generic and thread-safe batcher for Go, designed to accumulate items over time or until a capacity is reached, then process them in bulk.

## ✨ Features

- ✅ Support for generics
- ✅ Thread-safe
- ✅ Auto-flush on timeout or capacity
- ✅ Supports batch add: `Add(x1, x2, ..., xn)`
- ✅ Graceful shutdown via context

## 🔧 Installation

```bash
go get github.com/LullNil/go-batcher
```

## 📦 Usage

```go
import "github.com/LullNil/go-batcher"

func main() {
	b := batcher.New(3, 500*time.Millisecond, func(batch []string) {
		fmt.Printf("Batch processed: %v\n", batch)
	})

	b.Add("alpha", "beta", "gamma")
	b.Add("delta")

	b.Close()
}
```

## 🛠 Use Cases

- Database inserts in bulk
- REST API batching
- Event aggregation
- Logging optimization
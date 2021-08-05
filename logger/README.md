# logger

Logger leverages [uber/zap](https://github.com/uber-go/zap) for best performance and availability.

## Example

Production JSON format in error level:

```go
logger := logger.New(logger.Level("error")).Sugar()
defer logger.Sync()
logger.Errorw("failed to fetch URL",
  // Structured context as loosely typed key-value pairs.
  "url", url,
  "attempt", 3,
  "backoff", time.Second,
)
// {"level":"error","ts":1589648882.6028602,"msg":"failed to fetch URL","url":"https://localhost/foo","attempt":3,"backoff":1}
```

Development console print format in debug level:

```go
logger := logger.New(logger.Level("debug"), logger.Encoder("console")).Sugar()
defer logger.Sync()
logger.Debugw("failed to fetch URL",
  // Structured context as loosely typed key-value pairs.
  "url", url,
  "attempt", 3,
  "backoff", time.Second,
)
// 2020-05-16T19:09:41.991+0200    Debug   failed to fetch URL     {"url": "https://localhost/foo", "attempt": 3, "backoff": "1s"}
```

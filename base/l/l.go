package l

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/k0kubun/pp"
	"github.com/uber-go/zap"
)

const prefix = "github.com/ng-vu/go-grpc-sample/"

var ll = New()

// Logger wraps zap.Logger
type Logger struct {
	zap.Logger
}

// Short-hand functions for logging.
var (
	Base64    = zap.Base64
	Bool      = zap.Bool
	Float64   = zap.Float64
	Int       = zap.Int
	Int64     = zap.Int64
	Marshaler = zap.Marshaler
	Nest      = zap.Nest
	Skip      = zap.Skip
	Stack     = zap.Stack
	String    = zap.String
	Stringer  = zap.Stringer
	Time      = zap.Time
	Uint      = zap.Uint
	Uint64    = zap.Uint64
	Uintptr   = zap.Uintptr
	Error     = zap.Error
)

// Interface ...
func Interface(key string, val interface{}) zap.Field {
	if val, ok := val.(fmt.Stringer); ok {
		return zap.Stringer(key, val)
	}
	return zap.Object(key, val)
}

// Duration formats time.Duration in human-readable time
func Duration(key string, val time.Duration) zap.Field {
	return zap.Stringer(key, val)
}

// Int32 ...
func Int32(key string, val int32) zap.Field {
	return zap.Int(key, int(val))
}

// Object ...
func Object(key string, val interface{}) zap.Field {
	return zap.Stringer(key, dump(val))
}

type dd struct {
	v interface{}
}

func (d dd) String() string {
	return pp.Sprint(d.v)
}

func dump(v interface{}) dd {
	return dd{v}
}

const development = true

// New returns new zap.Logger
func New() Logger {
	_, filename, _, _ := runtime.Caller(1)
	name := filepath.Dir(truncFilename(filename))

	var enabler zap.AtomicLevel
	if e, ok := enablers[name]; ok {
		enabler = e
	} else {
		enabler = zap.DynamicLevel()
		enablers[name] = enabler
	}

	setLogLevelFromEnv(name, enabler)
	return Logger{
		zap.New(
			zap.NewTextEncoder(zap.TextNoTime()),
			enabler, addHook(),
		),
	}
}

// ServeHTTP supports logging level with an HTTP request.
func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	type errorResponse struct {
		Error string `json:"error"`
	}
	type payload struct {
		Name  string     `json:"name"`
		Level *zap.Level `json:"level"`
	}

	enc := json.NewEncoder(w)

	switch r.Method {
	case "GET":
		var payloads []payload
		for k, e := range enablers {
			lvl := e.Level()
			payloads = append(payloads, payload{
				Name:  k,
				Level: &lvl,
			})
		}
		enc.Encode(payloads)

	case "PUT":
		var req payload
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			enc.Encode(errorResponse{
				Error: fmt.Sprintf("Request body must be valid JSON: %v", err),
			})
			return
		}

		enabler, ok := enablers[req.Name]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			enc.Encode(errorResponse{
				Error: errEnablerNotFound.Error(),
			})
			return
		}

		if req.Level == nil {
			w.WriteHeader(http.StatusBadRequest)
			enc.Encode(errorResponse{
				Error: errLevelNil.Error(),
			})
			return
		}

		enabler.SetLevel(*req.Level)
		enc.Encode(req)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		enc.Encode(errorResponse{
			Error: "Only GET and PUT are supported.",
		})
	}
}

var (
	errHookNilEntry    = errors.New("can't call a hook on a nil *Entry")
	errCaller          = errors.New("failed to get caller")
	errEnablerNotFound = errors.New("enabler not found")
	errLevelNil        = errors.New("must specify a logging level")

	enablers = make(map[string]zap.AtomicLevel)

	bufPool = sync.Pool{New: func() interface{} {
		buf := make([]byte, 0, 1024)
		return &buf
	}}
)

const (
	_callerSkip = 4
	resetColor  = "\x1b[0m"

	black   = "\x1b[30m"
	red     = "\x1b[31m"
	green   = "\x1b[32m"
	yellow  = "\x1b[33m"
	blue    = "\x1b[34m"
	magenta = "\x1b[35m"
	cyan    = "\x1b[36m"
	white   = "\x1b[37m"
	gray    = "\x1b[90m"
)

func addHook() zap.Option {
	return zap.Hook(func(e *zap.Entry) error {
		if e == nil {
			return errHookNilEntry
		}
		_, filename, line, ok := runtime.Caller(_callerSkip)
		if !ok {
			return errCaller
		}

		// Re-use a buffer from the pool.
		buf := bufPool.Get().(*[]byte)

		t := e.Time
		year, month, day := t.Date()
		itoa(buf, year, 4)
		*buf = append(*buf, '/')
		itoa(buf, int(month), 2)
		*buf = append(*buf, '/')
		itoa(buf, day, 2)
		*buf = append(*buf, ' ')

		hour, min, sec := t.Clock()
		itoa(buf, hour, 2)
		*buf = append(*buf, ':')
		itoa(buf, min, 2)
		*buf = append(*buf, ':')
		itoa(buf, sec, 2)

		*buf = append(*buf, '.')
		itoa(buf, t.Nanosecond()/1e6, 3)
		*buf = append(*buf, ' ')

		if development {
			switch e.Level {
			case zap.ErrorLevel, zap.PanicLevel, zap.DPanicLevel, zap.FatalLevel:
				*buf = append(*buf, red...)
			case zap.WarnLevel:
				*buf = append(*buf, yellow...)
			case zap.InfoLevel:
				*buf = append(*buf, cyan...)
			case zap.DebugLevel:
				*buf = append(*buf, blue...)
			}
		}

		*buf = append(*buf, e.Message...)
		if development {
			*buf = append(*buf, gray...)
		}
		*buf = append(*buf, " → "...)
		*buf = append(*buf, truncFilename(filename)...)
		*buf = append(*buf, ':')
		*buf = strconv.AppendInt(*buf, int64(line), 10)
		if development {
			*buf = append(*buf, resetColor...)
		}
		*buf = append(*buf, ' ')

		newMsg := string(*buf)
		*buf = (*buf)[:0]
		bufPool.Put(buf)
		e.Message = newMsg
		return nil
	})
}

func truncFilename(filename string) string {
	index := strings.Index(filename, prefix)
	return filename[index+len(prefix):]
}

// Cheap integer to fixed-width decimal ASCII.  Give a negative width to avoid zero-padding.
func itoa(buf *[]byte, i int, wid int) {
	// Assemble decimal in reverse order.
	var b [20]byte
	bp := len(b) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	// i < 10
	b[bp] = byte('0' + i)
	*buf = append(*buf, b[bp:]...)
}

var envPatterns []*regexp.Regexp

func init() {
	envLog := os.Getenv("LOG_DEBUG")
	if envLog == "" {
		return
	}

	var errPattern string
	envPatterns, errPattern = initPatterns(envLog)
	if errPattern != "" {
		ll.Fatal("Unable to parse LOG_DEBUG. Please set it to a proper value.", String("invalid", errPattern))
	}

	ll.Info("Enable debug log", String("LOG_DEBUG", envLog))
}

func initPatterns(envLog string) ([]*regexp.Regexp, string) {
	patterns := strings.Split(envLog, ",")
	result := make([]*regexp.Regexp, len(patterns))
	for i, p := range patterns {
		r, err := parsePattern(p)
		if err != nil {
			return nil, p
		}

		result[i] = r
	}
	return result, ""
}

func parsePattern(p string) (*regexp.Regexp, error) {
	p = strings.Replace(strings.Trim(p, " "), "*", ".*", -1)
	return regexp.Compile(p)
}

func setLogLevelFromEnv(name string, enabler zap.AtomicLevel) {
	for _, r := range envPatterns {
		if r.MatchString(name) {
			enabler.SetLevel(zap.DebugLevel)
		}
	}
}

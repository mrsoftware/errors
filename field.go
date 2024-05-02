package errors

import (
	"context"
	"fmt"
	"math"
	"time"
)

var (
	minTimeInt64 = time.Unix(0, math.MinInt64)
	maxTimeInt64 = time.Unix(0, math.MaxInt64)
)

// Field is a stack-based field and used for items that we want to store in error.
type Field struct {
	Key       string
	Type      FieldType
	Integer   int64
	Str       string
	Interface interface{}
}

// Format is implement the fmt.Formatter for field.
func (f Field) Format(state fmt.State, verb rune) {
	switch verb {
	case 'v':
		if state.Flag('+') {
			fmt.Fprintf(state, "{Key: %s, Type: %s, Value: %+v}", f.Key, f.Type, f.Value())

			return
		}

		if state.Flag('#') {
			fmt.Fprintf(state, "{%s: %#v}", f.Key, f.Value())

			return
		}

		fmt.Fprintf(state, "{Key: %s, Value: %+v}", f.Key, f.Value())
	case 's':
		fmt.Fprintf(state, "[%s: %s]", f.Key, f.Value())
	case 'q':
		fmt.Fprintf(state, "%q", f.Value())
	}
}

// Value of Field.
func (f Field) Value() interface{} {
	switch f.Type {
	case FieldTypeString:
		return f.Str
	case FieldTypeInt64:
		return f.Integer
	case FieldTypeFloat64:
		return math.Float32frombits(uint32(f.Integer))
	case FieldTypeBinary:
		return f.Interface
	case FieldTypeByteString:
		return f.Interface
	case FieldTypeError:
		return f.Interface
	case FieldTypeTimeFull:
		return f.Interface.(time.Time)
	case FieldTypeTime:
		return time.Unix(0, f.Integer).In(f.Interface.(*time.Location))
	case FieldTypeDuration:
		return time.Duration(f.Integer)
	case FieldTypeBool:
		var b bool
		if f.Integer == 1 {
			b = true
		}

		return b
	default:
		return f.Interface
	}
}

// String version of Field.
func (f Field) String() string {
	return fmt.Sprintf("Key: %s, Type: %s, Value: %s", f.Key, f.Type, f.Value())
}

// Is compare field type.
func (f Field) Is(fieldType FieldType) bool {
	return f.Type == fieldType
}

// FieldType is a Field data type.
type FieldType int

const (
	// FieldTypeUnknown is used if the data type is unknown.
	FieldTypeUnknown FieldType = iota

	// FieldTypeReflect is used for fields that store Reflect.
	FieldTypeReflect

	// FieldTypeString is used for fields that store String.
	FieldTypeString

	// FieldTypeInt64 is used for fields that store Int/Int64.
	FieldTypeInt64

	// FieldTypeFloat64 is used for fields that store Float/Float64.
	FieldTypeFloat64

	// FieldTypeBinary is used for fields that store binary data as []byte.
	FieldTypeBinary

	// FieldTypeBool is used for fields that store Bool.
	FieldTypeBool

	// FieldTypeByteString is used for fields that store string data as []byte.
	FieldTypeByteString

	// FieldTypeTime is used for fields that store UnixNano time.
	FieldTypeTime

	// FieldTypeTimeFull is used for fields that store time.Time.
	FieldTypeTimeFull

	// FieldTypeDuration is used for fields that store time.Duration.
	FieldTypeDuration

	// FieldTypeError is used for fields that store error.
	FieldTypeError

	// FieldTypeContext is used for fields that store context.Context.
	FieldTypeContext
)

// String version of FieldType.
func (f FieldType) String() string {
	switch f {
	case FieldTypeReflect:
		return "Reflect"
	case FieldTypeString:
		return "String"
	case FieldTypeInt64:
		return "Int64"
	case FieldTypeFloat64:
		return "Float64"
	case FieldTypeBinary:
		return "Binary"
	case FieldTypeByteString:
		return "ByteString"
	case FieldTypeError:
		return "Error"
	case FieldTypeTimeFull:
		return "TimeFull"
	case FieldTypeTime:
		return "Time"
	case FieldTypeDuration:
		return "Duration"
	case FieldTypeBool:
		return "Bool"
	case FieldTypeContext:
		return "Context"
	case FieldTypeUnknown:
		fallthrough
	default:
		return "Unknown"
	}
}

// Any constructs a field with the given key and value.
func Any(key string, val interface{}) Field { // nolint: cyclop
	if val == nil {
		return nilField(key)
	}

	switch value := val.(type) {
	case string:
		return String(key, value)
	case int:
		return Int(key, value)
	case int64:
		return Int64(key, value)
	case float64:
		return Float64(key, value)
	case []byte:
		return Binary(key, value)
	case bool:
		return Bool(key, value)
	case time.Time:
		return Time(key, value)
	case time.Duration:
		return Duration(key, value)
	case error:
		return NamedError(key, value)
	case context.Context:
		return NamedContext(key, value)
	}

	return Field{Key: key, Type: FieldTypeUnknown, Interface: val}
}

// String constructs a field with the given key and value.
func String(key string, val string) Field {
	return Field{Key: key, Type: FieldTypeString, Str: val}
}

// Int constructs a field that carries an int.
func Int(key string, value int) Field {
	return Int64(key, int64(value))
}

// Int64 constructs a field that carries an int64.
func Int64(key string, value int64) Field {
	return Field{Key: key, Type: FieldTypeInt64, Integer: value}
}

// Binary constructs a field that carries an opaque binary blob.
//
// Binary data is serialized in an encoding-appropriate format. For example,
// zap's JSON encoder base64-encodes binary blobs. To log UTF-8 encoded text,
// use ByteString.
func Binary(key string, val []byte) Field {
	return Field{Key: key, Type: FieldTypeBinary, Interface: val}
}

// Bool constructs a field that carries a bool.
func Bool(key string, val bool) Field {
	var iVal int64
	if val {
		iVal = 1
	}

	return Field{Key: key, Type: FieldTypeBool, Integer: iVal}
}

// nilField returns a field which will marshal explicitly as nil.
func nilField(key string) Field { return Reflect(key, nil) }

// Reflect constructs a field with the given key and an arbitrary object. It uses
// an encoding-appropriate, reflection-based function to lazily serialize nearly
// any object into the logging context, but it's relatively slow and
// allocation-heavy. Outside tests, Any is always a better choice.
//
// If encoding fails (e.g., trying to serialize a map[int]string to JSON), Reflect
// includes the error message in the final log output.
func Reflect(key string, val interface{}) Field {
	return Field{Key: key, Type: FieldTypeReflect, Interface: val}
}

// ByteString constructs a field that carries UTF-8 encoded text as a []byte.
// To log opaque binary blobs (which aren't necessarily valid UTF-8), use
// Binary.
func ByteString(key string, val []byte) Field {
	return Field{Key: key, Type: FieldTypeByteString, Interface: val}
}

// Float64 constructs a field that carries a float64. The way the
// floating-point value is represented is encoder-dependent, so marshaling is
// necessarily lazy.
func Float64(key string, val float64) Field {
	return Field{Key: key, Type: FieldTypeFloat64, Integer: int64(math.Float64bits(val))}
}

// Time constructs a Field with the given key and value. The encoder
// controls how the time is serialized.
func Time(key string, val time.Time) Field {
	if val.Before(minTimeInt64) || val.After(maxTimeInt64) {
		return Field{Key: key, Type: FieldTypeTimeFull, Interface: val}
	}

	return Field{Key: key, Type: FieldTypeTime, Integer: val.UnixNano(), Interface: val.Location()}
}

// Duration constructs a field with the given key and value. The encoder
// controls how the duration is serialized.
func Duration(key string, val time.Duration) Field {
	return Field{Key: key, Type: FieldTypeDuration, Integer: int64(val)}
}

// ErrorField is shorthand for the common idiom NamedError("error", err).
func ErrorField(err error) Field {
	return NamedError("error", err)
}

// NamedError constructs a field that lazily stores err.Error() under the
// provided key. Errors which also implement fmt.Formatter (like those produced
// by github.com/pkg/errors) will also have their verbose representation stored
// under key+"Verbose". If passed a nil error, the field is a no-op.
//
// For the common case in which the key is simply "error", the Error function
// is shorter and less repetitive.
func NamedError(key string, err error) Field {
	return Field{Key: key, Type: FieldTypeError, Interface: err}
}

// ContextField is shorthand for the common idiom NamedContext("ctx", ctx).
func ContextField(ctx context.Context) Field {
	return NamedContext("ctx", ctx)
}

// NamedContext constructs a field that carries a bool.
func NamedContext(key string, ctx context.Context) Field {
	return Field{Key: key, Type: FieldTypeContext, Interface: ctx}
}

// Stack constructs a field that stores a stacktrace of the current goroutine
// under provided key. Keep in mind that taking a stacktrace is eager and
// expensive (relatively speaking); this function both makes an allocation and
// takes about two microseconds.
func Stack(key string) Field {
	return StackSkip(key, 1) // skip Stack
}

// StackDepth is like Stack with support of StacktraceDepth.
func StackDepth(key string, depth int) Field {
	return StackSkipDepth(key, 1, StacktraceDepth(depth)) // skip Stack
}

// StackSkip constructs a field similarly to Stack, but also skips the given
// number of frames from the top of the stacktrace.
func StackSkip(key string, skip int) Field {
	// Returning the stacktrace as a string costs an allocation, but saves us
	// from expanding the zapcore.Field union struct to include a byte slice. Since
	// taking a stacktrace is already so expensive (~10us), the extra allocation
	// is okay.
	return String(key, takeStacktrace(skip+1)) // skip StackSkip
}

// StackSkipDepth is like StackSkip but also support StacktraceDepth.
func StackSkipDepth(key string, skip int, depth StacktraceDepth) Field {
	return String(key, TakeStacktraceDepth(skip+1, depth)) // skip StackSkip
}

// IsNilField check of field is nilField.
func IsNilField(field Field) bool {
	if field.Type != FieldTypeReflect {
		return false
	}

	return field.Interface == nil
}

// GetChainFields finds all filed from error chain.
func GetChainFields(err error) []Field {
	fields := make([]Field, 0)

	for {
		eErr := GetError(err)
		fields = append(fields, eErr.fields...)

		if eErr.cause == nil {
			break
		}

		err = eErr.cause
	}

	return fields
}

// FindFieldInChain finds requested filed from error chain.
func FindFieldInChain(key string, err error) Field {
	for {
		eErr := GetError(err)
		for _, field := range eErr.fields {
			if field.Key == key {
				return field
			}
		}

		if eErr.cause == nil {
			break
		}

		err = eErr.cause
	}

	return nilField(key)
}

// GetFields from passed error.
func GetFields(err error) []Field {
	return GetError(err).fields
}

// GetField find you field based on the key.
func GetField(err error, key string) Field {
	for _, field := range GetError(err).fields {
		if field.Key == key {
			return field
		}
	}

	return nilField(key)
}

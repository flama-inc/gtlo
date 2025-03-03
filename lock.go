package tlo

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	DefaultLockFileMode = 0666
	//
	ErrObjectIsNil                = "object is nil"
	ErrMetadataIsNil              = "metadata is nil"
	ErrInstanceCouldBeInitialized = "instance cloud be initialized"
	//
	ErrfObjectInvalidPath   = "invalid object file path: %s"
	ErrfObjectReadPath      = "object file read error: %s"
	ErrfObjectUnmarshal     = "ohject unmarshal error: %s"
	ErrfObjectMarshal       = "object marshal error: %s"
	ErrfObjectWrite         = "object write error: %s"
	ErrfMetadataKeyNotFound = "metadata key not found: %s"
)

type Lock struct {
	Path              string
	createIfNotExists bool
	object            *LockObject
}

type LockObject struct {
	Lock      bool
	Timestamp int64 // unix micro sec
	Metadata  map[string][]byte
}

type LockOption func(*Lock)

func New(path string, options ...LockOption) (o *Lock) {
	o = &Lock{
		Path:              path,
		createIfNotExists: true,
	}
	for _, option := range options {
		option(o)
	}
	o.NewObject()
	return
}

func WithCreateIfNotExists(b bool) LockOption {
	return func(c *Lock) {
		c.SetCreateIfNotExists(b)
	}
}
func (c *Lock) SetCreateIfNotExists(b bool) {
	c.createIfNotExists = b
}

func (c Lock) IsCreateIfNotExists() bool {
	return c.createIfNotExists
}

func (c *Lock) Unmarshal(b []byte) (e error) {
	if err := bson.Unmarshal(b, c.object); err != nil {
		e = fmt.Errorf(ErrfObjectUnmarshal, err)
		return
	}
	return
}

func (c *Lock) Load() (e error) {
	if len(c.Path) < 1 {
		e = fmt.Errorf(ErrfObjectInvalidPath, c.Path)
		return
	}
	b, err := os.ReadFile(c.Path)
	if err != nil {
		if c.createIfNotExists {
			return c.Save()
		}
		return err
	}

	if err := c.Unmarshal(b); err != nil {
		return err
	}
	return

}

func (c *Lock) SetLock() (e error) {
	if c == nil {
		return errors.New(ErrInstanceCouldBeInitialized)
	}
	if c.object == nil {
		return errors.New(ErrObjectIsNil)
	}
	c.Lock()
	return c.Save()
}

func (c *Lock) SetUnlock() (e error) {
	if c == nil {
		return errors.New(ErrInstanceCouldBeInitialized)
	}
	if c.object == nil {
		return errors.New(ErrObjectIsNil)
	}
	c.Unlock()
	return c.Save()
}

func (c *Lock) GetMetadataAll() (m map[string][]byte, e error) {
	if c == nil {
		e = errors.New(ErrInstanceCouldBeInitialized)
		return
	}
	if c.object == nil {
		e = errors.New(ErrObjectIsNil)
		return
	}
	m = c.object.Metadata
	if m == nil {
		e = errors.New(ErrMetadataIsNil)
	}
	return
}

func (c *Lock) GetMetadata(key string) (b []byte, e error) {
	if c == nil {
		e = errors.New(ErrInstanceCouldBeInitialized)
		return
	}
	if c.object == nil {
		e = errors.New(ErrObjectIsNil)
		return
	}
	m := c.object.Metadata
	if m == nil {
		e = errors.New(ErrMetadataIsNil)
		return
	}

	var ok bool
	b, ok = m[key]
	if !ok {
		e = fmt.Errorf(ErrfMetadataKeyNotFound, key)
		return
	}
	return
}

func (c *Lock) SetMetadata(key string, value []byte) (e error) {
	if c == nil {
		e = errors.New(ErrInstanceCouldBeInitialized)
		return
	}
	if c.object == nil {
		e = errors.New(ErrObjectIsNil)
		return
	}
	if c.object.Metadata == nil {
		c.object.Metadata = make(map[string][]byte)
	}
	c.object.Metadata[key] = value
	return
}

func (o *Lock) NewObject() *LockObject {
	o.object = &LockObject{}
	o.SetTimestampNow()

	return o.object
}

func (o *Lock) Object() *LockObject {
	if o.object == nil {
		o.NewObject()
	}
	return o.object
}

func (o *Lock) Reset() (e error) {
	o.NewObject()
	o.Save()

	return
}

func (o *Lock) TimeCompare(t time.Time) int {
	if o == nil {
		return -1
	}
	if o.object == nil {
		return -1
	}
	return time.UnixMicro(o.object.Timestamp).Compare(t)
}

func (o *Lock) SetTimestampNow() *Lock {
	return o.SetTimestamp(time.Now())
}

func (o *Lock) SetTimestamp(t time.Time) *Lock {
	o.object.Timestamp = t.UnixMicro()
	return o
}

func (o *Lock) Touch() *Lock {
	return o.SetTimestampNow()
}

func (o *Lock) IsLocked() (b bool) {
	if o == nil {
		return
	}
	if o.object == nil {
		return
	}
	return o.object.Lock
}
func (o *Lock) IsUnlocked() (b bool) {
	if o == nil {
		return
	}
	if o.object == nil {
		return
	}
	return !o.IsLocked()
}

func (o *Lock) Lock() *Lock {
	if o == nil {
		return o
	}
	if o.object == nil {
		o.NewObject()
	}
	o.object.Lock = true
	return o
}

func (o *Lock) Unlock() *Lock {
	if o == nil {
		return o
	}
	if o.object == nil {
		o.NewObject()
	}
	o.object.Lock = false
	return o
}

func (o *Lock) Save() (e error) {
	if o.object == nil {
		e = errors.New(ErrObjectIsNil)
		return
	}
	b, err := bson.Marshal(o.object)
	if err != nil {
		e = fmt.Errorf(ErrfObjectMarshal, err.Error())
		return
	}
	if err := os.WriteFile(o.Path, b, DefaultLockFileMode); err != nil {
		e = fmt.Errorf(ErrfObjectWrite, err.Error())
	}

	return
}

func (o *Lock) Dump() (s string) {
	if o.object == nil {
		s = ErrObjectIsNil
		return
	}
	b, _ := json.MarshalIndent(o.object, "", "  ")
	s = string(b)
	return
}

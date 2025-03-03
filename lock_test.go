package tlo

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"gopkg.in/go-playground/assert.v1"
	//
)

func TestNew(t *testing.T) {

	c := New("")

	actual := c.Path
	expected := ""

	assert.Equal(t, actual, expected)
}

func TestLoadLock(t *testing.T) {
	lock := New("")
	if err := lock.Load(); err == nil {
		t.Fatal("expected: invalid args error")
	}

	d := t.TempDir()
	lock = New(d)
	if err := lock.Load(); err == nil {
		t.Fatalf("%s is directory", d)
	}

	f := filepath.Join(d, "file")

	lock = New(f, WithCreateIfNotExists(false))
	if err := lock.Load(); err == nil {
		t.Fatalf("expected: file not found: %s", f)
	}

	lock = New(f)
	if err := lock.Load(); err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}

	os.WriteFile(f, []byte(""), 0666)
	if err := lock.Load(); err == nil {
		t.Fatalf("expected: unmarshal error")
	}

}

func TestSetLock(t *testing.T) {

	lock := New("")
	if err := lock.Load(); err == nil {
		t.Fatal("expected: invalid args error")
	}

	if err := lock.SetLock(); err == nil {
		t.Fatal("expected: error")
	}

	d := t.TempDir()
	f := filepath.Join(d, "file")
	lock = New(f)
	if err := lock.Load(); err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}

	if err := lock.SetLock(); err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}

}

func TestSetUnlock(t *testing.T) {

	lock := New("")
	if err := lock.Load(); err == nil {
		t.Fatal("expected: invalid args error")
	}
	if err := lock.SetUnlock(); err == nil {
		t.Fatal("expected: error")
	}

	d := t.TempDir()
	f := filepath.Join(d, "file")
	lock = New(f)
	if err := lock.Load(); err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	if err := lock.SetUnlock(); err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}

}

func TestLockTimeCompare(t *testing.T) {

	d := t.TempDir()
	f := filepath.Join(d, "file")
	lock := New(f)
	if err := lock.Load(); err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}

	now := time.Now()

	actual := lock.TimeCompare(now)
	expected := -1
	assert.Equal(t, actual, expected)

	lock.SetTimestamp(now) // set new locktime ( = now)
	actual = lock.TimeCompare(now)
	expected = 0
	assert.Equal(t, actual, expected)

	lock.SetTimestampNow() // set new locktime ( > now)
	actual = lock.TimeCompare(now)
	expected = 1
	assert.Equal(t, actual, expected)

}

func TestGetMetadataAll(t *testing.T) {
	d := t.TempDir()
	f := filepath.Join(d, "file")
	lock := New(f)
	if err := lock.Load(); err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	if _, err := lock.GetMetadataAll(); err == nil {
		t.Fatal("error expected")
	}
	lock.SetMetadata("foo", []byte("bar"))
	m, err := lock.GetMetadataAll()
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	actual := m["foo"]
	expected := []byte("bar")
	assert.Equal(t, actual, expected)

}

func TestGetMetadata(t *testing.T) {
	d := t.TempDir()
	f := filepath.Join(d, "file")
	lock := New(f)
	if err := lock.Load(); err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	if _, err := lock.GetMetadata("foo"); err == nil {
		t.Fatal("error expected")
	}

	lock.SetMetadata("foo", []byte("bar"))

	m, err := lock.GetMetadata("foo")
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	actual := m
	expected := []byte("bar")
	assert.Equal(t, actual, expected)

}

func TestSetMetadata(t *testing.T) {
	d := t.TempDir()
	f := filepath.Join(d, "file")
	lock := New(f)
	if err := lock.Load(); err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	if err := lock.SetMetadata("foo", []byte("bar")); err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}

}

func TestNewObject(t *testing.T) {
	d := t.TempDir()
	f := filepath.Join(d, "file")
	lock := New(f)
	if err := lock.Load(); err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	if lock.NewObject() == nil {
		t.Fatal("object is nil")
	}
}

func TestIsLocked(t *testing.T) {
	d := t.TempDir()
	f := filepath.Join(d, "file")
	lock := New(f)
	if err := lock.Load(); err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}

	if err := lock.SetLock(); err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}

	actual := lock.IsLocked()
	expected := true
	assert.Equal(t, actual, expected)

	if err := lock.SetUnlock(); err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	actual = lock.IsLocked()
	expected = false
	assert.Equal(t, actual, expected)

}

func TestReset(t *testing.T) {
	d := t.TempDir()
	f := filepath.Join(d, "file")
	lock := New(f)
	if err := lock.Load(); err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}

	if err := lock.SetLock(); err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}

	if err := lock.Reset(); err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}

	actual := lock.IsLocked()
	expected := false
	assert.Equal(t, actual, expected)

}

func TestSave(t *testing.T) {
	d := t.TempDir()
	f := filepath.Join(d, "file")
	lock := New(f)
	if err := lock.Load(); err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}

	os.Remove(lock.Path)

	if err := lock.Save(); err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}

	if _, err := os.Stat(lock.Path); err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}

}

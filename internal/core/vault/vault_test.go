package vault

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestInit_CreatesMetaDir(t *testing.T) {
	root := filepath.Join(t.TempDir(), "my-vault")

	v, err := Init(root)
	if err != nil {
		t.Fatal(err)
	}

	info, err := os.Stat(v.MetaDir())
	if err != nil {
		t.Fatalf("meta dir missing: %v", err)
	}
	if !info.IsDir() {
		t.Error("meta path is not a directory")
	}
	if !filepath.IsAbs(v.Root()) {
		t.Errorf("Root() = %q, want an absolute path", v.Root())
	}
}

func TestInit_IsIdempotent(t *testing.T) {
	root := t.TempDir()

	if _, err := Init(root); err != nil {
		t.Fatal(err)
	}
	if _, err := Init(root); err != nil {
		t.Fatalf("second Init failed: %v", err)
	}
}

func TestOpen_AfterInit(t *testing.T) {
	root := t.TempDir()

	if _, err := Init(root); err != nil {
		t.Fatal(err)
	}

	v, err := Open(root)
	if err != nil {
		t.Fatalf("Open failed: %v", err)
	}
	if v.Root() == "" {
		t.Error("Root() is empty")
	}
}

func TestOpen_PlainDirectoryIsNotAVault(t *testing.T) {
	_, err := Open(t.TempDir())

	if !errors.Is(err, ErrNotAVault) {
		t.Fatalf("Open error = %v, want ErrNotAVault", err)
	}
}

func TestOpen_MissingDirectoryIsNotAVault(t *testing.T) {
	_, err := Open(filepath.Join(t.TempDir(), "nope"))

	if !errors.Is(err, ErrNotAVault) {
		t.Fatalf("Open error = %v, want ErrNotAVault", err)
	}
}

func TestOpen_MetaPathThatIsAFileIsNotAVault(t *testing.T) {
	root := t.TempDir()

	if err := os.WriteFile(filepath.Join(root, MetaDirName), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	if _, err := Open(root); !errors.Is(err, ErrNotAVault) {
		t.Fatalf("Open error = %v, want ErrNotAVault", err)
	}
}

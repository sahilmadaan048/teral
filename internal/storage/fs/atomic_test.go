package fs

import (
	"os"
	"path/filepath"
	"testing"
)

func assertFileContent(t *testing.T, path, want string) {
	t.Helper()

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	if string(got) != want {
		t.Errorf("content of %s = %q, want %q", path, got, want)
	}
}

func TestWriteFileAtomic_CreatesFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "note.md")

	if err := WriteFileAtomic(path, []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}

	assertFileContent(t, path, "hello")
}

func TestWriteFileAtomic_ReplacesExistingFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "note.md")

	if err := os.WriteFile(path, []byte("old"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := WriteFileAtomic(path, []byte("new"), 0o644); err != nil {
		t.Fatal(err)
	}

	assertFileContent(t, path, "new")
}

func TestWriteFileAtomic_LeavesNoTempFiles(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "note.md")

	for _, content := range []string{"one", "two", "three"} {
		if err := WriteFileAtomic(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 || entries[0].Name() != "note.md" {
		t.Errorf("directory has %d entries, want only note.md: %v", len(entries), entries)
	}
}

func TestWriteFileAtomic_MissingDirectoryFails(t *testing.T) {
	path := filepath.Join(t.TempDir(), "does-not-exist", "note.md")

	if err := WriteFileAtomic(path, []byte("x"), 0o644); err == nil {
		t.Fatal("expected an error when the directory is missing")
	}
}

func TestWriteFileAtomic_FailureCleansUpTempFile(t *testing.T) {
	dir := t.TempDir()

	// Make the target path a directory. Renaming a file over it must
	// fail, which exercises the cleanup path.
	target := filepath.Join(dir, "note.md")
	if err := os.Mkdir(target, 0o755); err != nil {
		t.Fatal(err)
	}

	if err := WriteFileAtomic(target, []byte("x"), 0o644); err == nil {
		t.Fatal("expected an error when the target is a directory")
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Errorf("directory has %d entries after failure, want 1 (temp file leaked): %v", len(entries), entries)
	}
}

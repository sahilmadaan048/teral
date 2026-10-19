// Package vault represents a directory of notes and the hidden metadata
// directory Teral keeps inside it.
package vault

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// MetaDirName is the hidden directory inside a vault that holds Teral's
// own data such as configuration, the index and history.
const MetaDirName = ".teral"

// ErrNotAVault is returned by Open when the directory has no metadata
// directory, so it was never initialized as a vault.
var ErrNotAVault = errors.New("vault: not a teral vault (missing .teral directory)")

// Vault is a handle to an initialized vault directory.
type Vault struct {
	root string
}

// Init creates the vault directory and its metadata directory if they do
// not exist. Calling Init on an existing vault is safe and changes nothing.
func Init(root string) (*Vault, error) {
	abs, err := filepath.Abs(root)
	if err != nil {
		return nil, fmt.Errorf("resolve vault path: %w", err)
	}

	if err := os.MkdirAll(filepath.Join(abs, MetaDirName), 0o750); err != nil {
		return nil, fmt.Errorf("create vault: %w", err)
	}

	return &Vault{root: abs}, nil
}

// Open returns a handle to an existing vault. It fails with ErrNotAVault
// if the directory was never initialized.
func Open(root string) (*Vault, error) {
	abs, err := filepath.Abs(root)
	if err != nil {
		return nil, fmt.Errorf("resolve vault path: %w", err)
	}

	info, err := os.Stat(filepath.Join(abs, MetaDirName))
	if errors.Is(err, os.ErrNotExist) {
		return nil, ErrNotAVault
	}
	if err != nil {
		return nil, fmt.Errorf("inspect vault: %w", err)
	}
	if !info.IsDir() {
		return nil, ErrNotAVault
	}

	return &Vault{root: abs}, nil
}

// Root returns the absolute path of the vault directory.
func (v *Vault) Root() string {
	return v.root
}

// MetaDir returns the absolute path of the vault's metadata directory.
func (v *Vault) MetaDir() string {
	return filepath.Join(v.root, MetaDirName)
}

package fs

import (
	"fmt"
	"io"
	"io/fs"
	"os"

	"github.com/golang-migrate/migrate/v4/source"
)

type FSSource struct {
	fs         fs.FS
	migrations *source.Migrations
}

func (e *FSSource) Open(_ string) (source.Driver, error) {
	return nil, fmt.Errorf("unimplemented: use `WithFiles` to create in code")
}

func (e *FSSource) Close() error {
	return nil
}

func (e *FSSource) First() (version uint, err error) {
	v, ok := e.migrations.First()
	if !ok {
		return 0, &os.PathError{Op: "first", Path: "", Err: os.ErrNotExist}
	}

	return v, nil
}

func (e *FSSource) Prev(version uint) (prevVersion uint, err error) {
	v, ok := e.migrations.Prev(version)
	if !ok {
		return 0, &os.PathError{Op: fmt.Sprintf("prev for version %v", version), Path: "", Err: os.ErrNotExist}
	}

	return v, nil
}

func (e *FSSource) Next(version uint) (nextVersion uint, err error) {
	v, ok := e.migrations.Next(version)
	if !ok {
		return 0, &os.PathError{Op: fmt.Sprintf("next for version %v", version), Path: "", Err: os.ErrNotExist}
	}

	return v, nil
}

func (e *FSSource) ReadUp(version uint) (r io.ReadCloser, identifier string, err error) {
	if m, ok := e.migrations.Up(version); ok {
		f, err := e.fs.Open(m.Raw)
		if err != nil {
			return nil, "", err
		}
		return f, m.Identifier, nil
	}
	return nil, "", &os.PathError{Op: fmt.Sprintf("read version %v", version), Path: "", Err: os.ErrNotExist}
}

func (e *FSSource) ReadDown(version uint) (r io.ReadCloser, identifier string, err error) {
	if m, ok := e.migrations.Down(version); ok {
		f, err := e.fs.Open(m.Raw)
		if err != nil {
			return nil, "", err
		}
		return f, m.Identifier, nil
	}
	return nil, "", &os.PathError{Op: fmt.Sprintf("read version %v", version), Path: "", Err: os.ErrNotExist}
}

func WithFS(files fs.FS) (source.Driver, error) {
	src := &FSSource{
		fs:         files,
		migrations: source.NewMigrations(),
	}
	entries, err := fs.ReadDir(files, ".")
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		m, err := source.DefaultParse(entry.Name())
		if err != nil {
			return nil, err
		}

		if !src.migrations.Append(m) {
			return nil, fmt.Errorf("unable to parse migration file %q", entry.Name())
		}
	}
	return src, nil
}

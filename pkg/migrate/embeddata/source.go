package embeddata

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/golang-migrate/migrate/v4/source"
)

type EmbedDataSource struct {
	data       map[string][]byte
	migrations *source.Migrations
}

func (e *EmbedDataSource) Open(_ string) (source.Driver, error) {
	return nil, fmt.Errorf("unimplemented: use `WithFiles` to create in code")
}

func (e *EmbedDataSource) Close() error {
	return nil
}

func (e *EmbedDataSource) First() (version uint, err error) {
	v, ok := e.migrations.First()
	if !ok {
		return 0, &os.PathError{Op: "first", Path: "", Err: os.ErrNotExist}
	}

	return v, nil
}

func (e *EmbedDataSource) Prev(version uint) (prevVersion uint, err error) {
	v, ok := e.migrations.Prev(version)
	if !ok {
		return 0, &os.PathError{Op: fmt.Sprintf("prev for version %v", version), Path: "", Err: os.ErrNotExist}
	}

	return v, nil
}

func (e *EmbedDataSource) Next(version uint) (nextVersion uint, err error) {
	v, ok := e.migrations.Next(version)
	if !ok {
		return 0, &os.PathError{Op: fmt.Sprintf("next for version %v", version), Path: "", Err: os.ErrNotExist}
	}

	return v, nil
}

func (e *EmbedDataSource) ReadUp(version uint) (r io.ReadCloser, identifier string, err error) {
	if m, ok := e.migrations.Up(version); ok {
		body := e.data[m.Raw]
		return ioutil.NopCloser(bytes.NewReader(body)), m.Identifier, nil
	}
	return nil, "", &os.PathError{Op: fmt.Sprintf("read version %v", version), Path: "", Err: os.ErrNotExist}
}

func (e *EmbedDataSource) ReadDown(version uint) (r io.ReadCloser, identifier string, err error) {
	if m, ok := e.migrations.Down(version); ok {
		body := e.data[m.Raw]
		return ioutil.NopCloser(bytes.NewReader(body)), m.Identifier, nil
	}
	return nil, "", &os.PathError{Op: fmt.Sprintf("read version %v", version), Path: "", Err: os.ErrNotExist}
}

func WithFiles(files map[string][]byte) (source.Driver, error) {
	src := &EmbedDataSource{
		data:       files,
		migrations: source.NewMigrations(),
	}
	for name := range files {
		m, err := source.DefaultParse(name)
		if err != nil {
			return nil, err
		}

		if !src.migrations.Append(m) {
			return nil, fmt.Errorf("unable to parse migration file %q", name)
		}
	}
	return src, nil
}

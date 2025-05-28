package vectorsdb

import (
	"encoding/gob"
	"io"

	"github.com/umk/llmservices/pkg/vectors"
)

func Marshal[V any](w io.Writer, db *Database[V]) error {
	db.repackVectors()

	e := gob.NewEncoder(w)

	if err := e.Encode(&db.header); err != nil {
		return err
	}

	if err := e.Encode(&db.Data); err != nil {
		return err
	}

	if err := vectors.Marshal(w, db.vectors); err != nil {
		return err
	}

	return nil
}

func Unmarshal[V any](r io.Reader) (*Database[V], error) {
	e := gob.NewDecoder(r)

	db := new(Database[V])

	if err := e.Decode(&db.header); err != nil {
		return db, err
	}

	if err := e.Decode(&db.Data); err != nil {
		return db, err
	}

	vectors, err := vectors.Unmarshal(r)
	if err != nil {
		return db, err
	}
	db.vectors = vectors

	return db, nil
}

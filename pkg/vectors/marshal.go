package vectors

import (
	"encoding/gob"
	"errors"
	"io"
)

func Marshal(w io.Writer, vectors *Vectors) error {
	e := gob.NewEncoder(w)

	if err := e.Encode(&vectors.header); err != nil {
		return err
	}

	for _, c := range vectors.chunks {
		if err := e.Encode(c); err != nil {
			return err
		}
	}

	return nil
}

func Unmarshal(r io.Reader) (*Vectors, error) {
	e := gob.NewDecoder(r)

	v := new(Vectors)

	if err := e.Decode(&v.header); err != nil {
		return nil, err
	}

	var chunk *vectorsChunk
	v.chunks = make([]*vectorsChunk, 0)

	for {
		chunk = &vectorsChunk{
			Records: make([]*chunkRecord, 0, v.header.ChunkSize),
		}
		if err := e.Decode(chunk); err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		v.chunks = append(v.chunks, chunk)
	}

	if len(v.chunks) == 0 {
		return nil, errors.New("no chunks found in the encoded data")
	}

	v.currentChunk = chunk
	return v, nil
}

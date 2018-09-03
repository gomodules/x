package types

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

/*
IntHash represents as int64 Generation and string Hash. It is json serialized into <int64>$<hash_string>.
*/
type IntHash struct {
	generation int64
	hash       string
}

func ParseIntHash(s string) (*IntHash, error) {
	idx := strings.IndexRune(s, '$')
	switch {
	case idx <= 0:
		return nil, errors.New("missing generation")
	case idx == len(s)-1:
		return nil, errors.New("missing hash")
	default:
		i, err := strconv.ParseInt(s[:idx], 10, 64)
		if err != nil {
			return nil, err
		}
		h := s[idx+1:]
		return &IntHash{generation: i, hash: h}, nil
	}
}

func NewIntHash(i int64, h string) IntHash { return IntHash{generation: i, hash: h} }

func IntHashForGeneration(i int64) IntHash { return IntHash{generation: i} }

func IntHashForHash(h string) IntHash { return IntHash{hash: h} }

func (m IntHash) Generation() int64 {
	return m.generation
}

func (m IntHash) Hash() string {
	return m.hash
}

func (m *IntHash) Equal(u *IntHash) bool {
	if m == nil {
		return u == nil
	}
	if u == nil { // t != nil
		return false
	}
	if m == u {
		return true
	}
	if m.generation == u.generation {
		return m.hash == u.hash
	}
	return false
}

func (m IntHash) String() string {
	return fmt.Sprintf(`%d$%s`, m.generation, m.hash)
}

func (m *IntHash) MarshalJSON() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	if m.hash == "" {
		return json.Marshal(m.generation)
	}
	return json.Marshal(m.String())
}

func (m *IntHash) UnmarshalJSON(data []byte) error {
	if m == nil {
		return errors.New("jsontypes.IntHash: UnmarshalJSON on nil pointer")
	}

	if data[0] == '"' {
		var s string
		err := json.Unmarshal(data, &s)
		if err != nil {
			return err
		}
		ih, err := ParseIntHash(s)
		if err != nil {
			return err
		}
		*m = *ih
		return nil
	} else if bytes.Equal(data, []byte("null")) {
		return nil
	}

	var i int64
	err := json.Unmarshal(data, &i)
	if err != nil {
		return err
	}
	m.generation = i
	return nil
}

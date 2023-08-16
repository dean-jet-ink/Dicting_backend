package algo

import (
	"math/rand"
	"time"

	"github.com/oklog/ulid"
)

func GenerateULID() (string, error) {
	t := time.Now()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	ulid, err := ulid.New(ulid.Timestamp(t), entropy)
	if err != nil {
		return "", err
	}

	return ulid.String(), nil
}

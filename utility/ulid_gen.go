package utility

import (
	"math/rand"
	"time"

	"github.com/oklog/ulid/v2"
)

var entropy = ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)

func GenID() string {
	t := time.Now().UTC() // Use the current time for a unique timestamp
	return ulid.MustNew(ulid.Timestamp(t), entropy).String()
}

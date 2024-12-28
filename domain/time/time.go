package time

import (
	"cloud.google.com/go/spanner"
	"time"
)

const (
	CreateAt = "CreateAt"
	UpdateAt = "UpdateAt"
)

var offset time.Duration

func Now() time.Time {
	return time.Now().Add(offset).UTC()
}

func SetOffset(d time.Duration) {
	offset = d
}

func CommitTimeStamp() time.Time {
	if offset != 0 {
		return Now() // オフセットがある場合はカスタム時刻
	}
	return spanner.CommitTimestamp
}

type RecordTime struct {
	CreateAt time.Time `spanner:"created_at"`
	UpdateAt time.Time `spanner:"updated_at"`
}

type RecordCreateTime struct {
	CreateAt time.Time `spanner:"created_at"`
}

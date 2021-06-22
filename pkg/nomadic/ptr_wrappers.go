package nomadic

import (
	"time"
)

func Int(v int) *int {
	return &v
}

func String(v string) *string {
	return &v
}

func Duration(v time.Duration) *time.Duration {
	return &v
}

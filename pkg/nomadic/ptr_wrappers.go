package nomadic

import (
	"time"
)

func Bool(v bool) *bool {
	return &v
}

func Int(v int) *int {
	return &v
}

func String(v string) *string {
	return &v
}

func Duration(v time.Duration) *time.Duration {
	return &v
}

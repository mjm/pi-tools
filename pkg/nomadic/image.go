package nomadic

import (
	"strings"
)

func Image(repo, tagOrDigest string) string {
	if strings.HasPrefix(tagOrDigest, "sha256:") {
		return repo + "@" + tagOrDigest
	} else {
		return repo + ":" + tagOrDigest
	}
}

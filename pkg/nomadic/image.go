package nomadic

import (
	_ "embed"
	"sort"
	"strings"
)

func Image(repo, tagOrDigest string) string {
	if tagOrDigest == "latest" {
		if uri, ok := digestMap[repo]; ok {
			return uri
		}

		repoWithRegistry := "index.docker.io/" + repo
		if uri, ok := digestMap[repoWithRegistry]; ok {
			return uri
		}
	}

	if strings.HasPrefix(tagOrDigest, "sha256:") {
		return repo + "@" + tagOrDigest
	} else {
		return repo + ":" + tagOrDigest
	}
}

func RegisteredImageURIs() []string {
	var uris []string
	for _, imageURI := range digestMap {
		uris = append(uris, imageURI)
	}
	sort.Strings(uris)
	return uris
}

//go:embed digests
var allDigests string

var digestMap map[string]string

func init() {
	digestMap = map[string]string{}
	for _, imageURI := range strings.Split(allDigests, "\n") {
		if imageURI == "" {
			continue
		}

		uriParts := strings.Split(imageURI, "@")
		digestMap[uriParts[0]] = imageURI
	}
}

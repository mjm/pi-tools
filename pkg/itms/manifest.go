package itms

type Manifest struct {
	Items []*Item `plist:"items"`
}

type Item struct {
	Assets   []*Asset `plist:"assets"`
	Metadata Metadata `plist:"metadata"`
}

type Asset struct {
	Kind string `plist:"kind"`
	URL  string `plist:"url"`
}

type Metadata struct {
	BundleIdentifier string `plist:"bundle-identifier"`
	BundleVersion    string `plist:"bundle-version"`
	Kind             string `plist:"kind"`
	Title            string `plist:"title"`
}

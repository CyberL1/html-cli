package types

type GithubRelease struct {
	TagName    string `json:"tag_name"`
	Prerelease bool
	Assets     []releaseAsset `json:"assets"`
}

type releaseAsset struct {
	Name               string `json:"name"`
	BrowserDownloadUrl string `json:"browser_download_url"`
}

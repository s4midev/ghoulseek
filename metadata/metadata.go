package metadata

type Track struct {
	Title    string  `json:"title"`
	Duration float64 `json:"duration"`
}
type Release struct {
	Title       string `json:"title"`
	Information string `json:"information"`
	// might switch to time.Time at some point
	ReleaseDate   string  `json:"releaseDate"`
	Tracks        []Track `json:"tracks"`
	MusicBrainzId string  `json:"musicBrainzId"`
	// album, ep, single
	ReleaseType string `json:"releaseType"`
}

type Artist struct {
	Name          string    `json:"name"`
	MusicBrainzId string    `json:"musicBrainzId"`
	Releases      []Release `json:"releases"`
}

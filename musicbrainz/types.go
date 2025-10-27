package musicbrainz

type Artist struct {
	ID             string    `json:"id"`
	Type           string    `json:"type,omitempty"`
	Name           string    `json:"name"`
	SortName       string    `json:"sort-name,omitempty"`
	Disambiguation string    `json:"disambiguation,omitempty"`
	Country        string    `json:"country,omitempty"`
	Area           *Area     `json:"area,omitempty"`
	LifeSpan       *LifeSpan `json:"life-span,omitempty"`
	Aliases        []Alias   `json:"aliases,omitempty"`
	Tags           []Tag     `json:"tags,omitempty"`
}

type Area struct {
	ID       string `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	SortName string `json:"sort-name,omitempty"`
}

type LifeSpan struct {
	Begin string `json:"begin,omitempty"`
	End   string `json:"end,omitempty"`
	Ended *bool  `json:"ended,omitempty"`
}

type Alias struct {
	Name     string `json:"name,omitempty"`
	Locale   string `json:"locale,omitempty"`
	Type     string `json:"type,omitempty"`
	SortName string `json:"sort-name,omitempty"`
}

type Tag struct {
	Name  string `json:"name,omitempty"`
	Count int    `json:"count,omitempty"`
}

type ReleaseGroup struct {
	ID               string   `json:"id"`
	Title            string   `json:"title"`
	FirstReleaseDate string   `json:"first-release-date,omitempty"`
	PrimaryType      string   `json:"primary-type,omitempty"`
	SecondaryTypes   []string `json:"secondary-types,omitempty"`
}

type ReleaseGroupResponse struct {
	ReleaseGroups []ReleaseGroup `json:"release-groups"`
}

type Release struct {
	ID           string        `json:"id"`
	Title        string        `json:"title"`
	Status       string        `json:"status,omitempty"`
	Date         string        `json:"date,omitempty"`
	Country      string        `json:"country,omitempty"`
	ReleaseGroup *ReleaseGroup `json:"release-group,omitempty"`
	Media        []Media       `json:"media,omitempty"`
}

type Media struct {
	Format     string  `json:"format,omitempty"`
	TrackCount int     `json:"track-count,omitempty"`
	Tracks     []Track `json:"tracks,omitempty"`
}

type Track struct {
	ID        string     `json:"id"`
	Number    string     `json:"number,omitempty"`
	Title     string     `json:"title"`
	Length    int        `json:"length,omitempty"` // ms
	Recording *Recording `json:"recording,omitempty"`
}

type Recording struct {
	ID            string         `json:"id"`
	Title         string         `json:"title"`
	Length        int            `json:"length,omitempty"`
	ArtistCredits []ArtistCredit `json:"artist-credit,omitempty"`
}

type ArtistCredit struct {
	Name   string     `json:"name,omitempty"`
	Artist *ArtistRef `json:"artist,omitempty"`
}

type ArtistRef struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

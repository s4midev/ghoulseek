package slsk

// the response of the /searches endpoint, NOT to be confused with `SearchResponses`
type SearchResponse struct {
	FileCount int `json:"fileCount"`
	// passed to searches/{id}/responses for response data
	SearchId   string        `json:"id"`
	Responses  []interface{} `json:"responses"`
	SearchText string        `json:"searchText"`
	// is this used for anything ??
	Token      float64 `json:"token"`
	IsComplete bool    `json:"isComplete"`
}

type File struct {
	BitDepth   int    `json:"bitDepth"`
	Extension  string `json:"extension"`
	FileName   string `json:"filename"`
	Length     int    `json:"length"`
	SampleRate int    `json:"sampleRate"`
	Size       int    `json:"size"`
	IsLocked   bool   `json:"isLocked"`
}
type SearchResponses struct {
	FileCount         int    `json:"fileCount"`
	Files             []File `json:"files"`
	HasFreeUploadSlot bool   `json:"hasFreeUploadSlot"`
	UploadSpeed       int64  `json:"uploadSpeed"`
	Username          string `json:"username"`
}

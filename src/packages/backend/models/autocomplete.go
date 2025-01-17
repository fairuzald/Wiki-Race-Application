package models

// AutoComplete struct to decode JSON response from Wikipedia API
type AutoComplete struct {
	BatchComplete string `json:"batchcomplete"`
	Continue      *struct {
		GpsOffset int    `json:"gpsoffset"`
		Continue  string `json:"continue,omitempty"`
	} `json:"continue,omitempty"`
	Query struct {
		Continue struct {
			PicContinue int    `json:"picontinue"`
			Continue    string `json:"continue,omitempty"`
		}
		Redirects []struct {
			Index int    `json:"index"`
			From  string `json:"from"`
			To    string `json:"to"`
		} `json:"redirects"`
		Pages map[string]Page `json:"pages"`
	} `json:"query"`
}

type Page struct {
	PageID    int `json:"pageid"`
	NS        int `json:"ns"`
	Title     string
	Index     int `json:"index"`
	Thumbnail struct {
		Source string
		Width  int
		Height int
	}
	Terms struct {
		Description []string
	}
	FullURL string `json:"fullurl"`
}

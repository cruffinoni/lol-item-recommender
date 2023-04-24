package gamedata

type ItemData struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Colloq      string   `json:"colloq"`
	Plaintext   string   `json:"plaintext"`
	Into        []string `json:"into"`
	Image       struct {
		Full   string `json:"full"`
		Sprite string `json:"sprite"`
		Group  string `json:"group"`
		X      int    `json:"x"`
		Y      int    `json:"y"`
		W      int    `json:"w"`
		H      int    `json:"h"`
	} `json:"image"`
	Gold struct {
		Base        int  `json:"base"`
		Purchasable bool `json:"purchasable"`
		Total       int  `json:"total"`
		Sell        int  `json:"sell"`
	} `json:"gold"`
	Tags  []string     `json:"tags"`
	Maps  map[int]bool `json:"maps"`
	Stats struct {
		FlatMovementSpeedMod int `json:"FlatMovementSpeedMod"`
	} `json:"stats"`
}

type ItemResponse struct {
	Type    string              `json:"type"`
	Version string              `json:"version"`
	Data    map[string]ItemData `json:"data"`
	Groups  []struct {
		Id              string `json:"id"`
		MaxGroupOwnable string `json:"MaxGroupOwnable"`
	} `json:"groups"`
	Tree []struct {
		Header string   `json:"header"`
		Tags   []string `json:"tags"`
	} `json:"tree"`
}

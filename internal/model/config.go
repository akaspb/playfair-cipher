package model

type Config struct {
	GridConfig GridConfig `json:"grid_config"`
}

type GridConfig struct {
	Chars  []rune `json:"chars"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
	Key    string `json:"key"`
}

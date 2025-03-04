package main

type ErrResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type TrackData struct {
	Ip        string `json:"ip"`
	UserAgent string `json:"user-agent"`
	Url       string `json:"url,omitempty"`
	Timestamp string `json:"timestamp"`
}

type TrackInfo struct {
	Id          string      `json:"id"`
	Name        string      `json:"name"`
	Email       string      `json:"email,omitempty"`
	CreatedAt   string      `json:"created_at"`
	CreatedBy   string      `json:"created_by"`
	PixelEvents []TrackData `json:"pixel_events,omitempty"`
	UrlEvents   []TrackData `json:"url_events,omitempty"`
}

package main

type ErrResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type TrackData struct {
	Ip        string `json:"ip"`
	UserAgent string `json:"user-agent"`
	Timestamp string `json:"timestamp"`
}

type TrackInfo struct {
	Id        string      `json:"id"`
	Name      string      `json:"name"`
	Email     string      `json:"email,omitempty"`
	CreatedAt string      `json:"created_at"`
	CreatedBy string      `json:"created_by"`
	Events    []TrackData `json:"events,omitempty"`
}

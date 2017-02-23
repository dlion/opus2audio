package main

//Structs taken from http://apiv2.online-convert.com/ doc

type JsonResponsePost struct {
	ID     string `json:"id"`
	Token  string `json:"token"`
	Type   string `json:"type"`
	Status struct {
		Code string `json:"code"`
		Info string `json:"info"`
	} `json:"status"`
	Errors     []interface{} `json:"errors"`
	Process    bool          `json:"process"`
	Conversion []struct {
		ID       string `json:"id"`
		Target   string `json:"target"`
		Category string `json:"category"`
		Options  struct {
			Frequency            interface{} `json:"frequency"`
			AudioBitrate         interface{} `json:"audio_bitrate"`
			DownloadPassword     interface{} `json:"download_password"`
			AllowMultipleOutputs bool        `json:"allow_multiple_outputs"`
			Channels             interface{} `json:"channels"`
			Normalize            bool        `json:"normalize"`
			Start                interface{} `json:"start"`
			End                  interface{} `json:"end"`
		} `json:"options"`
	} `json:"conversion"`
	Input []struct {
		ID          string `json:"id"`
		Type        string `json:"type"`
		Source      string `json:"source"`
		Filename    string `json:"filename"`
		Size        int    `json:"size"`
		Hash        string `json:"hash"`
		Checksum    string `json:"checksum"`
		ContentType string `json:"content_type"`
		CreatedAt   string `json:"created_at"`
		ModifiedAt  string `json:"modified_at"`
	} `json:"input"`
	Output       []interface{} `json:"output"`
	Callback     string        `json:"callback"`
	NotifyStatus bool          `json:"notify_status"`
	Server       string        `json:"server"`
	Spent        int           `json:"spent"`
	CreatedAt    string        `json:"created_at"`
	ModifiedAt   string        `json:"modified_at"`
}

type JsonResponseGet []struct {
	ID     string `json:"id"`
	Source struct {
		Conversion string   `json:"conversion"`
		Input      []string `json:"input"`
	} `json:"source"`
	URI              string `json:"uri"`
	Size             int    `json:"size"`
	Status           string `json:"status"`
	ContentType      string `json:"content_type"`
	DownloadsCounter int    `json:"downloads_counter"`
	Checksum         string `json:"checksum"`
	CreatedAt        string `json:"created_at"`
}
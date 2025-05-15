package constants

const CONFIG_VERSION int = 2
const APP_VERSION string = "4.9.1"

var AudioType = struct {
	M4a  string
	Mp3  string
	Flac string
}{M4a: ".m4a", Mp3: ".mp3", Flac: ".flac"}

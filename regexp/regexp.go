package regexp

import "regexp"

var (
	USERSSECID = regexp.MustCompile(`"secUid":"([^"]+)"`)

	VIDEOURL   = regexp.MustCompile(`"playAddr":\s*"([^"]+)"`)
	VIDEOTITLE = regexp.MustCompile(`"desc":"([^"]+)"`)
)

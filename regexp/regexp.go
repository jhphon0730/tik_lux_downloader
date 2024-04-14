package regexp

import "regexp"

var (
	VIDEOURL   = regexp.MustCompile(`"playAddr":\s*"([^"]+)"`)
	VIDEOTITLE = regexp.MustCompile(`"desc":"([^"]+)"`)
)

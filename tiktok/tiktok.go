package tiktok

import (
	"log"
	"regexp"
	"strings"
	"sync"
	"tik_downloader/request"
	"tik_downloader/utils"
)

func GetVideoURL(url string) (string, string, map[string]string) {
	html, video_url_headers, err := request.Get(url, url, nil)
	if err != nil {
		log.Fatalln(err.Error())
	}

	urlMatcherRegExp := regexp.MustCompile(`"downloadAddr":\s*"([^"]+)"`)
	matchDownloadAddr := urlMatcherRegExp.FindStringSubmatch(html)

	if len(matchDownloadAddr) == 0 {
		return "", "", nil
	}

	videoURL := strings.ReplaceAll(matchDownloadAddr[1], `\u002F`, "/")

	titleMatcherRegExp := regexp.MustCompile(`"desc":"([^"]+)"`)
	titleMatcher := titleMatcherRegExp.FindStringSubmatch(html)
	if len(titleMatcher) == 0 {
		return "", "", nil
	}

	title := titleMatcher[1]
	titleArr := strings.Split(title, "|")
	if len(titleArr) == 1 {
		title = titleArr[0]
	} else {
		title = strings.TrimSpace(strings.Join(titleArr[:len(titleArr)-1], "|"))
	}
	return videoURL, title, video_url_headers
}

func GetVideoByteData(url string, video_url_headers map[string]string) ([]byte, error) {
	video_body, _, err := request.GetByte(url, url, video_url_headers)
	if err != nil {
		return nil, err
	}

	return video_body, nil
}

func Download(videoList []string, downloadWg *sync.WaitGroup) {
	downloadWg.Add(len(videoList))

	for _, video := range videoList {
		go func(video string) {
			defer downloadWg.Done()
			video_url, video_title, cookies := GetVideoURL(video)
			if len(video_url) == 0 || len(video_title) == 0 || cookies == nil {
				log.Fatalln("Failed to get video url")
			}

			log.Println("Downloading video: ", video_title)

			video_bytes, err := GetVideoByteData(video_url, cookies)
			if err != nil {
				log.Fatalln(err.Error())
			}

			err = utils.Save(video_bytes, video_title)
			if err != nil {
				log.Fatalln(err.Error())
			}
		}(video)
	}
}

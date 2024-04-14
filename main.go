package main

import (
	"log"
	"sync"
	"tik_downloader/tiktok"
)

func DownloadVideo() {
	var wg sync.WaitGroup
	var videoList []string = []string{"https://www.tiktok.com/@xz_s2_xz/video/7356933138792172816?is_from_webapp=1"}

	tiktok.Download(videoList, &wg)

	wg.Wait()
}

func main() {
	_, _, err := tiktok.GetOriginUserURL("https://www.tiktok.com/@p7apan")
	if err != nil {
		log.Fatalln(err.Error())
	}
}

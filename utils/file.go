package utils

import (
	"bufio"
	"os"
)

func Save(videoBytes []byte, filename string) error {
	file, err := os.Create(filename + ".mp4")
	if err != nil {
		return err
	}
	defer file.Close()

	// bufio.NewWriter를 사용하여 파일 쓰기를 버퍼링합니다.
	bufferedWriter := bufio.NewWriter(file)

	// 비디오 바이트를 파일에 쓰고 버퍼를 비웁니다.
	if _, err := bufferedWriter.Write(videoBytes); err != nil {
		return err // 쓰기 오류 처리
	}
	return bufferedWriter.Flush() // 버퍼 비우기를 확인
}

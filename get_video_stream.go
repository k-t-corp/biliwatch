package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type BiliBiliVideoStreamDataUrl struct {
	Url       string `json:"url"`
	BackupUrl string `json:"backup_url"`
}

type BiliBiliVideoStreamData struct {
	Durl []BiliBiliVideoStreamDataUrl `json:"durl"`
}

type BiliBiliVideoStream struct {
	Code    int                     `json:"code"`
	Message string                  `json:"message"`
	Data    BiliBiliVideoStreamData `json:"data"`
}

type VideoStream struct {
	urls []string
}

func GetVideoStream(bvid string, cid int) (*VideoStream, error) {
	params := url.Values{}
	params.Set("bvid", bvid)
	params.Set("cid", strconv.Itoa(cid))
	// https://github.com/SocialSisterYi/bilibili-API-collect/blob/master/video/videostream_url.md#qn视频清晰度标识
	params.Set("qn", "64")   // 720P
	params.Set("fnval", "1") // mp4 format
	params.Set("fnver", "0")
	params.Set("fourk", "0") // no 4K
	endpoint := "https://api.bilibili.com/x/player/playurl?" + params.Encode()

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	log.Println(fmt.Sprintf("getVideoStream %s %v %s", bvid, cid, res.Status))
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var bVideoStream BiliBiliVideoStream
	err = json.Unmarshal(body, &bVideoStream)
	if err != nil {
		return nil, err
	}

	if bVideoStream.Code != 0 {
		return nil, errors.New(bVideoStream.Message)
	}

	var urls []string
	for _, durl := range bVideoStream.Data.Durl {
		urls = append(urls, durl.Url)
	}

	return &VideoStream{urls: urls}, nil
}

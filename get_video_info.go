package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type BiliBiliVideoInfoData struct {
	Title string `json:"title"`
	Cid   int    `json:"cid"`
}

type BiliBiliVideoInfo struct {
	Code    int                   `json:"code"`
	Message string                `json:"message"`
	Data    BiliBiliVideoInfoData `json:"data"`
}

type VideoInfo struct {
	Title string
	Cid   int
}

func GetVideoInfo(bvid string) (*VideoInfo, error) {
	params := url.Values{}
	params.Set("bvid", bvid)
	endpoint := "https://api.bilibili.com/x/web-interface/view?" + params.Encode()

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	log.Println(fmt.Sprintf("getVideoInfo %s %s", bvid, res.Status))
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var bVideoInfo BiliBiliVideoInfo
	err = json.Unmarshal(body, &bVideoInfo)
	if err != nil {
		return nil, err
	}

	if bVideoInfo.Code != 0 {
		return nil, errors.New(bVideoInfo.Message)
	}

	return &VideoInfo{
		Title: bVideoInfo.Data.Title,
		Cid:   bVideoInfo.Data.Cid,
	}, nil
}

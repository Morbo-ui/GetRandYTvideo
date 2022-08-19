package youtube

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

const YOUTUBE_SEARCH_URL = "https://www.googleapis.com/youtube/v3/search"
const YOUTUBE_API_TOKEN = ""
const YOUTUBE_VIDEO_URL = "https://www.youtube.com/watch?v="

func GetLastVideo(channelUrl string, count int) (string, error) {
	items, err := retrieveVideos(channelUrl)
	if err != nil {
		return "", err
	}
	if len(items) < 1 {
		return "", errors.New("No video found")
	}
	return YOUTUBE_VIDEO_URL + items[count].Id.VideoId, nil
}

func retrieveVideos(channelUrl string) ([]Item, error) {
	req, err := makeRequest(channelUrl, 20)
	if err != nil {
		return nil, err
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var restResponse RestResponse
	err = json.Unmarshal(body, &restResponse)
	if err != nil {
		return nil, err
	}
	return restResponse.Items, nil
}

func makeRequest(channelUrl string, maxResults int) (*http.Request, error) {
	lastSlashIndex := strings.LastIndex(channelUrl, "/")
	channelId := channelUrl[lastSlashIndex+1:]
	//channelId := "UCsJn2vDKhf3YrCnSoKRv0lw"

	req, err := http.NewRequest("GET", YOUTUBE_SEARCH_URL, nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	query.Add("part", "id")
	query.Add("channelId", channelId)
	query.Add("maxResults", strconv.Itoa(maxResults))
	query.Add("order", "date")
	query.Add("key", YOUTUBE_API_TOKEN)
	req.URL.RawQuery = query.Encode()
	fmt.Println(req.URL.String())
	return req, nil
}

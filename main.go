package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

func downloadVideo(videoURL string, apiKey string, startTime string, endTime string) error {
	// Parse the video URL
	u, err := url.Parse(videoURL)
	if err != nil {
		return fmt.Errorf("error parsing URL: %v", err)
	}

	// Extract the video ID from the URL
	query, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return fmt.Errorf("error parsing query string: %v", err)
	}
	videoID := query.Get("v")
	if videoID == "" {
		return fmt.Errorf("invalid video URL")
	}

	// Send a GET request to the YouTube Data API to search for the video
	resp, err := http.Get(fmt.Sprintf("https://www.googleapis.com/youtube/v3/search?key=%s&type=video&q=%s&part=id,snippet", apiKey, url.QueryEscape(videoURL)))
	if err != nil {
		return fmt.Errorf("error sending GET request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %v", err)
	}

	// Parse the JSON response to extract the video ID
	var data struct {
		Items []struct {
			ID struct {
				VideoID string `json:"videoId"`
			} `json:"id"`
		} `json:"items"`
	}
	if err := json.Unmarshal(body, &data); err != nil {
		return fmt.Errorf("error parsing JSON response: %v", err)
	}

	// Check if the video was found
	if len(data.Items) == 0 {
		return fmt.Errorf("video not found")
	}

	// Get the video ID
	videoID = data.Items[0].ID.VideoID

	// Send a GET request to the YouTube Data API to retrieve information about the video
	resp, err = http.Get(fmt.Sprintf("https://www.googleapis.com/youtube/v3/videos?id=%s&key=%s&part=contentDetails,statistics&startTime=%s&endTime=%s", videoID, apiKey, startTime, endTime))
	if err != nil {
		return fmt.Errorf("error sending GET request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %v", err)
	}
	defer resp.Body.Close()

	// Create a file to save the video
	file, err := os.Create("video.mp4")
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	// Write the video data to the file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}

	return nil
}

err := downloadVideo("VIDEO_ID", "YOUR_API_KEY", "00:00:10", "00:00:20")
if err != nil {
	// handle error
}
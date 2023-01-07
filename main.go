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
	file, err := os.Create(fmt.Sprintf("%s.mp4",videoID))
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

func uploadVideo(filePath string) (*youtube.Video, error) {
	ctx := context.Background()

	// Set up the YouTube Data API client
	client, err := youtube.NewService(ctx, option.WithAPIKey("YOUR_API_KEY"))
	if err != nil {
		return nil, fmt.Errorf("Error creating YouTube client: %v", err)
	}

	// Read the file that you want to upload
	videoBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("Error reading file: %v", err)
	}

	// Set up the metadata for the video
	video := &youtube.Video{
		Snippet: &youtube.VideoSnippet{
			Title:       "My Video",
			Description: "This is my video",
			CategoryId:  "22",
		},
		Status: &youtube.VideoStatus{
			PrivacyStatus: "private",
		},
	}

	// Create a request to insert the video
	insertCall := client.Videos.Insert("snippet,status", video)

	// Set the file to be uploaded
	media := &youtube.Media{
		Body: bytes.NewReader(videoBytes),
	}
	insertCall.Media(media)

	// Execute the request and return the response
	response, err := insertCall.Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("Error executing request: %v", err)
	}
	return response, nil
}

func main() {
	videoURL := "URL"
	videoID := "ID" //TODO: parseID from url (v=) and get ride of the parsing in the download video, use id instead
	apiKey := ="KEY"
	startTime := "00:00:10"
	endTime := "00:00:20"
	err := downloadVideo(videoURL, apiKey, startTime, endTime)
	if err != nil {
		// handle error
	}

	// Upload a video
	video, err := uploadVideo(fmt.Sprintf("%s.mp4",videoID)) //TODO: pass the metadata and media info as a param struct
	if err != nil {
		log.Fatalf("Error uploading video: %v", err)
	}
	fmt.Printf("Uploaded video: %v\n", video.Id)
}


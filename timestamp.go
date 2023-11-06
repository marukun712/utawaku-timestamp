package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type TimeStamp struct {
	Time  string `json:"time"`
	Title string `json:"title"`
}

type VideoDetails struct {
	Items TimeStamp `json:"items"`
	Title string    `json:"title"`
	ID    string    `json:"id"`
	Image string    `json:"image"`
}

func main() {
	ctx := context.Background()
	youtubeService, err := youtube.NewService(ctx, option.WithAPIKey(""))

	if err != nil {
		log.Fatal(err)
	}

	videoIDs := fetchVideoIDs("PLhu18ozRJ5d3XLfoeUw6WQxGJAydvC-iL", youtubeService)

	comments := fetchVideoComments(videoIDs, youtubeService)

	checkTimeStamp(comments)
}

// プレイリストから動画IDをすべて取得
func fetchVideoIDs(playlistID string, service *youtube.Service) (videoIDs []string) {
	playlistsCall := service.PlaylistItems.List([]string{"contentDetails"}).PlaylistId(playlistID).MaxResults(50)

	playlistsResponse, err := playlistsCall.Do()
	if err != nil {
		log.Fatalf("Error call YouTube API: %v", err)
	}

	for _, item := range playlistsResponse.Items {
		videoIDs = append(videoIDs, item.ContentDetails.VideoId)
	}

	if playlistsResponse.NextPageToken != "" {
		nextPageToken := playlistsResponse.NextPageToken
		for {
			nextCall := service.PlaylistItems.List([]string{"contentDetails"}).PlaylistId(playlistID).PageToken(nextPageToken).MaxResults(50)
			nextResponse, err := nextCall.Do()
			if err != nil {
				log.Fatalf("Error call YouTube API: %v", err)
			}
			for _, nextItem := range nextResponse.Items {
				videoIDs = append(videoIDs, nextItem.ContentDetails.VideoId)
			}
			nextPageToken = nextResponse.NextPageToken
			if nextPageToken == "" {
				break
			}
		}
	}

	return videoIDs
}

// 動画IDからコメント情報を取得
func fetchVideoComments(videoIDs []string, service *youtube.Service) (videoComments []string) {
	for _, id := range videoIDs {
		videoCommentsCall := service.CommentThreads.List([]string{"id", "snippet"}).VideoId(id).MaxResults(100)

		videoCommentsResponse, err := videoCommentsCall.Do()
		if err != nil {
			log.Fatalf("Error call YouTube API: %v", err)
		}

		for _, item := range videoCommentsResponse.Items {
			text := item.Snippet.TopLevelComment.Snippet.TextOriginal

			videoComments = append(videoComments, text)
		}
	}

	return videoComments
}

func checkTimeStamp(videoComments []string) {
	for _, item := range videoComments {
		pattern := `[0-9]{1,}:[0-9]{1,}:[0-9]{1,}(.*)|[0-9]{1,}:[0-9]{1,}(.*)`

		// 正規表現をコンパイル
		re, err := regexp.Compile(pattern)
		if err != nil {
			fmt.Println("Error compiling regex:", err)
			return
		}

		// マッチングを行う
		matches := re.FindAllString(item, -1)

		if len(matches) > 5 {
			for _, item := range matches {
				pattern := `[0-9]{1,}:[0-9]{1,}:[0-9]{1,}|[0-9]{1,}:[0-9]{1,}`

				re, err := regexp.Compile(pattern)
				if err != nil {
					fmt.Println("Error compiling regex:", err)
					return
				}
				timestamp := re.FindAllString(item, -1)
				title := re.ReplaceAllString(item, "")

				t := TimeStamp{strings.TrimSpace(timestamp[0]), strings.TrimSpace(title)}
				bytes, err := json.Marshal(t)
				if err != nil {
					log.Fatal(err)
				}

				fmt.Println(string(bytes))
			}
		}
	}
}

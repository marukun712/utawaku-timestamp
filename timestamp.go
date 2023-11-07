package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

// タイムスタンプオブジェクト
type TimeStamp struct {
	Time  string `json:"time"`
	Title string `json:"title"`
}

// 動画の詳細情報オブジェクト
type VideoDetail struct {
	TimeStamp []TimeStamp `json:"timeStamps"`
	Title     string      `json:"title"`
	ID        string      `json:"id"`
	Image     string      `json:"image"`
}

func main() {
	//Youtube Serviceの作成
	ctx := context.Background()
	youtubeService, err := youtube.NewService(ctx, option.WithAPIKey(""))

	if err != nil {
		log.Fatal(err)
	}

	//プレイリストIDからvideoIDを全て取得
	videoIDs := fetchVideoIDs("PLhu18ozRJ5d3XLfoeUw6WQxGJAydvC-iL", youtubeService)

	//videoIDから動画の詳細情報とタイムスタンプを取得
	details := fetchVideoDetails(videoIDs, youtubeService)

	//JSONファイルに書き込み
	file, err := os.Create("data.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(details); err != nil {
		log.Fatal(err)
	}
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

	//NextPageTokenがあれば次ページのデータも取得
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

// 動画IDから詳細情報を取得
func fetchVideoDetails(videoIDs []string, service *youtube.Service) (videoDetails []VideoDetail) {
	for _, id := range videoIDs {
		var videoComments []string

		//動画ごとにコメントを取得
		videoCommentsCall := service.CommentThreads.List([]string{"id", "snippet"}).VideoId(id).MaxResults(100).Order("relevance")

		videoCommentsResponse, err := videoCommentsCall.Do()
		if err != nil {
			log.Fatalf("Error call YouTube API: %v", err)
		}

		//動画の詳細情報を取得
		videosCall := service.Videos.List([]string{"snippet", "statistics"}).Id(id)
		response, err := videosCall.Do()
		if err != nil {
			log.Fatalf("Error call YouTube API: %v", err)
		}
		item := response.Items[0]

		//動画タイトル、サムネイルURL
		title := item.Snippet.Title
		image := "https://img.youtube.com/vi/" + id + "/maxresdefault.jpg"

		for _, item := range videoCommentsResponse.Items {
			text := item.Snippet.TopLevelComment.Snippet.TextOriginal

			videoComments = append(videoComments, text)
		}

		//タイムスタンプコメントを抽出
		timeStamp := getTimeStamp(videoComments)

		//VideoDetailオブジェクトの作成
		d := VideoDetail{timeStamp, title, id, image}

		videoDetails = append(videoDetails, d)
	}

	return videoDetails
}

// コメントから正規表現でタイムスタンプコメントを抽出
func getTimeStamp(videoComments []string) (timeStamps []TimeStamp) {
	for _, item := range videoComments {
		//タイムスタンプコメントの全文を取得
		pattern := `[0-9]{1,}:[0-9]{1,}:[0-9]{1,}(.*)|[0-9]{1,}:[0-9]{1,}(.*)`

		//コンパイル
		re, err := regexp.Compile(pattern)
		if err != nil {
			fmt.Println("Error compiling regex:", err)
		}

		//マッチング
		matches := re.FindAllString(item, -1)

		//タイムスタンプが5件以上なら
		if len(matches) > 5 {
			for _, item := range matches {
				//タイムスタンプのみを取得
				pattern := `[0-9]{1,}:[0-9]{1,}:[0-9]{1,}|[0-9]{1,}:[0-9]{1,}`

				//コンパイル
				re, err := regexp.Compile(pattern)
				if err != nil {
					fmt.Println("Error compiling regex:", err)
				}
				//マッチング
				time := re.FindAllString(item, -1)

				//タイムスタンプ部分を削除してタイトルを取得
				title := re.ReplaceAllString(item, "")

				//TimeStampオブジェクトの作成
				t := TimeStamp{strings.TrimSpace(time[0]), strings.TrimSpace(title)}

				timeStamps = append(timeStamps, t)
			}

			break
		}
	}

	return timeStamps
}

package engine

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type lastFMRecentTracksResponse struct {
	RecentTracks struct {
		Track []lastFMTrack `json:"track"`
	} `json:"recenttracks"`
}

type lastFMTrack struct {
	Name   string `json:"name"`
	Artist struct {
		Text string `json:"#text"`
	} `json:"artist"`
}

func fetchRecentTracks() []Track {
	username := strings.TrimSpace(os.Getenv("LASTFM_USERNAME"))
	apiKey := strings.TrimSpace(os.Getenv("LASTFM_API_KEY"))
	if username == "" || apiKey == "" {
		return nil
	}

	requestURL := fmt.Sprintf(
		"https://ws.audioscrobbler.com/2.0/?method=user.getrecenttracks&user=%s&api_key=%s&format=json&limit=5",
		url.QueryEscape(username),
		url.QueryEscape(apiKey),
	)

	client := &http.Client{Timeout: 5 * time.Second}
	response, err := client.Get(requestURL)
	if err != nil {
		return nil
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil
	}

	tracks, err := parseRecentTracks(response.Body)
	if err != nil {
		return nil
	}

	return tracks
}

func parseRecentTracks(reader io.Reader) ([]Track, error) {
	var payload lastFMRecentTracksResponse
	if err := json.NewDecoder(reader).Decode(&payload); err != nil {
		return nil, err
	}

	tracks := make([]Track, 0, len(payload.RecentTracks.Track))
	for _, item := range payload.RecentTracks.Track {
		artist := strings.TrimSpace(item.Artist.Text)
		name := strings.TrimSpace(item.Name)
		if artist == "" || name == "" {
			continue
		}

		tracks = append(tracks, Track{Artist: artist, Name: name})
		if len(tracks) == 5 {
			break
		}
	}

	return tracks, nil
}

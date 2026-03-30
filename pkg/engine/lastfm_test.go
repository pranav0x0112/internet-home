package engine

import (
	"strings"
	"testing"
)

func TestParseRecentTracks(t *testing.T) {
	reader := strings.NewReader(`{
  "recenttracks": {
    "track": [
      {"name": "Song One", "artist": {"#text": "Artist One"}},
      {"name": "Song Two", "artist": {"#text": "Artist Two"}},
      {"name": "Song Three", "artist": {"#text": "Artist Three"}},
      {"name": "Song Four", "artist": {"#text": "Artist Four"}},
      {"name": "Song Five", "artist": {"#text": "Artist Five"}},
      {"name": "Song Six", "artist": {"#text": "Artist Six"}}
    ]
  }
}`)

	tracks, err := parseRecentTracks(reader)
	if err != nil {
		t.Fatalf("parseRecentTracks() returned error: %v", err)
	}

	if len(tracks) != 5 {
		t.Fatalf("expected 5 tracks, got %d", len(tracks))
	}

	if tracks[0].Artist != "Artist One" || tracks[0].Name != "Song One" {
		t.Fatalf("unexpected first track: %+v", tracks[0])
	}

	if tracks[4].Artist != "Artist Five" || tracks[4].Name != "Song Five" {
		t.Fatalf("unexpected fifth track: %+v", tracks[4])
	}
}

func TestFetchRecentTracksWithoutCredentialsSkipsWidget(t *testing.T) {
	t.Setenv("LASTFM_USERNAME", "")
	t.Setenv("LASTFM_API_KEY", "")

	tracks := fetchRecentTracks()
	if len(tracks) != 0 {
		t.Fatalf("expected no tracks when credentials are missing, got %d", len(tracks))
	}
}

(function () {
    const LASTFM_USER = "prawns_baka";
    const LASTFM_API_KEY = "986d6f61bbc6c1438ca6ac2c481857de";

    var widget = document.querySelector('.listening-now');
    var list = document.getElementById('lastfm-tracks');

    if (!widget || !list) return;

    if (!LASTFM_API_KEY) {
        widget.hidden = true;
        return;
    }

    var apiURL = 'https://ws.audioscrobbler.com/2.0/?method=user.getrecenttracks&user=' + encodeURIComponent(LASTFM_USER) + '&api_key=' + encodeURIComponent(LASTFM_API_KEY) + '&format=json&limit=5';

    fetch(apiURL, { cache: 'no-store' })
        .then(function (response) {
            if (!response.ok) {
                throw new Error('Last.fm request failed');
            }
            return response.json();
        })
        .then(function (data) {
            var tracks = (((data || {}).recenttracks || {}).track) || [];
            if (!Array.isArray(tracks)) {
                tracks = [tracks];
            }

            var items = tracks
                .map(function (track) {
                    var artist = (track && track.artist && track.artist['#text']) ? track.artist['#text'].trim() : '';
                    var name = track && track.name ? track.name.trim() : '';
                    var art = track && track.image && track.image[2] && track.image[2]['#text'] ? track.image[2]['#text'].trim() : '';
                    var url = track && track.url ? track.url.trim() : '';
                    var isNowPlaying = track && track['@attr'] && track['@attr'].nowplaying === 'true';
                    if (!artist || !name) return null;
                    return { artist: artist, name: name, art: art, url: url, isNowPlaying: isNowPlaying };
                })
                .filter(Boolean)
                .slice(0, 5);

            if (!items.length) {
                widget.hidden = true;
                return;
            }

            list.textContent = '';
            items.forEach(function (track) {
                var row = document.createElement('a');
                row.className = 'lastfm-track';
                if (track.isNowPlaying) {
                    row.classList.add('lastfm-nowplaying');
                }
                row.href = track.url || '#';
                row.target = '_blank';
                row.rel = 'noopener noreferrer';

                if (track.art) {
                    var img = document.createElement('img');
                    img.className = 'lastfm-art';
                    img.src = track.art;
                    img.alt = '';
                    row.appendChild(img);
                }

                var text = document.createElement('span');
                text.className = 'lastfm-text';
                text.textContent = (track.isNowPlaying ? '▶ ' : '') + track.artist + ' - ' + track.name;
                row.appendChild(text);

                list.appendChild(row);
            });
        })
        .catch(function () {
            widget.hidden = true;
        });
})();

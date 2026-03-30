(function () {
    var script = document.currentScript;
    if (!script) return;

    var apiKey = (script.dataset.lastfmKey || '').trim();
    var widget = document.querySelector('.listening-now');
    var list = document.getElementById('lastfm-tracks');

    if (!widget || !list) return;

    if (!apiKey) {
        widget.hidden = true;
        return;
    }

    var apiURL = 'https://ws.audioscrobbler.com/2.0/?method=user.getrecenttracks&user=prawns_baka&api_key=' + encodeURIComponent(apiKey) + '&format=json&limit=5';

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
                    if (!artist || !name) return null;
                    return { artist: artist, name: name };
                })
                .filter(Boolean)
                .slice(0, 5);

            if (!items.length) {
                widget.hidden = true;
                return;
            }

            list.textContent = '';
            items.forEach(function (track) {
                var li = document.createElement('li');
                li.textContent = track.artist + ' — ' + track.name;
                list.appendChild(li);
            });
        })
        .catch(function () {
            widget.hidden = true;
        });
})();

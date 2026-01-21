const PLAYLISTS_URL = window.APP_CONFIG.apiBase + "/video/";
const SEGMENT_URL = window.APP_CONFIG.apiBase + "/segment/";

const videoPlayer = document.getElementById("video");
const videoIdInput = document.getElementById("videoId");

class Player {
    constructor(videoIdInput, videoPlayer) {
        this.videoPlayer = videoPlayer;
        this.videoIdInput = videoIdInput;
        this.playlist = null;
    }

    async getPlaylist() {
        const headers = new Headers();
        headers.append("Accept", "application/vnd.apple.mpegurl");

        const playlistReq = await fetch(PLAYLISTS_URL + this.videoIdInput.value, {
            method: "GET",
            headers: headers
        });

        const playlistText = await playlistReq.text();
        this.playlist = new Playlist(playlistText, this.videoIdInput.value);
        this.playlist.parseSegments();
    }

    async processVideo() {
        const mimeCodec = 'video/mp2t; codecs="avc1.640028, mp4a.40.2"';
        if (!window.MediaSource) {
            console.error("MediaSource API not supported");
            return;
        }

        const mediaSource = new MediaSource();
        this.videoPlayer.src = URL.createObjectURL(mediaSource);

        mediaSource.addEventListener('sourceopen', async () => {
            const sourceBuffer = mediaSource.addSourceBuffer(mimeCodec);

            await this.getPlaylist();
            const segments = this.playlist.getSegments();

            if (segments.length === 0) {
                console.error("No video segments found");
                return;
            }

            const appendNextSegment = async (index) => {
                if (index >= segments.length) {
                    console.log("All segments appended");
                    return;
                }

                const segmentBytes = await segments[index].getSegmentBytes();
                sourceBuffer.appendBuffer(segmentBytes);

                sourceBuffer.addEventListener("updateend", () => {
                    appendNextSegment(index + 1);
                }, { once: true });
            };

            appendNextSegment(0);
        });
    }
}

class Playlist {
    constructor(playlist, videoId) {
        this.playlist = playlist;
        this.segments = [];
        this.videoId = videoId;
    }

    parseSegments() {
        this.segments = [];  // Clear previous segments

        for (const line of this.playlist.split("\n")) {
            if (line.startsWith("#") || line.trim() === "") {
                continue;
            }
            const segmentUrl = SEGMENT_URL + this.videoId + "/" + line.split(".")[0];
            this.segments.push(new Segment(segmentUrl));
        }
    }

    getSegments() {
        return this.segments;
    }
}

class Segment {
    constructor(segmentURL) {
        this.segmentURL = segmentURL;
    }

    async getSegmentBytes() {
        const segmentReq = await fetch(this.segmentURL);
        const segmentBlob = await segmentReq.blob();
        return segmentBlob.arrayBuffer();
    }
}

const player = new Player(videoIdInput, videoPlayer);


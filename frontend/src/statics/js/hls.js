class Loader extends Hls.DefaultConfig.loader {
    constructor(config) {
        super(config);
        this.config = config;
    }

    load(context, config, callbacks) {
        const url = window.location.href.split("/")

        if (context.type === "manifest") {
            console.log("Manifest is being requested.");
        } else if (context.type === "segment") {
            console.log("Segment is being requested.");
        }

        // Rewrite segment URLs if needed
        if (context.type === "fragment" || context.type === "segment" || context.type == undefined) {
            const segmentPath = new URL(context.url).pathname.split('/').pop();
            context.url = `http://0.0.0.0:8081/segment/${url[url.length - 1]}/${segmentPath.split(".")[0]}`;
            console.log(`Rewritten URL: ${context.url}`);
        }

        // Pass the request to the parent loader
        super.load(context, config, callbacks);
    }
}

document.addEventListener("DOMContentLoaded", () => {

    if (Hls.isSupported()) {
        const video = document.getElementById('video');
        const hls = new Hls({
            loader: Loader
        });
        const url = window.location.href.split("/")
        hls.loadSource('http://0.0.0.0:8081/video/' + url[url.length - 1]); // HLS playlist that references .ts files
        hls.attachMedia(video);
    }
})

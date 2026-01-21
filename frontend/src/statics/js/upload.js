const CHUNK_SIZE = 1 * 1024 * 1024

const fileInput = document.getElementById("file-input")
const titleInput = document.getElementById("title-input");
const sendStatus = document.getElementById("sendStatus-bar")
const proccesStatus = document.getElementById("proccesStatus-bar")

function sleep(ms) {
    return new Promise(res => setTimeout(res, ms));
}

function arrayBufferToBase64(buffer) {
    let binary = '';
    const bytes = new Uint8Array(buffer);
    const len = bytes.byteLength;
    for (let i = 0; i < len; i++) {
        binary += String.fromCharCode(bytes[i]);
    }
    return window.btoa(binary); // Use btoa to encode to base64
}

class VideoFile {
    constructor(fileInput, sendStatus, proccesStatus) {
        this.fileInput = fileInput
        this.sendStatus = sendStatus
        this.proccesStatus = proccesStatus

        this.chunks = []
        this.totalChunks = 0

        this.uploaded = 0
        this.sended = 0

        this.uploadId = 0
    }

    async proccesVideo() {
        const file = this.fileInput.files[0]
        const reader = new FileReader()

        reader.readAsArrayBuffer(file)
        this.totalChunks = Math.max(1, Math.ceil(file.size / CHUNK_SIZE))

        const headers = new Headers();
        headers.append("Content-type", "application/json")

        fetch("http://0.0.0.0:8081/uploadMetadata", {
            method: "POST",
            body: JSON.stringify(
                {
                    "title": titleInput.value,
                    "fileSize": file.size,
                    "chunks": this.totalChunks
                }
            ),
            headers: headers
        }).then(res => res.json())
            .then(res => {
                this.uploadId = res.id
                console.log("ID: ", this.uploadId)
            })

        reader.onloadend = async (e) => {
            const arrayBuffer = e.target.result
            const byteArray = new Uint8Array(arrayBuffer)

            for (let i = 0; i < file.size; i += CHUNK_SIZE) {
                var start = 0
                var end = 0

                if ((i + CHUNK_SIZE) > file.size) {
                    start = i
                    end = file.size
                } else {
                    start = i
                    end = i + CHUNK_SIZE
                }

                this.uploaded += end - start

                const chunk = byteArray.slice(start, end)
                this.chunks.push(chunk)
                this.updateProccesBar()
            }

            if (this.uploadId != 0) {
                this.sendVideo()
            } else {
                await sleep(1000)
                this.sendVideo()
            }
        }
    }

    async sendVideo() {
        const headers = new Headers();
        headers.append("Content-type", "application/json")

        let activeUploads = 0;
        const uploadQueue = [...this.chunks]
        const promises = []

        while (uploadQueue.length > 0 || activeUploads > 0) {
            if (activeUploads < 5 && uploadQueue.length > 0) {
                const chunk = uploadQueue.shift()
                const msg = {
                    id: this.uploadId,
                    chunkId: this.sended,
                    chunk: arrayBufferToBase64(chunk)
                }

                activeUploads++
                const promise = sendVideoServer(msg, headers).finally(() => {
                    activeUploads--;
                })
                promises.push(promise)

                this.sended++;
                this.updateSendBar()
            }

            await sleep(10)
        }

        await Promise.all(promises)
    }

    updateProccesBar() {
        let value = (this.chunks.length / this.totalChunks) * 100
        this.proccesStatus.style.width = `${Math.min(value, 100)}%`
    }

    updateSendBar() {
        let value = (this.sended / this.chunks.length) * 100
        this.sendStatus.style.width = `${Math.min(value, 100)}%`
    }
}

const video = new VideoFile(fileInput, sendStatus, proccesStatus)

async function sendVideoServer(msg, headers) {
    fetch("http://0.0.0.0:8081/uploadVideoChunk", {
        method: "POST",
        body: JSON.stringify(msg),
        headers: headers
    })
}

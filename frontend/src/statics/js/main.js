const videoElement = document.getElementById("videos");

document.addEventListener("DOMContentLoaded", () => {
    fetch(window.APP_CONFIG.apiBase + "/videos").then((res) => res.json()).then((data) => {
        for (let i = 0; i < data.length; i++) {
            const element = document.createElement("a")
            element.textContent = data[i].title
            element.href = window.APP_CONFIG.apiBase + "/videos/" + data[i].id

            videoElement.appendChild(element)
        }
    }).catch((err) => console.log(err))
})

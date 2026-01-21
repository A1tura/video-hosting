const videoElement = document.getElementById("videos");


document.addEventListener("DOMContentLoaded", () => {
    fetch("http://127.0.0.1:8081/videos").then((res) => res.json()).then((data) => {
        for (let i = 0; i < data.length; i++) {
            const element = document.createElement("a")
            element.textContent = data[i].title
            element.href = "http://0.0.0.0:8080/videos/" + data[i].id

            videoElement.appendChild(element)
        }
    }).catch((err) => console.log(err))
})


console.log("PixiJS loaded successfully!");

const app = new PIXI.Application({
    width: window.innerWidth,
    height: window.innerHeight,
    backgroundColor: 0x1099bb,
})
document.body.appendChild(app.view)

const sprite = PIXI.Sprite.from('./assets/circle.png')
sprite.x = 300
sprite.y = 300
sprite.anchor.set(0.5)
app.stage.addChild(sprite)

const socket = new WebSocket("ws://localhost:8080/ws")

time = 0

app.ticker.add((delta) => {
    sprite.rotation += 0.01 * delta

    sprite.x += Math.sin(app.ticker.lastTime / 1000) * 2 * delta

    time += 0.01

    const red = Math.abs(Math.sin(time + 0))
    const green = Math.abs(Math.sin(time + (Math.PI / 3)))
    const blue = Math.abs(Math.sin(time + (2 * Math.PI / 3)))

    app.renderer.backgroundColor = PIXI.utils.rgb2hex([red, green, blue])
});

function resize() {
    app.renderer.resize(window.innerWidth, window.innerHeight)
}

window.addEventListener('resize', resize)

socket.onopen = () => {
    console.log("Connected to WebSocket server")
    socket.send("Hello, Server!")
};

socket.onmessage = (event) => {
    console.log("Received message from server:", event.data)
};

socket.onclose = () => {
    console.log("WebSocket connection closed")
};

socket.onerror = (error) => {
    console.error("WebSocket error:", error)
};
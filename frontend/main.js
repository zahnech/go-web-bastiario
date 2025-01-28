
console.log("PixiJS loaded successfully!");

const app = new PIXI.Application({
    width: window.innerWidth,
    height: window.innerHeight,
    backgroundColor: 0x1099bb,
})
document.body.appendChild(app.view)

const playerSprites = {};

const socket = new WebSocket("ws://localhost:8080/ws")

time = 0

app.ticker.add((delta) => {
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
    const data = JSON.parse(event.data);
    if (data.type !== undefined) {
        switch (data.type) {
            case "pl_init":
                if (!playerSprites[data.id]) {
                    const sprite = PIXI.Sprite.from('./assets/circle.png')
                    sprite.x = data.x
                    sprite.y = data.y
                    sprite.anchor.set(0.5)
                    app.stage.addChild(sprite)

                    playerSprites[data.id] = sprite;

                    console.log(`Player ${data.id} added at position: (${data.x}, ${data.y})`)
                }
                break

            case "pl_del":
                if (playerSprites[data.id]) {
                    app.stage.removeChild(playerSprites[data.id])
                    delete playerSprites[data.id]
                    console.log(`Player ${data.id} removed from the game`)
                }
                break
            
            case "pl_chng_loc":
                if (playerSprites[data.id]) {
                    // add later
                }
                break
        }
    } else {
        console.log("Received data:", data);
    }
};

socket.onclose = () => {
    console.log("WebSocket connection closed")
};

socket.onerror = (error) => {
    console.error("WebSocket error:", error)
};
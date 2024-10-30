import { MessageTypes } from "./myEvents.js"
import { GameManager } from "./gameManager.js"
import { GraphicsManger } from "./graphicsManager.js"
import { Tile, TileSide } from "./tile.js"
import { WS_driver } from "./websockets.js"


var driver:WS_driver = new WS_driver()

var canvasElem = document.getElementById('canvas') as HTMLCanvasElement
canvasElem.width = window.innerWidth
canvasElem.height = window.innerHeight - 50

var manager = new GameManager(null, driver)
var graphManager = new GraphicsManger(canvasElem.getContext("2d"), manager.tiles)


var butt = document.getElementById('make') as HTMLButtonElement
butt.onclick = (ev) =>{
    manager.makeRoom()
}
var joinButt = document.getElementById("join") as HTMLButtonElement
var inputTxt = document.getElementById("roomId") as HTMLInputElement
joinButt.onclick = (ev) =>{
    manager.roomId = inputTxt.value
    manager.joinRoom()
}

var startGame = document.getElementById("startGame") as HTMLButtonElement
startGame.onclick = (ev) => {
    console.log("CLICKED START")
    manager.startGame()
}

Tile.origin_x = Math.floor((canvasElem.width / 2) / Tile.width) * Tile.width
Tile.origin_y = Math.floor(canvasElem.height / (2 * Tile.height)) * Tile.height

graphManager.redraw()

driver.attach(MessageTypes.sendTile, manager)
driver.attach(MessageTypes.sendTile, graphManager)

canvasElem.addEventListener("wheel", (ev:WheelEvent)=>{
    Tile.resize(Math.floor((ev.deltaY)*( - 1) / 50))
    Tile.origin_x = Math.floor((canvasElem.width / 2) / Tile.width) * Tile.width
    Tile.origin_y = Math.floor(canvasElem.height / (2 * Tile.height)) * Tile.height
    graphManager.redraw()
})

canvasElem.addEventListener("mousedown", (ev:MouseEvent)=>{
    let x = Math.floor(ev.clientX / Tile.width) * Tile.width
    let y = Math.floor(ev.clientY / Tile.height) * Tile.height
    manager.addTile(x, y, [TileSide.road, TileSide.city, TileSide.grass, TileSide.city])
    graphManager.redraw()
})



//manager.makeRoom()
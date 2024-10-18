import { GameManager } from "./gameManager.js"
import { WS_driver } from "./websockets.js"


var driver:WS_driver = new WS_driver()

var canvasElem = document.getElementById('canvas') as HTMLCanvasElement
canvasElem.width = window.innerWidth
canvasElem.height = window.innerHeight - 50

var ctx = canvasElem.getContext('2d')
ctx?.fillRect(0,0,window.innerWidth, window.innerHeight)

var manager = new GameManager(null, driver)

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

//manager.makeRoom()
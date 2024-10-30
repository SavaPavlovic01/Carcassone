import { GameManager } from "./gameManager.js";
import { GraphicsManger } from "./graphicsManager.js";
import { WS_driver } from "./websockets.js";

export class Game{

    gameManager:GameManager
    graphicsMangaer:GraphicsManger
    driver:WS_driver
    canvasElem:HTMLCanvasElement

    constructor(){
        this.driver = new WS_driver();
        this.canvasElem = document.getElementById('canvas') as HTMLCanvasElement
        this.canvasElem.width = window.innerWidth
        this.canvasElem.height = window.innerHeight - 50
        
        this.gameManager = new GameManager(null, this.driver)
        this.graphicsMangaer = new GraphicsManger(this.canvasElem.getContext("2d"), this.gameManager.tiles)

    }
}
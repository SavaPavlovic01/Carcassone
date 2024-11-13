import { GameManager } from "./gameManager.js";
import { GraphicsManger } from "./graphicsManager.js";
import { Meeple } from "./meeple.js";
import { MessageTypes } from "./myEvents.js";
import { Tile, TileSide } from "./tile.js";
import { WS_driver } from "./websockets.js";

export class Game{

    gameManager:GameManager
    graphicsMangaer:GraphicsManger
    driver:WS_driver
    canvasElem:HTMLCanvasElement
    activeMeeple:[Meeple|null, number]
    activeMeeplePlaced:boolean = true

    

    constructor(){
        this.driver = new WS_driver();
        this.canvasElem = document.getElementById('canvas') as HTMLCanvasElement
        this.canvasElem.width = window.innerWidth
        this.canvasElem.height = window.innerHeight - 50
        this.activeMeeple = [null, -1]
        
        this.gameManager = new GameManager(null, this.driver)
        this.graphicsMangaer = new GraphicsManger(this.canvasElem.getContext("2d"), this.gameManager.tiles)

        this.tempButtons()

        Tile.origin_x = Math.floor((this.canvasElem.width / 2) / Tile.width) * Tile.width
        Tile.origin_y = Math.floor(this.canvasElem.height / (2 * Tile.height)) * Tile.height

        this.graphicsMangaer.redraw()

        this.driver.attach(MessageTypes.sendTile, this.gameManager)
        this.driver.attach(MessageTypes.sendTile, this.graphicsMangaer)
        this.driver.attach(MessageTypes.sendMeeple, this.graphicsMangaer)
        this.driver.attach(MessageTypes.removeMeeple, this.graphicsMangaer)
        this.driver.attach(MessageTypes.movedMeeple, this.graphicsMangaer)
        this.driver.attach(MessageTypes.tempTilePlaced, this.graphicsMangaer)
        this.driver.attach(MessageTypes.startGame, this.gameManager)
        this.driver.attach(MessageTypes.startGame, this.graphicsMangaer)
        this.driver.attach(MessageTypes.joinRoom, this.graphicsMangaer)

        this.addCanvasListeners()

    }

    private tempButtons(){
        let butt = document.getElementById('make') as HTMLButtonElement
        butt.onclick = (ev) =>{
            this.gameManager.makeRoom()
        }
        let joinButt = document.getElementById("join") as HTMLButtonElement
        let inputTxt = document.getElementById("roomId") as HTMLInputElement
        joinButt.onclick = (ev) =>{
            this.gameManager.roomId = inputTxt.value
            this.gameManager.joinRoom()
        }

        let startGame = document.getElementById("startGame") as HTMLButtonElement
        startGame.onclick = (ev) => {
            console.log("CLICKED START")
            this.gameManager.startGame()
        }

        let endTurn = document.getElementById("endTurn") as HTMLButtonElement
        endTurn.onclick = (ev) =>{
            if(!this.graphicsMangaer.tempTileValidPostion || !this.graphicsMangaer.tempTile){
                console.log("NE MOZ")
                return
            }
            if(!this.gameManager.myTurn) return;
            let index = [this.graphicsMangaer.tempTile.x, this.graphicsMangaer.tempTile.y]
            let id = GameManager.indexToStr(index)
            this.graphicsMangaer.tempTile.overlayColor = ""
            this.gameManager.tiles.set(id, this.graphicsMangaer.tempTile)
            console.log(this.graphicsMangaer.tempTile)
            const msg = {
                msgType:MessageTypes.sendTile,
                roomId:this.gameManager.roomId,
                playerId:this.gameManager.playerId,
                sides:this.graphicsMangaer.tempTile.sides,
                x:index[0],
                y:index[1]
            }

            this.driver.send_msg(JSON.stringify(msg))
        }

        let rotate = document.getElementById("rotate") as HTMLButtonElement
        rotate.onclick = (ev) =>{
            this.graphicsMangaer.drawnTile?.rotate()
            if(this.graphicsMangaer.tempTile && this.graphicsMangaer.drawnTile){
                this.graphicsMangaer.tempTile.sides = this.graphicsMangaer.drawnTile.sides
                let coords = this.graphicsMangaer.tempTile.indexToCoord()
                this.graphicsMangaer.tempTileValidPostion = this.gameManager.checkIfTileValid(coords[0], coords[1], this.graphicsMangaer.drawnTile.sides)
            }
            this.graphicsMangaer.redraw()
        }
    }

    private addCanvasListeners(){
        this.canvasElem.addEventListener("wheel", (ev:WheelEvent)=>{
            Tile.resize(Math.floor((ev.deltaY)*( - 1) / 50))
            Tile.origin_x = Math.floor((this.canvasElem.width / 2) / Tile.width) * Tile.width
            Tile.origin_y = Math.floor(this.canvasElem.height / (2 * Tile.height)) * Tile.height
            this.graphicsMangaer.redraw()
        })
        
        this.canvasElem.addEventListener("mousedown", (ev:MouseEvent)=>{
            if(!this.gameManager.myTurn) return;
            let meep = this.graphicsMangaer.clickedMeeple(ev.clientX, ev.clientY)
            if(meep[0] != null){
                this.activeMeeple = meep
                this.activeMeeplePlaced = true
                return
            }
            if(ev.clientX >= this.graphicsMangaer.meepleUiX && ev.clientY >= this.graphicsMangaer.meepleUiY){
                this.activeMeeple[0] = new Meeple(ev.clientX, ev.clientY)
                this.activeMeeple[1] = this.graphicsMangaer.meeples.length 
                this.graphicsMangaer.addMeeple(this.activeMeeple[0])
                this.activeMeeplePlaced = false
                this.graphicsMangaer.redraw()
                return
            }
            let x = Math.floor(ev.clientX / Tile.width) * Tile.width
            let y = Math.floor(ev.clientY / Tile.height) * Tile.height
            //this.gameManager.addTile(x, y, [TileSide.road, TileSide.city, TileSide.grass, TileSide.city])
            if(!this.graphicsMangaer.drawnTile) return
            this.graphicsMangaer.tempTileValidPostion = this.gameManager.checkIfTileValid(x, y, this.graphicsMangaer.drawnTile.sides)
            let index = Tile.coordToIndex(x, y)
            this.graphicsMangaer.tempTile = new Tile(index[0], index[1], this.graphicsMangaer.drawnTile.sides)
            const msg = {
                msgType: MessageTypes.tempTilePlaced,
                roomId:this.gameManager.roomId,
                playerId:this.gameManager.playerId,
                sides:this.graphicsMangaer.drawnTile.sides,
                x:index[0],
                y:index[1],
                isValid:this.graphicsMangaer.tempTileValidPostion
            }
            this.driver.send_msg(JSON.stringify(msg))
            this.graphicsMangaer.redraw()
        })

        this.canvasElem.addEventListener("mouseup", (ev:MouseEvent) =>{
            if(!this.gameManager.myTurn) return;
            if(!this.activeMeeple[0]) return
            let tile = this.graphicsMangaer.checkIfMeepleValidPosition(this.activeMeeple[0])
            if(!tile){
                this.graphicsMangaer.removeMeeple(this.activeMeeple[1])
                this.graphicsMangaer.redraw()
                
                if(this.activeMeeplePlaced){
                    // poruka za premestanje
                    const msg = {
                        msgType:MessageTypes.removeMeeple,
                        roomId:this.gameManager.roomId,
                        playerId:this.gameManager.playerId,
                        index:this.activeMeeple[1]
                    }
                    console.log("Sklonjen sa table vec postojeci")
                    this.driver.send_msg(JSON.stringify(msg))
                } 
            }else{
                
                if(this.activeMeeplePlaced){
                    // vec postojeci premesten
                    const msg = {
                        msgType:MessageTypes.movedMeeple,
                        roomId:this.gameManager.roomId,
                        playerId:this.gameManager.playerId,
                        index:this.activeMeeple[1],
                        x:this.activeMeeple[0].x,
                        y:this.activeMeeple[0].y
                    }
                    this.driver.send_msg(JSON.stringify(msg))
                } else {
                    // dodat nov
                    const msg = {
                        msgType: MessageTypes.sendMeeple,
                        roomId:this.gameManager.roomId,
                        playerId:this.gameManager.playerId,
                        x:this.activeMeeple[0].x,
                        y:this.activeMeeple[0].y,
                        color:"red",
                        isPriest:false
                    }
                    this.driver.send_msg(JSON.stringify(msg))
                }                
            }
            this.activeMeeple = [null, -1]
        })

        this.canvasElem.addEventListener("mousemove", (ev:MouseEvent) =>{
            if(!this.activeMeeple[0]) return
            if(!this.gameManager.myTurn) return
            this.activeMeeple[0].x = ev.clientX
            this.activeMeeple[0].y = ev.clientY
            this.graphicsMangaer.redraw()
        })
    }
}
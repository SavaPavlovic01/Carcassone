import { Meeple } from "./meeple.js"
import { Listener, MessageTypes } from "./myEvents.js"
import { Tile } from "./tile.js"

export class GraphicsManger implements Listener{

    ctx:CanvasRenderingContext2D | null
    backgroundColor:string = "black"
    readonly gameState:Map<string, Tile>

    meepleUiX:number = 0
    meepleUiY:number = 0
    meepleUiSize:number = 0

    meeples:Meeple[] = []

    tempTile:Tile | null = null
    tempTileValidPostion:boolean = false

    drawnTile:Tile | null = null

    playerCnt:number = 0
    gameStarted:boolean = false

    constructor(ctx:CanvasRenderingContext2D | null, gameState:Map<string, Tile>){
        this.ctx = ctx
        this.gameState = gameState
        this.meepleUiSize = 50
        if(!ctx) return
        this.meepleUiX = ctx.canvas.width - this.meepleUiSize
        this.meepleUiY = ctx.canvas.height - this.meepleUiSize
    }

    public clear(){
        if(!this.ctx) return
        let tempStyle = this.ctx.fillStyle
        this.ctx.fillStyle = this.backgroundColor
        this.ctx.fillRect(0, 0, this.ctx.canvas.width, this.ctx.canvas.height)
        this.ctx.fillStyle = tempStyle
    }

    public drawTiles(){
        if(!this.gameState) return
        for(let [_, tile] of this.gameState){
            tile.draw(this.ctx)
        }
        
        if(this.tempTile) {
            // check if tempTile in valid position
            if(this.tempTileValidPostion){
                this.tempTile.overlayColor = "rgb(0,255,0,0.5)"
            } else{
                this.tempTile.overlayColor = "rgb(255,0,0,0.5)"
            }
            this.tempTile.draw(this.ctx)
        }
    }
    // meepleLeft cemo da dobijemo iz websocket poruke
    public drawMeepleUI(meepleLeft:number){
        
        if(!this.ctx) return
        let tempFill = this.ctx?.fillStyle
        this.ctx.fillStyle = "rgba(255,0,0,0.5)"
        this.ctx?.fillRect(this.meepleUiX, this.meepleUiY, this.meepleUiSize, this.meepleUiSize)
        this.ctx?.strokeText(meepleLeft.toString(), this.meepleUiX + this.meepleUiSize / 2, this.meepleUiY + this.meepleUiSize /2, this.meepleUiSize / 2)
        this.ctx.fillStyle = tempFill
    }

    public drawDrawnTile(){
        if(!this.drawnTile || !this.ctx) return
        this.drawnTile.overlayColor = "rgb(255,255,255,0.5)"
        this.drawnTile.drawAsUI(this.ctx) 
    }

    public drawPlayerCnt(){
        if(!this.ctx) return
        let temp = this.ctx.fillStyle
        this.ctx.fillStyle = "white"
        this.ctx.fillText("Players in room:" + this.playerCnt.toString(), 0, 25, 100)
        this.ctx.fillStyle = temp
    }

    public redraw(){
        if(!this.gameState) return
        this.clear()
        this.drawMeepleUI(5)
        this.drawDrawnTile()
        if(!this.gameStarted)this.drawPlayerCnt()
        this.drawTiles()
        this.drawMeeple()
    }

    notify(eventType: number, msg: any) {
        if(eventType == MessageTypes.sendTile){
            this.tempTile = null
            this.tempTileValidPostion = false
            if(!this.gameState) return
            this.redraw()
        }

        if(eventType == MessageTypes.sendMeeple){
            let curMeeple = new Meeple(msg["x"], msg["y"])
            this.meeples.push(curMeeple)
            this.redraw()
        }

        if(eventType == MessageTypes.removeMeeple){
            this.removeMeeple(msg["index"])
            this.redraw()
        }

        if(eventType == MessageTypes.movedMeeple){
            this.meeples[msg["index"]].x = msg["x"]
            this.meeples[msg["index"]].y = msg["y"]
            this.redraw()
        }

        if(eventType == MessageTypes.tempTilePlaced){
            this.tempTile = new Tile(msg["tileX"], msg["tileY"], msg["tileSides"])
            this.tempTileValidPostion = msg["isValid"]
            this.redraw()
        }

        if(eventType == MessageTypes.startGame){
            this.drawnTile = new Tile(0, 0, msg["tileSides"])
            this.gameStarted = true
            this.redraw()
        }

        if(eventType == MessageTypes.joinRoom){
            this.playerCnt = msg.playerCnt
            this.redraw()
        }
    }

    public addMeeple(meeple:Meeple){
        this.meeples.push(meeple)
    }

    public removeMeeple(index:number){
        if(index == -1) return
        this.meeples.splice(index, 1)
    }

    public drawMeeple(){
        if(!this.ctx) return
        for(let meeple of this.meeples){
            meeple.draw(this.ctx)
        }
    }

    public checkIfMeepleValidPosition(meeple:Meeple):Tile|null{
        /*
        for(let [_, tile] of this.gameState){
            let realCoords = tile.indexToCoord()
            let checkX = meeple.x >= realCoords[0] && meeple.x <= realCoords[0] + Tile.width
            let checkY = meeple.y >= realCoords[1] && meeple.y <= realCoords[1] + Tile.height
            if(checkX && checkY){
                return tile
            }
        }*/
        if(!this.tempTile || !this.tempTileValidPostion) return null
        let realCoords = this.tempTile.indexToCoord()
        let checkX = meeple.x >= realCoords[0] && meeple.x <= realCoords[0] + Tile.width
        let checkY = meeple.y >= realCoords[1] && meeple.y <= realCoords[1] + Tile.height
        if(checkX && checkY){
            return this.tempTile
        }
        return null
    }

    public clickedMeeple(x:number, y:number) : [Meeple|null, number] {
        for(let [i, meeple] of this.meeples.entries()){
            let xOk = x >= meeple.x && x <= meeple.x + Meeple.meepleSize 
            let yOk = y >= meeple.y && y <= meeple.y + Meeple.meepleSize
            if(xOk && yOk){
                return [meeple, i]
            }
        }
        return [null, -1]
    }
}
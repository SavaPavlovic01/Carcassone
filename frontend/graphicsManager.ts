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

    public redraw(){
        if(!this.gameState) return
        this.clear()
        this.drawTiles()
        this.drawMeepleUI(5)
        this.drawMeeple()
    }

    notify(eventType: number, msg: any) {
        if(eventType == MessageTypes.sendTile){
            if(!this.gameState) return
            this.redraw()
        }
    }

    public addMeeple(meeple:Meeple){
        this.meeples.push(meeple)
    }

    public removeMeeple(index:number){
        this.meeples.splice(index, 1)
    }

    public drawMeeple(){
        if(!this.ctx) return
        for(let meeple of this.meeples){
            meeple.draw(this.ctx)
        }
    }

    public checkIfMeepleValidPosition(meeple:Meeple):Tile|null{
        for(let [_, tile] of this.gameState){
            let realCoords = tile.indexToCoord()
            let checkX = meeple.x >= realCoords[0] && meeple.x <= realCoords[0] + Tile.width
            let checkY = meeple.y >= realCoords[1] && meeple.y <= realCoords[1] + Tile.height
            if(checkX && checkY){
                return tile
            }
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
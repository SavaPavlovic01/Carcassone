import { Listener, MessageTypes } from "./myEvents.js"
import { Tile } from "./tile.js"

export class GraphicsManger implements Listener{

    ctx:CanvasRenderingContext2D | null
    backgroundColor:string = "black"
    readonly gameState:Map<string, Tile>

    constructor(ctx:CanvasRenderingContext2D | null, gameState:Map<string, Tile>){
        this.ctx = ctx
        this.gameState = gameState
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
        const size = 50
        if(!this.ctx) return
        let tempFill = this.ctx?.fillStyle
        this.ctx.fillStyle = "rgba(255,0,0,0.5)"
        this.ctx?.fillRect(this.ctx.canvas.width - size, this.ctx.canvas.height - size, size, size)
        this.ctx?.strokeText(meepleLeft.toString(), this.ctx.canvas.width - size / 2, this.ctx.canvas.height - size /2, size/2)
        this.ctx.fillStyle = tempFill
    }

    public redraw(){
        if(!this.gameState) return
        this.clear()
        this.drawTiles()
        this.drawMeepleUI(5)
    }

    notify(eventType: number, msg: any) {
        if(eventType == MessageTypes.sendTile){
            if(!this.gameState) return
            this.redraw()
        }
    }
}
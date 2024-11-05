export enum TileSide{
    grass,
    road,
    city,
    citty_connected
}

export enum Orientation{
    left,
    top,
    right,
    bottom
}

export class Tile{
    static width:number = 100
    static height:number = 100

    static origin_x:number = 200
    static origin_y:number = 200

    sides:TileSide[]

    x:number;
    y:number;

    overlayColor:string = ""

    constructor(_x:number, _y:number, _sides:TileSide[]){
        this.x = _x
        this.y = _y
        this.sides = _sides
    }

    // TODO: Maybe move to graphicsManager
    public draw(ctx:CanvasRenderingContext2D | null){
        if(ctx !== null){
            ctx.strokeStyle = "white"
            ctx.strokeRect(Tile.origin_x + this.x*Tile.width, Tile.origin_y - this.y* Tile.height, Tile.width, Tile.height)
            let number_coord = this.indexToCoord()
            ctx.strokeText(String(this.sides[0]), number_coord[0], number_coord[1] + Tile.height / 2)
            ctx.strokeText(String(this.sides[1]), number_coord[0] + Tile.width / 2, number_coord[1])
            ctx.strokeText(String(this.sides[2]), number_coord[0] + Tile.width , number_coord[1] + Tile.height / 2)
            ctx.strokeText(String(this.sides[3]), number_coord[0] + Tile.width / 2, number_coord[1] + Tile.height)

            if(this.overlayColor != ""){
                let tempStyle = ctx.fillStyle
                ctx.fillStyle = this.overlayColor
                ctx.fillRect(number_coord[0], number_coord[1], Tile.width, Tile.height)
                ctx.fillStyle = tempStyle
            }
        }
        
    }

    public drawAsUI(ctx:CanvasRenderingContext2D){
        ctx.strokeStyle = "white"
        ctx.strokeRect(0, ctx.canvas.height - 50, 50, 50)
        ctx.strokeText(String(this.sides[0]), 0, ctx.canvas.height - 25)
        ctx.strokeText(String(this.sides[1]), 0 + 25, ctx.canvas.height - 50)
        ctx.strokeText(String(this.sides[2]), 0 + 50 , ctx.canvas.height - 25)
        ctx.strokeText(String(this.sides[3]), 0 + 25, ctx.canvas.height )

        if(this.overlayColor != ""){
            let tempStyle = ctx.fillStyle
            ctx.fillStyle = this.overlayColor
            ctx.fillRect(0, ctx.canvas.height, 50, 50)
            ctx.fillStyle = tempStyle
        }
    }

    static resize(amount:number){
        this.width += amount
        this.height += amount
    }

    static coordToIndex(x:number, y:number){
        return [(x - Tile.origin_x) / Tile.width, (Tile.origin_y - y) / Tile.height]
    }

    indexToCoord(){
        return [(Tile.origin_x + this.x * Tile.width), (Tile.origin_y - this.y * Tile.height)]
    }

    public rotate(){
        this.sides = [this.sides[3], ...this.sides.slice(0, 3)]
    }
}
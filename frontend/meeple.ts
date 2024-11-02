export class Meeple{
    static meepleSize = 50
    x:number
    y:number

    constructor(_x:number, _y:number){
        this.x = _x
        this.y = _y
    }

    public draw(ctx:CanvasRenderingContext2D){
        let tempFill = ctx.fillStyle
        ctx.fillStyle = "red"
        ctx.fillRect(this.x, this.y, Meeple.meepleSize, Meeple.meepleSize)
        ctx.fillStyle = tempFill
    }
}
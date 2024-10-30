import { Constants } from "./constants.js";
import { Listener, MessageTypes } from "./myEvents.js";
import { Orientation, Tile, TileSide } from "./tile.js";
import { WS_driver } from "./websockets.js";

export class GameManager implements Listener{
    roomId:string|null;
    wsDriver:WS_driver;
    playerId:string
    tiles:Map<string, Tile>
    meepleLeft:number = 5

    constructor(_roomId:string|null, _wsDriver:WS_driver){
        this.roomId = _roomId
        this.wsDriver = _wsDriver
        this.playerId = crypto.randomUUID()
        this.tiles = new Map()
        this.tiles.set("0|0", new Tile(0, 0, [TileSide.road, TileSide.city, TileSide.road, TileSide.grass]))
    }

    public notify(eventType:number, msg: any) {
        if(eventType == MessageTypes.sendTile){
            let x = msg["tileX"]
            let y = msg["tileY"]
            this.tiles.set(GameManager.indexToStr([x, y]), new Tile(x,y, msg["tileSides"]))
        }
    }

    public joinRoom(){
        if(!this.wsDriver.is_open){
            return false
        }

        const msg = {
            msgType:MessageTypes.joinRoom,
            gameVersion:Constants.gameVersion,
            roomId:this.roomId,
            playerId:this.playerId
        }
        
        this.wsDriver.send_msg(JSON.stringify(msg))
        return true
    }

    public makeRoom(){
        if(!this.wsDriver.is_open){
            return false
        }

        const msg = {
            msgType:MessageTypes.makeRoom,
            gameVersion:Constants.gameVersion,
            roomId:this.playerId,
            playerId:this.playerId
        }
        this.roomId = this.playerId
        console.log(JSON.stringify(msg))
        this.wsDriver.send_msg(JSON.stringify(msg))
        console.log("POSLATO")
        return true
    }

    public startGame(){
        const msg = {
            msgType : MessageTypes.startGame,
            gameVersion:Constants.gameVersion,
            roomId:this.roomId,
            playerId:this.playerId
        }

        this.wsDriver.send_msg(JSON.stringify(msg))
        console.log("sent start game msg")
        return true
    }

    private static indexToStr(index:number[]){
        return `${index[0]}|${index[1]}`
    }

    public addTile(x:number, y:number, sides:TileSide[]){
        let index = Tile.coordToIndex(x, y)
        let hasNeighbor = false

        if(this.tiles.get(GameManager.indexToStr(index)) != undefined){
            console.log("OCCUPIED")
            return
        }
        

        let orientations = [Orientation.left, Orientation.top, Orientation.right, Orientation.bottom]

        
        for(let curSide of orientations){
            let curTile = this.getNeighbor(index[0], index[1], curSide)
            console.log(curTile)
            if(curTile != undefined){
                hasNeighbor = true
                if(!this.sameSideType(curTile.sides[this.getOppositeSide(curSide)], sides[curSide])) return
            }
        }

        if(hasNeighbor){
            let addedTile = new Tile(-index[0], index[1], sides)
            this.tiles.set(GameManager.indexToStr([index[0], index[1]]), addedTile)
            this.sendTile(addedTile)
        }
    }

    private getNeighbor(x:number, y:number, pos:Orientation){
        let offsets = [[1, 0], [0, -1], [-1, 0], [0, 1]]
        console.log(pos)
        return this.tiles.get(GameManager.indexToStr([x + offsets[pos][0], y + offsets[pos][1]]))
    }

    private getOppositeSide(side:Orientation){
        return (side + 2) % 4 
    }

    private sameSideType(x:TileSide, y:TileSide){
        if(x == y) return true
        if((x == TileSide.city && y == TileSide.citty_connected) || (x == TileSide.citty_connected && y == TileSide.city)) return true
        return false
    }

    public sendTile(tile:Tile){
        const msg = {
            msgType:MessageTypes.sendTile,
            version:Constants.gameVersion,
            x:tile.x,
            y:tile.y,
            sides:tile.sides,
            roomId:this.roomId,
            playerId:this.playerId
        }
        this.wsDriver.send_msg(JSON.stringify(msg))
    }

}
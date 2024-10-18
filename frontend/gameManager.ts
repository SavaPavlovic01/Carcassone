import { Constants } from "./constants.js";
import { MessageTypes } from "./events.js";
import { WS_driver } from "./websockets.js";

export class GameManager{
    roomId:string|null;
    wsDriver:WS_driver;
    playerId:string

    constructor(_roomId:string|null, _wsDriver:WS_driver){
        this.roomId = _roomId
        this.wsDriver = _wsDriver
        this.playerId = crypto.randomUUID()
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

}
import { Listener } from "./myEvents.js";

export class WS_driver{
    socket:WebSocket;
    is_open:boolean = false

    subscribers:Map<number, Listener[]>

    constructor(){
        this.socket = new WebSocket("ws://localhost:8080/ws")
        this.socket.onopen = (ev:Event) =>{
            console.log("OPEN")
            this.is_open = true
        }
        this.socket.onmessage = (ev:MessageEvent) =>{
            let data = JSON.parse(ev.data)
            this.subscribers.get(data["msgType"])?.forEach((listener:Listener) =>{
                listener.notify(data["msgType"], data);
            })
            //alert(JSON.stringify(data))
            console.log(data)
        }
        this.subscribers = new Map<number, Listener[]>()
    }

    public send_msg(msg:string) {
        console.log("CLICK")
        if(this.is_open){
            this.socket.send(msg)
        }   
    }

    public attach(msgType:number, listener:Listener){
        let listeners = this.subscribers.get(msgType)
        if(!listeners){
            this.subscribers.set(msgType, [listener]);
        } else listeners?.push(listener);
    }

}
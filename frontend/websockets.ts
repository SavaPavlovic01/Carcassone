export class WS_driver{
    socket:WebSocket;
    is_open:boolean = false

    constructor(){
        this.socket = new WebSocket("ws://localhost:8080")
        this.socket.onopen = (ev:Event) =>{
            console.log("OPEN")
            this.is_open = true
        }
        this.socket.onmessage = (ev:MessageEvent) =>{
            alert(ev.data)
        }
    }

    public send_msg(msg:string) {
        console.log("CLICK")
        if(this.is_open){
            this.socket.send(msg)
        }   
    }

}
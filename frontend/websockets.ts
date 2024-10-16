window.onload = (event:Event) => {
    let socket = new WebSocket("ws://localhost:8080/")
    socket.onopen = (ev:Event) => {
        socket.onmessage = (msg_ev:MessageEvent) => {
            alert(msg_ev.data)
        }
        socket.send("Hello world")
    }
}
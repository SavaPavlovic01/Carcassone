export enum MessageTypes{
    joinRoom,
    makeRoom,
    startGame,
    sendTile
}

export interface Listener{
    notify(eventType:number, msg:any):any;
}
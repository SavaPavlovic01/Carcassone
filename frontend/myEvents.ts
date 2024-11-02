export enum MessageTypes{
    joinRoom,
    makeRoom,
    startGame,
    sendTile,
    sendMeeple,
    removeMeeple,
    movedMeeple
}

export interface Listener{
    notify(eventType:number, msg:any):any;
}
export enum MessageTypes{
    joinRoom,
    makeRoom,
    startGame,
    sendTile,
    sendMeeple,
    removeMeeple,
    movedMeeple,
    tempTilePlaced,
    pullTile
}

export interface Listener{
    notify(eventType:number, msg:any):any;
}
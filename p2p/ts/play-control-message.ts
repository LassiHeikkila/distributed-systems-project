enum ControlMessageType {
    Play = "play",
    Pause = "pause",
    SetTime = "setTime",
}

class ControlMessage {
    messageType: ControlMessageType;
    message: ControlMessagePlay | ControlMessagePause | ControlMessageSetTime

    constructor(
        public obj: {
            messageType: ControlMessageType
            message: ControlMessagePlay | ControlMessagePause | ControlMessageSetTime
        }
    ) {
        
        this.messageType = obj.messageType;
        this.message = obj.message;
    }
};

type ControlMessagePlay = {};
type ControlMessagePause = {};
type ControlMessageSetCurrentTime = {
    currentTimeSeconds: number;
};

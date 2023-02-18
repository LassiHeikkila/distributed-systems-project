class ChatMessage {
    content: string;
    from: string;
    to: string;
    
    constructor(
        public contentArg: string
    ) {
        this.content = contentArg;
    }
};
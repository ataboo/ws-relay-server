
class Client {
    wsClient;
    textDecoder = new TextDecoder();
    textEncoder = new TextEncoder();
    playerId = 0;
    messages = [];

    init() {
        this.connectButton.onclick = () => this.connect();
        this.sendButton.onclick = () => this.sendMsg();
    }

    connect() {
        console.log("connecting");

        this.setConnected(true);

        const wsPath = `wss://localhost:3000/ws`;

        this.wsClient = new WebSocket(wsPath);

        this.wsClient.addEventListener('close', (evt) => {console.log("ws closed", evt)});
        this.wsClient.addEventListener('open', (evt) => {console.log("ws connected", evt)});
        this.wsClient.addEventListener('message', async (evt) => {
            const bytes = await evt.data.bytes();
            const msg = this.decodeMsg(bytes);

            this.handleMsg(msg);
        });
        this.wsClient.addEventListener('error', (evt) => {
            console.log("ws error", evt);
            if (this.wsClient.connected) {
                this.setConnected(false);
            }
        });
    }

    sendJoin() {
        console.log('joining');

        const userName = document.getElementById('username-input').value;

        const payload = JSON.stringify({
            name: userName,
            room_code: 'ABCDEF',
        });

        const msg = {
            version: 1,
            code: 2,
            sender: this.playerId,
            payloadBytes: this.textEncoder.encode(payload),
        }

        const msgArrBuffer = this.encodeMsg(msg);

        this.wsClient.send(msgArrBuffer);
    }

    sendMsg() {
        const messageStr = this.messageInput.value;
        if(!messageStr) {
            return;
        }

        const msg = {
            version: 1,
            code: 3,
            sender: this.playerId,
            payloadBytes: this.textEncoder.encode(messageStr),
        };

        const msgArrBuffer = this.encodeMsg(msg);

        this.wsClient.send(msgArrBuffer);
    }

    encodeMsg(msg) {
        const length = 10 + (msg.payloadBytes?.length ?? 0);

        const bytes = [];
        bytes.push(...this.uintToBytes(length, 4));
        bytes.push(...this.uintToBytes(msg.version, 2));
        bytes.push(...this.uintToBytes(msg.code, 2));
        bytes.push(...this.uintToBytes(this.playerId, 2));
        bytes.push(...msg.payloadBytes)

        const typedArr = new Uint8Array(bytes);

        return typedArr.buffer;
    }

    decodeMsg(bytes) {
        const length = this.bytesToUint(bytes.slice(0, 4), 4);
        if (bytes.length !== length) {
            throw new Error("Unexpected message length");
        }

        const version = this.bytesToUint(bytes.slice(4, 6), 2);
        const code = this.bytesToUint(bytes.slice(6, 8), 2);
        const sender = this.bytesToUint(bytes.slice(8, 10), 2);
        const payloadBytes = bytes.slice(10);

        let payloadObj = undefined;
        if(length > 10) {
            const payloadStr = this.textDecoder.decode(payloadBytes);
            payloadObj = JSON.parse(payloadStr)
        }

        const msg = {
            version,
            code,
            sender,
            payloadBytes,
            payloadObj,
        }

        return msg;
    }

    pushMsg(msg) {
        this.messages.push(msg);
        if (this.messages.length > 20) {
            this.messages = this.messages.slice(1);
        }

        this.messageBox.innerHTML = this.messages.map(m => `<div>${m}</div>`).join(' ');
    }

    handleMsg(msg) {
        console.dir(msg);

        switch(msg.code) {
            case 1:
                if(msg.payloadObj.user_id) {
                    this.playerId = msg.payloadObj.user_id;
                    console.log("player: ", msg.payloadObj.user_id);

                    this.sendJoin();
                }
            break;
        }

        this.pushMsg(JSON.stringify(msg.payloadObj,undefined, '\t'));
    }

    uintToBytes(num, byteCount) {
        const out = [];
        
        for(let i=0; i<byteCount; i++) {
            out.push(num & (0xFF<<i*8))
        }

        return out;
    }

    bytesToUint(byteSlice, length) {
        let out = 0;
        
        for(let i=0; i<length; i++) {
            out |= byteSlice[i] << i*8
        }

        return out
    }

    setConnected(connected) {
        if(connected) {
            this.connectButton.setAttribute('disabled', 'disabled');
            this.userNameInput.setAttribute('disabled', 'disabled');
        } else {
            this.connectButton.removeAttribute('disabled');
            this.userNameInput.removeAttribute('disabled');
        }
    }

    get connectButton() {
        return document.getElementById('connect-btn');
    }

    get userNameInput() {
        return document.getElementById('username-input');
    }

    get messageBox() {
        return document.getElementById('message-box');
    }

    get messageInput() {
        return document.getElementById('message-input');
    }

    get sendButton() {
        return document.getElementById('send-btn');
    }
}

function initClient() {
    const client = new Client();
    client.init();
    // document.getElementById('send-btn').onclick = () => client.sendJoin();
}
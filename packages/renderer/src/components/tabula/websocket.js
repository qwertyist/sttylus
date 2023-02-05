import { store } from "../../store"
import  EventBus  from "../../eventbus"
const mt = {
    OK: 0,
    CreateSession: 1,
    JoinSession: 2,
    LeaveSession: 3,
    NoSession: 404,
    Info: 4,
    SessionData: 5,
    Identify: 6,
    TXDelta: 20,
    RXDelta: 21,
    TXClear: 22,
    RXClear: 23,
    TXAbb: 24,
    RXAbb: 25,
    TXManuscript: 26,
    RXManuscript: 27,
    ReadySignal:  28,
    RetrieveDoc: 30,
    Loss: 500,
    Ping: 200,
    Pong: 300
}
function waitForConnection(socket, callback) {
    setTimeout(
        function () {
            if (socket.readyState === 1) {
                console.log("Connection successful")
                this.status = "connected"
                if (callback != null) {
                    callback();
                }
            } else {
                console.log("Wait for connection...")
                waitForConnection(socket, callback)
            }
        }, 5
    )
}

export default class wsConnection {
    constructor(quill, endpoint) {
        this.status = "disconnected"
        this.endpoint = endpoint
        this.mt = mt
        this.quill = quill
        this.quill.index = 0
        this.quill.version = 0
        self.websocket = null
        this.connect(endpoint)
        this.status = "pending"
    }
    connect(endpoint) {
        try {
            self.websocket = new WebSocket(endpoint);
        } catch (err) {
            console.log("Couldn't connect:", err)
            EventBus.$emit("localConnection", false)
            EventBus.$emit("websocketFailed")
            return
        }
        self.websocket.onopen = (e) => this.onOpen(e)
        self.websocket.onerror = (e) => this.onError(e)
        self.websocket.onclose = (e) => this.onClose(e)
        self.websocket.onmessage = (e) => this.onMessage(e)
    }
    ping() {
        if (!self.websocket) { 
            this.status = "disconnected"
            console.log("No websocket, dont ping")
            return;
        }
        if (self.websocket.readyState !== 1) {
            this.status = "pending"
            console.log("ReadyState not true")
            return;
        }
        let pingMessage = JSON.stringify(
            {
                type: this.mt.Ping
            }
        )
        this.status = "connected"
        self.websocket.send(pingMessage)
        setTimeout(() => { this.ping() }, 5000);
    }

    reconnect(tries) {
        EventBus.$emit("websocketReconnecting", tries)
        this.connect(this.endpoint)
        this.status = "reconnecting"
        self.websocket.onerror = (e) => { 
            console.log("hantera error här inne")
            tries++
            setTimeout(() => { this.reconnect(tries)}, 500+(tries*250))
        }
        self.websocket.onclose = (e) => { console.log("hantera close här inne")}
        self.websocket.onopen = (e) => { 
        this.onOpen(e)
        self.websocket.onerror = (e) => this.onError(e)
        self.websocket.onclose = (e) => this.onClose(e)
        self.websocket.onmessage = (e) => this.onMessage(e)
        }
        self.websocket.onmessage = (e) => { console.log("hantera message här inne") }

    }

    close() {
        console.log("close")
        self.websocket.close()
        self.websocket = null;
    }
    onOpen(e) {
        EventBus.$emit("websocketConnected")
        self.websocket.send("interpreter")
        this.joinsession()
        this.ping();
        this.status = "connected"
        store.commit("setLocalSession", {connected: true})
    }

    onClose(e) {
        store.commit("setLocalSession", {connected: false})
        console.log("onClose connection status:", this.status)
        if (self.websocket) {
            if(this.status == "disconnected" || this.status == "pending") {
                EventBus.$emit("localConnection", false)
                setTimeout(() => {
                    EventBus.$emit("websocketFailed")
                }, 250)
                return
            }
        }
        
        if(this.status == "reconnecting") {
            return
        }

        if(!e.wasClean) {
            EventBus.$emit("websocketDropped") 
            setTimeout(() => { 
                this.reconnect(1)
            }, 500)
            return
        }
        EventBus.$emit("websocketClosed") 
    }

    onError(e) {
        console.log("onError connection status:", this.status)
        if (this.status != "reconnecting") {
            if (self.websocket.readyState == 1) {
                console.log("ws normal error", e.type)
            }
        }
        //EventBus.$emit("websocketError")
    }

    onMessage(e) {
        //console.log(e.data)
        if (e.data) {
            let rx = JSON.parse(e.data)
            switch (rx.type) {
                case this.mt.Pong:
                    break
                case this.mt.CreateSession:
                    console.log("Server: Create Session")
                    this.createsession()
                    break
                case this.mt.JoinSession:
                    console.log("JoinSession message", rx)
                    EventBus.$emit("clientConnected", rx)
                    break
                case this.mt.LeaveSession:
                    console.log("LeaveSession message", rx)
                    EventBus.$emit("clientDisconnected", rx)
                    break
                case this.mt.RXDelta:
                    //console.log("RXDelta (version: ", rx.body.version, "):", rx.body.delta, rx.body.index)
                    //console.log("local version:", this.quill.version)
                    if (rx.body.version > this.quill.version) {
                        this.quill.updateContents(rx.body.delta, "collab")
                        if(!store.getters.getModalOpen) {
                            this.quill.setSelection(rx.body.index, 0, "collab")
                        } 
                        this.quill.version = rx.body.version
                    }
                    break
                case this.mt.RXClear:
                    //console.log("RXClear should clear and reset version");
                    this.quill.version = 0
                    this.quill.setText("")
                    break
                case this.mt.RetrieveDoc:
                    //console.log("RetrieveDoc:", rx.body.version, rx.body.delta)
                    this.quill.version = rx.body.version
                    this.quill.setContents(rx.body.delta, "collab")
                    this.quill.setSelection(rx.body.index);
                    break
                case this.mt.RXAbb:
                    EventBus.$emit("sharedAbb", rx.abb)         
                    break
              case this.mt.ReadySignal:
                    EventBus.$emit("recvReadySignal")
                    break
            }
        }
        //this.quill.updateContents()
    }
    send(data) {
        console.log(data)
        self.websocket.send(data)
    }
    sendClear() {
        let clearMessage = JSON.stringify(
            { type: this.mt.TXClear }
        )
        self.websocket.send(clearMessage)
    }
    sendDelta(data) {
        //console.log("sending version:", this.quill.version)
        this.quill.version++
        let deltaMessage = JSON.stringify(
            {
                type: this.mt.TXDelta, body: {
                    delta: data, version: this.quill.version
                }
            }
        )
        self.websocket.send(deltaMessage)
    }
    sendSharedAbb(abb) {
        let sharedAbbMessage = JSON.stringify({
            type: this.mt.TXAbb,
            abb: abb
        })
        console.log(sharedAbbMessage)
        self.websocket.send(sharedAbbMessage)
    }
    sendSessionData(data) {
        let sessionDataMessage = JSON.stringify({
            type: this.mt.SessionData,
            zoom: data
        })
        self.websocket.send(sessionDataMessage)
    }   
    sendReadySignal() {
      let readySignalMessage = JSON.stringify({
            type: this.mt.ReadySignal,
      })
      self.websocket.send(readySignalMessage)
    }

    createsession() {
        console.log("version:", this.quill.version)
        let createMessage = JSON.stringify({ type: this.mt.CreateSession, body: { version: this.quill.version, delta: this.quill.getContents() } })
        waitForConnection(self.websocket, function () { self.websocket.send(createMessage) })
    }
    joinsession() {
        let JoinMessage = JSON.stringify({ type: this.mt.JoinSession, msg:"interpreter"})
        console.log("join:", JoinMessage)
        waitForConnection(self.websocket, function () { self.websocket.send(JoinMessage) })

    }
    leavesession() {
        self.websocket.send(JSON.stringify({ type: this.mt.LeaveSession }))
    }
}

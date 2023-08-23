import EventBus from '../../eventbus'
const mt = {
  OK: 0,
  CreateSession: 1,
  JoinSession: 2,
  LeaveSession: 3,
  NotAuthorized: 401,
  NoSession: 404,
  Info: 4,
  TXDelta: 20,
  RXDelta: 21,
  TXClear: 22,
  RXClear: 23,
  TXAbb: 24,
  RXAbb: 25,
  TXManuscript: 26,
  RXManuscript: 27,
  RetrieveDoc: 30,
  TXChat: 40,
  RXChat: 41,
  Loss: 500,
  Ping: 200,
  Pong: 300,
}
function waitForConnection(socket, callback) {
  setTimeout(function () {
    if (socket.readyState === 1) {
      console.log('Connection successful')
      if (callback != null) {
        callback()
      }
    } else {
      console.log('Wait for connection...')
      waitForConnection(socket, callback)
    }
  }, 5)
}

export default class wsConnection {
  constructor(quill, endpoint, password) {
    this.status = 'disconnected'
    this.endpoint = endpoint
    this.mt = mt
    this.quill = quill
    this.quill.index = 0
    this.quill.version = 0
    this.buffer = []
    this.n = 0
    this.ready = false
    this.password = password
    self.websocket = null
    this.id = ''
    this.connect(endpoint, password)
  }

  connect(endpoint, password) {
    try {
      self.websocket = new WebSocket(endpoint)
    } catch (err) {
      console.error('tried endpoint', endpoint)
      console.error("Couldn't connect:", err)
      EventBus.$emit('localConnection', false)
      EventBus.$emit('websocketFailed')
      return
    }
    this.password = password
    self.websocket.onopen = (e) => this.onOpen(e)
    self.websocket.onerror = (e) => this.onError(e)
    self.websocket.onclose = (e) => this.onClose(e)
    self.websocket.onmessage = (e) => this.onMessage(e)
  }

  ping() {
    if (!self.websocket) {
      this.status = 'disconnected'
      return
    }
    if (self.websocket.readyState !== 1) {
      this.status = 'pending'
      return
    }
    let pingMessage = JSON.stringify({
      type: this.mt.Ping,
    })

    this.status = 'connected'
    self.websocket.send(pingMessage)
    setTimeout(() => {
      this.ping()
    }, 5000)
  }

  reconnect(tries) {
    EventBus.$emit('websocketReconnecting', tries)
    this.connect(this.endpoint, this.password)
    this.status = 'reconnecting'
    self.websocket.onerror = (e) => {
      tries++
      setTimeout(() => {
        this.reconnect(tries)
      }, 500 + tries * 250)
    }
    self.websocket.onclose = () => {}
    self.websocket.onmessage = (e) => {
    }
    self.websocket.onopen = (e) => {
      this.onOpen(e)
      self.websocket.onerror = (e) => this.onError(e)
      self.websocket.onclose = (e) => this.onClose(e)
      self.websocket.onmessage = (e) => this.onMessage(e)
    }
  }
  close() {
    self.websocket.close()
    self.websocket = null
  }
  onOpen(e) {
    this.ready = false
    // TODO: Add user specific message
    self.websocket.send('user')
    this.joinsession(this.password)
    this.ping()
    this.status = 'connected'
  }

  onClose(e) {
    if (self.websocket) {
      if (this.status == 'disconnected' || this.status == 'pending') {
        setTimeout(() => {
          EventBus.$emit('websocketFailed', e)
        }, 250)
        return
      }
    }

    if (this.status == 'reconnecting') {
      return
    }
    if (!e.wasClean) {
      EventBus.$emit('websocketDropped')
      setTimeout(() => {
        this.reconnect(1)
      }, 500)
      return
    }
    EventBus.$emit('websocketClosed')
  }

  onError(e) {
    if (this.status != 'reconnecting') {
      if (self.websocket.readyState == 1) {
        console.log('ws normal error', e.type)
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
          this.createsession()
          this.ready = true
          break
        case this.mt.JoinSession:
          if (!rx.id) {
            this.id = rx.id
          }
          break
        case this.mt.LeaveSession:
          break
        case this.mt.NotAuthorized:
          EventBus.$emit('notAuthorized', this.password)
          break
        case this.mt.RXDelta:
          if (rx.msg == 'starting') {
            this.quill.setContents(rx.body.delta, 'collab')
            this.quill.setSelection(rx.body.index, 0, 'collab')
            this.quill.version = rx.body.version
            this.ready = true
            return
          }
          if (this.ready) {
            if (rx.body.version > this.quill.version) {
              this.quill.updateContents(rx.body.delta, 'collab')
              this.quill.setSelection(rx.body.index, 0, 'collab')
              this.quill.version = rx.body.version
            } else {
            }
          } else {
            this.buffer.push(rx.body)
          }
          //console.log("RXDelta (version: ", rx.body.version, "):", rx.body.delta, rx.body.index)
          //console.log("local version:", this.quill.version)
          break
        case this.mt.RXClear:
          //console.log("RXClear should clear and reset version");
          this.quill.version = 0
          this.quill.setText('')
          break
        case this.mt.RetrieveDoc:
          switch (rx.msg) {
            case 'waiting':
              EventBus.$emit('joinedEmptySession')
              return
            case 'started':
              EventBus.$emit('joinedRunningSession')
              break
            default:
              break
          }

          this.quill.version = rx.body.version
          this.quill.setContents(rx.body.delta, 'collab')
          this.quill.setSelection(rx.body.index)
          this.ready = true
          break
        case this.mt.RXChat:
          EventBus.$emit('RXChat', rx)
          break
      }
    }
    //this.quill.updateContents()
  }
  send(data) {
    self.websocket.send(data)
  }
  sendClear() {
    let clearMessage = JSON.stringify({ type: this.mt.TXClear })
    self.websocket.send(clearMessage)
  }
  sendDelta(data) {
    //console.log("sending version:", this.quill.version)
    this.quill.version++
    let deltaMessage = JSON.stringify({
      type: this.mt.TXDelta,
      body: {
        delta: data,
        version: this.quill.version,
      },
    })
    self.websocket.send(deltaMessage)
  }
  createsession() {
    let createMessage = JSON.stringify({
      type: this.mt.CreateSession,
      body: {
        version: this.quill.version,
        delta: this.quill.getContents(),
      },
    })
    waitForConnection(self.websocket, function () {
      self.websocket.send(createMessage)
    })
  }
  joinsession() {
    let JoinMessage = JSON.stringify({
      type: this.mt.JoinSession,
      msg: 'user',
      password: this.password,
    })
    waitForConnection(self.websocket, function () {
      self.websocket.send(JoinMessage)
    })
  }
  leavesession() {
    self.websocket.send(JSON.stringify({ type: this.mt.LeaveSession }))
  }
  sendChat(data) {
    let chatMessage = JSON.stringify({
      id: this.id,
      type: this.mt.TXChat,
      chat: data,
    })
    self.websocket.send(chatMessage)
  }
}

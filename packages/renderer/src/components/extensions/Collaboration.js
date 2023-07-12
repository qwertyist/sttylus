/* eslint-disable no-continue */
/* eslint-disable no-unused-vars */
/* eslint-disable no-empty-function */
/* eslint-disable no-undefined */
/* eslint-disable guard-for-in */
/* eslint-disable class-methods-use-this */
import { Extension } from 'tiptap'
import { Step } from 'prosemirror-transform'
import { Decoration, DecorationSet } from 'prosemirror-view'
import {
    collab,
    sendableSteps,
    getVersion,
    receiveTransaction,
} from 'prosemirror-collab'

export default class Collaboration extends Extension {
    get name() {
        return 'collaboration'
    }

    init() {
        this.editor.on('init', ({ state }) => {
            // let other participants know you are here
            this.options.me.cursor = state.selection.anchor
            this.options.me.focused = state.selection.focused
            this.options.socket.emit('init', this.options.me)
        })

        this.sendUpdate = this.debounce(({ state, transaction }) => {
            const sendable = sendableSteps(state)
            this.options.me.cursor = state.selection.anchor
            this.options.me.focused = state.selection.focused

            if (sendable) {
                this.options.socket.emit('update', {
                    version: sendable.version,
                    steps: sendable.steps.map((step) => step.toJSON()),
                    clientID: this.options.clientID,
                    participant: this.options.me,
                })
            } else if (transaction.updated > 0) {
                this.options.socket.emit('cursorchange', this.options.me)
            }
        }, this.options.debounce)

        this.updateLocalCursors = (state) => {
            const sendable = sendableSteps(state)
            if (sendable) {
                for (let participantID in this.participants) {
                    let cursor = this.participants[participantID].cursor
                    if (
                        cursor != undefined &&
                        sendable.steps[0].slice != undefined &&
                        cursor >= sendable.steps[0].from
                    ) {
                        let gap = sendable.steps[0].from - sendable.steps[0].to
                        this.participants[participantID].cursor =
                            cursor + gap + sendable.steps[0].slice.content.size
                    }
                }
                this.options.updateCursors({ participants: this.participants })
            }
        }

        this.editor.on('transaction', ({ state, transaction }) => {
            //  this.updateLocalCursors(state);
            this.sendUpdate({ state, transaction })
        })
    }

    get defaultOptions() {
        return {
            me: {
                displayname: '',
            },
            socket: '',
            version: 0,
            clientID: '',
            debounce: 5,
            onSendable: () => {},
            update: ({ steps, version, participants }) => {
                const { state, view, schema } = this.editor

                if (getVersion(state) > version) {
                    return
                }

                view.dispatch(
                    receiveTransaction(
                        state,
                        steps.map((item) => Step.fromJSON(schema, item.step)),
                        steps.map((item) => item.clientID)
                    ).scrollIntoView()
                )
            },

            updateCursors: ({ steps, version, participants }) => {
                if (this.active) {
                    const { state, view, schema } = this.editor
                    this.participants = participants

                    // Set the decorations in the editor
                    let clientID = this.options.clientID
                    let props = {
                        decorations(state) {
                            let decos = []
                            if (participants != undefined) {
                                for (const [id, dec] of Object.entries(
                                    participants
                                )) {
                                    if (dec.cursor == undefined) {
                                        continue
                                    }
                                    let cursorclass = 'cursor'
                                    let displayname = dec.displayname
                                    let displaycolor =
                                        'style="background-color:' +
                                        dec.displaycolor +
                                        '; border-top-color:' +
                                        dec.displaycolor +
                                        '"'

                                    const dom = document.createElement('div')
                                    if (dec.focused == false) {
                                        cursorclass += ' inactive'
                                    }

                                    if (dec.clientID == clientID) {
                                        cursorclass += ' me'
                                    }
                                    if (displayname == false) {
                                        displayname = dec.clientID
                                    }

                                    dom.innerHTML =
                                        '<span class="' +
                                        cursorclass +
                                        '" ' +
                                        displaycolor +
                                        '>' +
                                        displayname +
                                        '</span>'
                                    dom.style.display = 'inline'
                                    dom.class = 'tooltip'
                                    decos.push(
                                        Decoration.widget(dec.cursor, dom)
                                    )
                                }
                            }
                            return DecorationSet.create(state.doc, decos)
                        },
                    }
                    view.setProps(props)
                }
            },
        }
    }

    get plugins() {
        return [
            collab({
                active: this.options.active,
                version: this.options.version,
                clientID: this.options.clientID,
            }),
        ]
    }

    debounce(fn, delay) {
        let timeout
        return function (...args) {
            if (timeout) {
                clearTimeout(timeout)
            }
            timeout = setTimeout(() => {
                fn(...args)
                timeout = null
            }, delay)
        }
    }
}

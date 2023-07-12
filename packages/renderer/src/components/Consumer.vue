<template>
    <div>
        <!--<div class="tiptapcontainer" :class="{ navMargin: nav }">
   <b-badge variant="primary">Användarvy {{ endpoint }}</b-badge>
    <template v-if="editor">
      <b-container fluid class="no-gutters">
        <b-row class="d-block">
          -->
        <editor-content
            id="editor"
            autocorrect="off"
            autocapitalize="off"
            spellcheck="false"
            class="editor__content"
            auto-focus="true"
            :style="{
                color: fontSettings.foreground,
                backgroundColor: fontSettings.background,
                fontSize: fontSettings.size + 'px',
                fontFamily: 'Arial',
                height: widgetHeight,
            }"
            :editor="editor"
            ref="consumer"
        />
        <ConsumerSettings />
    </div>
    <!--
        </b-row>
      </b-container>
    </template>
    <em v-else>Laddar programmet...</em>
  </div>-->
</template>

<script>
//import io from "socket.io-client";
import EventBus from '../eventbus.js'
import { Editor, EditorContent } from 'tiptap'
import {
    HardBreak,
    Bold,
    Italic,
    Underline,
    History,
    // Collaboration,
} from 'tiptap-extensions'
import { Selection } from 'prosemirror-state'
import Collaboration from './extensions/Collaboration.js'
import ExpandWord from './extensions/ExpandWord.js'
import Manuscript from './extensions/Manuscript.js'
import Hotkeys from './extensions/Hotkeys.js'
import ConsumerSettings from './settings/ConsumerSettings.vue'

export default {
    name: 'Tiptapexample',
    components: {
        EditorContent,
        ConsumerSettings,
    },
    props: {
        nav: Boolean,
        endpoint: String,
    },
    data() {
        return {
            editor: null,
            socket: null,
            authed: false,
            count: 0,
            me: {
                displayname: '',
            },
            participants: {},
            oldPosition: 0,
            fontSettings: {
                family: '',
                size: 32,
                foreground: '#ffff00',
                background: '#000000',
            },
        }
    },
    methods: {
        onInit({ doc, version }) {
            if (this.editor) {
                this.editor.destroy()
            }
            this.editor = new Editor({
                content: doc,
                autoFocus: 'end',
                editable: false,
                fontSettings: {},
                editorProps: {
                    handleScrollToSelection: ({ state }) => {
                        this.$refs.consumer.$el.scrollTop =
                            this.$refs.consumer.$el.scrollHeight
                    },
                },
                extensions: [
                    new Hotkeys({
                        sizeUp: this.sizeUp,
                        sizeDown: this.sizeDown,
                        clearDoc: this.clearDoc,
                        addAbb: this.addAbb,
                        openNav: this.toggleNav,
                    }),
                    // new HardBreak(),
                    new Bold(),
                    new Italic(),
                    new Underline(),
                    new History(),
                    new Manuscript(),
                    new Collaboration({
                        socket: this.socket,
                        clientID: this.socket.id,
                        // version is an integer which is incremented with every change
                        version,
                        // debounce changes so we can save some requests
                        debounce: 0,
                        // onSendable is called whenever there are changed we have to send to our server
                        onSendable: ({ sendable }) => {
                            // this.socket.emit('update', sendable)
                        },
                    }),
                ],
            })
        },

        setCount(count) {
            this.count = count
        },

        setParticipants(participants) {
            this.participants = participants
        },
        focusAt(pos) {
            this.editor.focus(pos)
        },
        refocus() {
            this.editor.focus('end')
        },
        updateSettings() {
            const fontSettings = JSON.parse(
                localStorage.getItem('fontSettings')
            )
            this.fontSettings = fontSettings
            document.activeElement.blur()
            this.editor.focus()
            setTimeout(() => {
                const widget = document.querySelector('.ProseMirror')
                widget.click()
            }, 30)
        },
        sizeDown() {
            this.fontSettings.size -= 4
            this.$store.commit('setFontSettings', this.fontSettings)
            this.refreshEditor()
        },
        sizeUp() {
            this.fontSettings.size += 4
            this.$store.commit('setFontSettings', this.fontSettings)
            this.refreshEditor()
        },
        setCollabSocket(endpoint) {
            this.socket = WebSocket()
                // get the current document and its version
                .on('init', (data) => {
                    this.onInit(data)
                    let who = {
                        type: 'interpreter',
                    }
                    this.socket.emit('ident', who)
                })
                // send all updates to the collaboration extension
                .on('update', (data) => {
                    if (data.steps[0].clientID !== this.socket.id) {
                        // focus at this step
                        this.focusAt(data.steps[0].step.to)
                    }
                    this.editor.extensions.options.collaboration.update(data)
                    this.editor.extensions.options.collaboration.updateCursors(
                        data
                    )
                })
                // get count of connected users
                .on('getCount', (count) => this.setCount(count))
                // update Cursor position of collaborators
                .on('cursorupdate', (data) => {
                    this.editor.extensions.options.collaboration.updateCursors(
                        data
                    )
                    this.setParticipants(data.participants)
                })
                .on('session', (Id) => {
                    this.$store.commit('setMySessionId', Id)
                })
                // .on('connect', c => {})
                .on('disconnect', (c) => {
                    console.log('disconnected:', c)
                })
                .on('connect_failed', (c) => {
                    console.log('failed connecting:', c)
                })
                .on('connect_error', (c) => {
                    console.log('server is offline:', c)
                })
        },
        stopCollab() {
            this.socket = null
            this.$forceUpdate()
        },
    },
    computed: {
        widgetHeight() {
            return '100vh !important'
        },
    },
    mounted() {
        EventBus.$on('refocus', this.refocus)
        EventBus.$on('updateSettings', this.updateSettings)

        this.fontSettings = this.$store.state.settings.font
        this.currentLists = this.$store.state.selectedLists
        this.setCollabSocket(this.endpoint + ':3000')
    },
    beforeDestroy() {
        EventBus.$off('refocus')
        EventBus.$off('updateSettings')
        this.editor.destroy()
        if (this.socket) {
            this.socket.destroy()
        }
    },
}
</script>

<style scoped lang="scss">
.html,
body {
    height: 100%;
    white-space: nowrap;
}
.selected {
    font-weight: bold;
}
.count {
    display: flex;
    align-items: center;
    font-weight: bold;
    color: #27b127;
    margin-bottom: 1rem;
    text-transform: uppercase;
    font-size: 0.7rem;
    line-height: 1;
    &:before {
        content: '';
        display: inline-flex;
        background-color: #27b127;
        width: 0.4rem;
        height: 0.4rem;
        border-radius: 50%;
        margin-right: 0.3rem;
    }
}

.cursor.me {
    display: none;
    /*background-color: #F55;*/
}
.cursor.inactive {
    opacity: 0.5;
}
.cursor.me::after {
    display: none;
    border-color: inherit;
}
.cursor.inactive::after {
    opacity: inherit;
    border-color: inherit;
}

.cursor {
    /*background-color: #555;*/
    color: #fff;
    text-align: center;
    border-radius: 6px 6px 6px 0px;
    padding: 5px;
    margin-left: -4.5px;
    position: absolute;
    z-index: 1;
    bottom: 5px;
    left: -50%;
    opacity: 0.85;
}

.cursor::after {
    content: '';
    position: absolute;
    top: 100%;
    left: 0%;
    border-width: 5px;
    border-style: solid;
    border-color: inherit;
    /*border-color: #555 transparent transparent transparent;*/
    color: transparent;
}

.editor__content {
    border: 0px;
    padding: 0rem 0.75rem;
    background-color: #f5f5f5;
    font-size: 2em;
    overflow-y: scroll;
}

.bv-row {
    padding-top: 20px;
}

.ProseMirror-widget {
    position: absolute;
    width: 0.1px;
    /*border-style: solid;*/
}
.ProseMirror:focus {
    outline: none;
}
/*.ProseMirror {
  height: 95vh !important;
  width: 95vw !important;
}
*/
.navMargin {
    position: fixed;
    top: 5vh;
    height: 75vh !important;
    width: 100vw !important;
}

button.is-active {
    color: #fff;
    background-color: #ddd;
    font-weight: bold;
}

.p {
    margin-bottom: 0rem !important;
}

.btn {
    margin: 5px;
    border-color: #ddd;
}
</style>

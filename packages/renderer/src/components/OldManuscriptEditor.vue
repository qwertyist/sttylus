<template>
    <div>
        <b-container fluid class="no-gutters">
            <b-row class="d-block">
                <div class="editor">
                    <editor-content
                        autocomplete="off"
                        id="editor"
                        autocorrect="off"
                        autocapitalize="off"
                        spellcheck="false"
                        class="editor__content"
                        :style="{
                            height: 75 + 'vh',
                            color: fontSettings.foreground + '!important',
                            backgroundColor:
                                fontSettings.background + '!important',
                            fontSize: fontSettings.size + 'px',
                            fontFamily: fontSettings.fontFamily,
                        }"
                        :editor="editor"
                        ref="editor"
                    />

                    <editor-menu-bar :editor="editor">
                        <div class="menubar">
                            <!--
              <button
                class="menubar__button btn btn-light"
                :class="{ 'is-active': isActive.bold() }"
                @click="commands.bold"
              >
                <i class="fa fa-bold"></i>
              </button>

              <button
                class="menubar__button btn btn-light"
                :class="{ 'is-active': isActive.italic() }"
                @click="commands.italic"
              >
                <i class="fa fa-italic"></i>
              </button>

              <button
                class="menubar__button btn btn-light"
                :class="{ 'is-active': isActive.underline() }"
                @click="commands.underline"
              >
                <i class="fa fa-underline"></i>
              </button>
              -->
                            <b-button
                                variant="primary"
                                @click="saveManuscript(true)"
                                >Spara</b-button
                            >
                        </div>
                    </editor-menu-bar>
                </div>
            </b-row>
            <AddAbbreviation ref="addAbb" />
            <SaveManuscript ref="saveManuscript" />
        </b-container>
    </div>
</template>

<script>
import api from '../api/api.js'
import EventBus from '../eventbus.js'

import { Editor, EditorContent, EditorMenuBar, Extension } from 'tiptap'
import {
    HardBreak,
    Heading,
    Bold,
    Code,
    Italic,
    Strike,
    Underline,
    Blockquote,
    History,
} from 'tiptap-extensions'
import ExpandWord from './extensions/ExpandWord.js'
import Hotkeys from './extensions/Hotkeys.js'
import AddAbbreviation from './modals/AddAbbreviation.vue'
import SaveManuscript from './modals/SaveManuscript.vue'

export default {
    name: 'ManuscriptEditor',
    components: {
        EditorContent,
        EditorMenuBar,
        AddAbbreviation,
        SaveManuscript,
    },
    props: ['id'],
    data() {
        return {
            fontSettings: {
                family: '',
                size: 32,
                foreground: '#000000',
                background: '#ffffff',
            },
            original: null,
            export: null,
            editor: new Editor({
                autoFocus: true,
                extensions: [
                    new Hotkeys({
                        sizeUp: this.sizeUp,
                        sizeDown: this.sizeDown,
                        clearDoc: this.printEditorContent,
                        addAbb: this.addAbb,
                        openNav: this.toggleNav,
                    }),
                    //new HardBreak(),
                    /*          new Bold(),
          new Italic(),
          new Underline(),
          //new OrderedList(),
          //new BulletList(),
  
  */
                    new History(),
                    new ExpandWord({
                        userId: this.$store.state.userData.id,
                        manuscripts: [],
                        onKeyDown: ({ view, event, collapse, focusLast }) => {
                            this.oldPosition = view.state.doc.content.size - 1

                            if (focusLast) {
                                this.editor.setSelection(0, 0)

                                this.editor.focus('end')
                                return
                            }
                        },
                    }),
                ],
                onUpdate: ({ getJSON }) => {
                    this.export = getJSON()
                },
                onInit: () => {},
            }),
        }
    },
    methods: {
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
        clearDoc() {
            this.editor.setContent('')
        },
        addAbb(query) {
            this.$refs['addAbb'].showModal(query)
        },
        toggleNav() {
            this.$root.$emit('toggleNav')
        },
        refreshEditor(doc) {
            this.editor.focus()
        },
        refocus(wait) {
            console.log('Wait:', wait)
            this.editor.focus('end')
        },
        printEditorContent() {
            let json = this.editor.getJSON()

            console.log(json)
        },
        getEditorContent() {
            this.export = this.editor.getJSON()
        },
        leaveManuscriptEditor() {
            this.getEditorContent()
            console.log('original:', this.original, 'export:', this.export)
            if (
                this.export != null &&
                this.isEqual(this.original, this.export) == false &&
                this.export !== '<p></p>'
            ) {
                EventBus.$emit('unsavedManuscriptChanges', true)
            } else {
                EventBus.$emit('unsavedManuscriptChanges', false)
            }
        },
        stopEditingManuscript(close) {
            console.log('Stäng:', close)
            this.getEditorContent()
            // console.log("loaded content:", this.content);
            //  console.log("current content:", this.export);
            this.saveManuscript(close)
        },
        saveManuscript(close) {
            this.export = this.editor.getJSON()
            this.preview = this.editor.getHTML()
            console.log(this.export)
            if (this.export == '<p></p>') {
                this.$toast.warning('Du kan inte spara ett tomt manuskript')
                return
            }
            if (this.id == 'parsed') {
                this.$refs['saveManuscript'].showModal()
            } else if (this.id == 'new') {
                this.$refs['saveManuscript'].showModal({
                    id: '',
                    preview: this.preview,
                    body: this.export,
                    close: close,
                })
            } else {
                this.$refs['saveManuscript'].showModal({
                    id: this.id,
                    preview: this.preview,
                    body: this.export,
                    close: close,
                })
            }
        },
        addEventListeners() {
            EventBus.$on('leaveManuscriptEditor', this.leaveManuscriptEditor)
            EventBus.$on('refocus', this.refocus)
            EventBus.$on('stopEditingManuscript', this.stopEditingManuscript)
            EventBus.$on('saveManuscript', this.saveManuscript)
        },
        removeEventListeners() {
            EventBus.$off('leaveManuscriptEditor')
            EventBus.$off('refocus')
            EventBus.$off('stopEditingManuscript')
            EventBus.$off('saveManuscript')
        },
    },
    watch: {
        id(val) {
            console.log('content id:', val)
        },
    },
    mounted() {
        this.addEventListeners()
        api.cacheAbbs()
            .then((resp) => {
                EventBus.$emit('getAbbCache')
            })
            .catch((err) => {
                console.log("manuscript editor couldn't create abbs cache")
            })
    },
    beforeDestroy() {
        this.removeEventListeners()
    },
}
</script>

<style scoped lang="scss"></style>

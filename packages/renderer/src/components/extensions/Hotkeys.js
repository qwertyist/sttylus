/* eslint-disable no-control-regex */
/* eslint-disable no-empty-function */
/* eslint-disable class-methods-use-this */
import { Extension } from 'tiptap'
import EventBus from '../../eventbus'

export default class Hotkeys extends Extension {
    get name() {
        return 'hotkeys'
    }
    get defaultOptions() {
        return {
            sizeUp: () => {},
            sizeDown: () => {},
            clearDoc: () => {},
        }
    }
    keys() {
        return {
            ESC: () => {
                this.options.stopPrompt()

                this.editor.focus('end')
                return true
            },
            ENTER: () => {
                this.options.expand()
                return true
            },
            F4: () => {
                this.options.clearDoc()
                return true
            },
            F6: () => {
                // return true prevents default behaviour
                this.options.sizeDown()
                return true
            },
            F7: () => {
                this.options.sizeUp()

                // return true prevents default behaviour
                return true
            },
            F1: () => {
                this.options.openSupport()
                return true
            },
            'SHIFT-F2': (state) => {
                const $position = state.tr.selection.$from

                if (state.tr.selection.$from !== state.tr.selection.$to) {
                    const selection = $position.doc.textBetween(
                        state.tr.selection.from,
                        state.tr.selection.to,
                        '\0',
                        '\0'
                    )

                    this.options.addAbb(
                        selection
                            .trim()
                            .replace(/[\x00-\x1F\x7F-\x9F]/g, '')
                            .toLowerCase()
                    )
                    return true
                }

                const anchor = state.tr.selection.anchor

                try {
                    const textFrom = $position.before()
                    const text = $position.doc.textBetween(
                        textFrom,
                        anchor,
                        '\0',
                        '\0'
                    )

                    const wordsToAnchor = text.split(/[ :;\-()[\]_{}\uFEFF]+/)
                    const query = wordsToAnchor[wordsToAnchor.length - 1]
                    this.options.addAbb(query.trim().toLowerCase())
                } catch {
                    this.options.addAbb('')
                }

                return true
            },
            F2: (state) => {
                const $position = state.tr.selection.$from

                if (state.tr.selection.$from !== state.tr.selection.$to) {
                    const selection = $position.doc.textBetween(
                        state.tr.selection.from,
                        state.tr.selection.to,
                        '\0',
                        '\0'
                    )

                    this.options.addAbb(
                        selection.trim().replace(/[\x00-\x1F\x7F-\x9F]/g, '')
                    )
                    return true
                }

                const anchor = state.tr.selection.anchor

                try {
                    const textFrom = $position.before()
                    const text = $position.doc.textBetween(
                        textFrom,
                        anchor,
                        '\0',
                        '\0'
                    )

                    const wordsToAnchor = text.split(/[ :;\-()[\]_{}\uFEFF]+/)
                    const query = wordsToAnchor[wordsToAnchor.length - 1]
                    this.options.addAbb(query.trim())
                } catch {
                    this.options.addAbb('')
                }

                return true
            },
            F3: (state) => {
                console.log('Should open right pane')
                this.options.toggleChat()
                return true
            },
            F8: () => {
                this.options.changeTextColor()
                return true
            },
            F10: () => {
                this.options.openNav('')
                return true
            },
            Tab: () => {
                return true
            },
        }
    }
}

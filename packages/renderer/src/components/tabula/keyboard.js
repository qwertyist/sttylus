import Quill from "quill"
import Delta from "quill-delta"
import api from "../../api/api";
import { store } from "../../store"
import EventBus from "../../eventbus"
import { promptManuscript, queryManuscript, removePreviews } from "./manuscript";
var Keyboard = Quill.import('modules/keyboard');


let capitalizeNext = true;
let abbreviated = false;
let functionKeys = ["F1", "F2", "F3", "F4", "F5", "F6", "F7", "F8", "F12"]

let separators = [" ", ".", "!", "?", ":", ";", "-", "(", ")", "[", "]", "{", "}", "Enter"]

export default class keyboard extends Keyboard {
    static DEFAULTS = {
        ...Keyboard.DEFAULTS,
        bindings: {
            ...Keyboard.DEFAULTS.bindings,
            ['list autofill']: undefined,
        }
    }
    constructor(quill, options) {
        super(quill, options)
        this.abbreviated = false;
        this.querying = null;
        this.capitalizeNext = 1;
        this.scrollIntoView = false;
        this.manuscriptEditor = options.manuscriptEditor
        this.URL = false
        this.addKeybindings()
        this.scrollHandler(quill.root)
        this.prompt = ""
        this.lastKey = ""
        this.cache = new Map();
    }
    listen() {
        this.quill.root.addEventListener("keydown", e => {
            if (this.URL) {
                if (e.key != "." && separators.indexOf(e.key) !== -1) {
                    this.capitalizeNext = true
                    this.URL = false
                    return
                } else {
                    this.capitalizeNext = false
                    return
                }
            }
            if (e.key.length == 1) {
                if (this.capitalizeNext) {
                    if (!e.ctrlKey && !e.altKey) {
                        let match = e.key.match(/\p{Letter}+/gu)
                        if (match) {
                            let letter = match[0].toUpperCase()
                            e.preventDefault()
                            let sel = this.quill.getSelection()
                            this.quill.insertText(sel.index, letter)
                            this.capitalizeNext = false
                        }
                    }
                } else if (e.shiftKey) { this.capitalizeNext = false }
            }

        })
        super.listen()

    }

    cacheAbb(abb) {
        if (abb.capitalize) {
            abb.word = abb.word.slice(0, 1).toLowerCase() + abb.word.slice(1)
        }
        this.cache.set(abb.abb, abb.word)
    }

    unloadAbb(abb) {

    }

    wordBeforeCursor(prefix) {
        return prefix.split(/[\u200B\s-\/"'()]/).pop()
    }

    capitalize(word) {
        return word.charAt(0).toUpperCase() + word.slice(1);
    }


    abbreviate(index, abb, abbreviator, quill) {
        //     console.log("abbreviate at index:", index, "with abb:", abb, "and sep:", abbreviator)
        if (!this.manuscriptEditor && queryManuscript(abb)) {
            quill.deleteText(index - abb.length, abb.length)
            this.prompt = abb
        } else {
            setTimeout(() => {
                console.log("Querying");
                this.querying = false
            }, 125)
            if (!this.querying) {
                this.querying = true;
                api.abbreviate(abb)
                    .then(r => {
                        let word = r.data.word;
                        if (r.status == 204) {
                            console.log("No content")
                            word = abb
                        }
                        if (r.data.missed) {
                            console.log("Missed abb!")
                            store.commit("incrementMissedAbb", { word: abb, abb: word })
                            word = abb

                        }

                        this.insertAbbreviation(index, abb, abbreviator, word, quill)
                        setTimeout(() => quill.setSelection(quill.getSelection().index, 0), 0)
                    }).catch(err => {
                        this.insertAbbreviation(index, abb, abbreviator, abb, quill)
                    })
            } else {
                quill.insertText(index, abbreviator)
            }
        }
    }

    insertAbbreviation(index, abb, abbreviator, word, quill) {
        let format = quill.getFormat(index)
        if (abb == word) {
            quill.insertText(index, abbreviator, format)
            return
        }
        let delta = new Delta().retain(index - abb.length).delete(abb.length).insert(word + abbreviator, format)
        quill.updateContents(delta)

    }

    scrollHandler(node) {
        if (!this.manuscriptEditor) {
            node.addEventListener("keyup", ({ key: key, target: { lastChild } }) => {
                this.lastKey = key
                if (this.scrollIntoView) {
                    lastChild.scrollIntoView()
                    this.scrollIntoView = false
                }
            })
        }
    }

    addKeybindings() {
        // TAB
        this.bindings[9].unshift({
            key: 9,
            handler: function (range, context) {
                console.log("TAB");
                this.quill.insertText(range.index, "\u200B")
            }
        })
        //ESCAPE
        this.addBinding({
            key: 27,
            handler: function (range, context) {
                if (store.getters.getModalOpen) {
                  console.log("MODAL open - dont catch")
                }
                let end = this.quill.getText().length - 1
                if (this.lastKey == "Escape") {
                    this.prompt = ""
                    removePreviews(range.index, this.quill)
                }
                this.quill.setSelection(end)
                return true
            }
        })
        //F2
        this.addBinding({
            key: 113,
            shiftKey: false,
            handler: function (range, context) {
                let phrase = ""
                if (range.length > 0) {
                    phrase = this.quill.getText(range.index, range.length)
                } else {
                    phrase = context.prefix.split(" ").pop()
                }
                EventBus.$emit("addAbbreviation", phrase)
            }
        })
        //Shift+F2
        this.addBinding({
            key: 113,
            shiftKey: true,
            handler: function (range, context) {
                let phrase = ""
                if (range.length > 0) {
                    phrase = this.quill.getText(range.index, range.length)
                } else {
                    phrase = context.prefix.split(" ").pop()
                }
                EventBus.$emit("addAbbreviation", phrase.toLowerCase())
            }
        })
        //F3
        this.addBinding({
            key: 114,
            handler: function () {
              EventBus.$emit("sendReadySignal")
            }
        })
        //F4
        this.addBinding({
            key: 115,
            handler: function (range, context) {
                removePreviews(range.index, this.quill)
                console.log("F4")
                this.url = false;
                this.quill.setText("");
                this.capitalizeNext = true
                EventBus.$emit("clear")
                this.prompt = ""
            }
        })
        //F6 sizeDown
        this.addBinding({
            key: 117,
            ctrlKey: false,
            handler: function () {
                EventBus.$emit("sizeChange", { inc: false, send: false})
            }
        })
        //ctrl+F6 send sizeDown
        this.addBinding({
            key: 117,
            ctrlKey: true,
            handler: function () {
                EventBus.$emit("sizeChange", { inc: false, send: true})
            }
        })
        //F7 sizeUp
        this.addBinding({
            key: 118,
            ctrlKey: false,
            handler: function () {
                EventBus.$emit("sizeChange", { inc: true, send: false})
            }
        })
        //CTRL+F7 send sizeUp
        this.addBinding({
            key: 118,
            ctrlKey: true,
            handler: function () {
                console.log("keybinding ctrl+f7")
                EventBus.$emit("sizeChange", { inc: true, send: true })
            }
        })
        //F8 colorChange
        this.addBinding({
            key: 119,
            ctrlKey: false,
            handler: function () {
                EventBus.$emit("colorChange", false)
            }
        })
        //ctrl+F8 send colorChange
        this.addBinding({
            key: 119,
            ctrlKey: true,
            handler: function () {
                EventBus.$emit("colorChange", true)
            }
        })
        //F9 (Create) 
        this.addBinding({
            key: 120,
            handler: function () {
                console.log("keybinding F9")
                EventBus.$emit("createSession", true)
            }
        })
        //F10 (Join)
        this.addBinding({
            key: 121,
            handler: function () {
                EventBus.$emit("joinSession")
            }
        })

        //Space
        this.addBinding({
            key: 32,
            shiftKey: null,
            handler: function (range, context) {
                if (this.abbreviated) {
                    this.abbreviated = false;
                    return true
                }
                let abb = this.wordBeforeCursor(context.prefix)
                //console.log(abb.length)
                //console.log("word before cursor:", abb)
                if (!abb) {
                    this.abbreviated = false;
                    return true
                }
                //abb = this.capitalize(abb) 
                this.abbreviate(range.index, abb, " ", this.quill)
            }
        })
        //Ctrl+Space
        this.addBinding({
            key: 32,
            ctrlKey: true,
            handler: function (range, context) {
                let abb = this.wordBeforeCursor(context.prefix)
                //console.log(abb.length)
                //console.log("word before cursor:", abb)
                this.abbreviated = false
                this.capitalize = false
                this.abbreviate(range.index, abb, "\u200B", this.quill)
            }
        })
        //Enter
        this.bindings[13].unshift({
            key: 13,
            handler: function (range, context) {

                if (this.prompt != "") {
                    return
                }
                this.url = false
                let scroll = range.index == this.quill.getLength() - 1
                let abb = this.wordBeforeCursor(context.prefix)
                if (abb != "") {
                    this.abbreviate(range.index, abb, "\n", this.quill)
                } else {
                    this.quill.insertText(range.index, "\n")
                }
                if (this.options.capitalizeOnNewLine) {
                    this.capitalizeNext = true
                } else {
                    this.capitalizeNext = false
                }
                this.abbreviated = false
                this.scrollIntoView = scroll
                EventBus.$emit("newLine", scroll)
            }

        })

        //Shift+Enter
        this.bindings[13].unshift({
            key: 13,
            shiftKey: true,
            handler: function (range, context) {
                let abb = this.wordBeforeCursor(context.prefix)
                this.url = false
                if (abb) {
                    this.abbreviate(range.index, abb, "\n", this.quill)
                    this.capitalizeNext = true
                } else {
                    this.quill.insertText(range.index, "\n")
                }
                if (this.options.capitalizeOnNewLine) {
                    this.capitalizeNext = false
                } else {
                    this.capitalizeNext = true
                }
                this.scrollIntoView = true
            }
        })
      
        //Backspace
        this.addBinding({
          key: 8,
          handler: function () {
            this.capitalizeNext = false
            return true
          }
        })

        this.addBinding({
          key: 8,
          ctrlKey: true,
          handler: function (range, context) {
            this.capitalizeNext = false
            return true
          }
        })
        //Period
        this.addBinding({
            key: 190,
            handler: function (range, context) {
                let abb = this.wordBeforeCursor(context.prefix)
                if (!abb.trim()) {
                    return true
                }
                this.capitalizeNext = true
                this.abbreviated = true;
                this.URL = true
                this.abbreviate(range.index, abb, ".", this.quill)
            }
        })
        //Colon
        this.addBinding({
            key: 190,
            shiftKey: true,
            handler: function (range, context) {
                if (context.prefix.split(" ").length == 1) {
                    this.capitalizeNext = true;
                }
                let abb = this.wordBeforeCursor(context.prefix)
                this.abbreviate(range.index, abb, ":", this.quill)
            }
        })

        //Comma
        this.addBinding({
            key: 188,
            handler: function (range, context) {
                let abb = this.wordBeforeCursor(context.prefix)
                this.abbreviate(range.index, abb, ",", this.quill)
            }
        })
        //Semicolon
        this.addBinding({
            key: 188,
            shiftKey: true,
            handler: function (range, context) {
                let abb = this.wordBeforeCursor(context.prefix)
                this.abbreviate(range.index, abb, ";", this.quill)
            }
        })
        //Dash
        this.addBinding({
            key: 189,
            handler: function (range, context) {
                let abb = this.wordBeforeCursor(context.prefix)
                this.abbreviate(range.index, abb, "-", this.quill)
            }
        })
        //Exclamation
        this.addBinding({
            key: 49,
            shiftKey: true,
            handler: function (range, context) {
                let abb = this.wordBeforeCursor(context.prefix)
                this.abbreviate(range.index, abb, "!", this.quill)
                this.capitalizeNext = true
                this.abbreviated = true;
            }
        })
        //Question mark
        this.addBinding({
            key: 187,
            shiftKey: true,
            handler: function (range, context) {
                let abb = this.wordBeforeCursor(context.prefix)
                this.abbreviate(range.index, abb, "?", this.quill)
                this.capitalizeNext = true
                this.abbreviated = true;
            }
        })
        // Forwardslash
        this.addBinding({
            key: 55,
            shiftKey: true,
            handler: function (range, context) {
                let abb = this.wordBeforeCursor(context.prefix)
                this.abbreviate(range.index, abb, "/", this.quill)
                this.capitalizeNext = false;
                this.abbreviated = true;
            }
        })
        // Double quotation mark " 
        this.addBinding({
            key: 50,
            shiftKey: true,
            handler: function (range, context) {
                let abb = this.wordBeforeCursor(context.prefix)
                this.abbreviate(range.index, abb, "\"", this.quill)
                this.capitalizeNext = false;
                this.abbreviated = true;
            }
        })
        // Quotation mark ' 
        this.addBinding({
            key: 191,
            handler: function (range, context) {
                console.log("\' quotation mark")
                let abb = this.wordBeforeCursor(context.prefix)
                this.abbreviate(range.index, abb, "'", this.quill)
                this.capitalizeNext = false;
                this.abbreviated = true;
            }
        })
        //Ctrl+Backspace
        this.bindings[8].unshift({
            key: 8,
            ctrlKey: true,
            handler: function (range, context) {
                console.log("prefix:", context.prefix, "length.", context.prefix.length)
                if (context.prefix.trim().split(" ").length == 1) {
                    this.capitalizeNext = true
                }
                return true
            }
        })
        // Right →
        this.bindings[39].unshift({
            key: 39,
            handler: function (range, context) {
                if (this.prompt != "") {
                    let index = this.quill.getText().length - 1
                    let offset = promptManuscript(index, this.quill, "insertWord")

                    if (offset == -1) {
                        this.prompt = false
                        return false
                    }
                    this.quill.setSelection(index + offset)
                    return false
                }
                return true
            }
        })
        // Ctrl+Right →
        this.bindings[39].unshift({
            key: 39,
            ctrlKey: true,
            handler: function (range, context) {
              if (range.index == this.quill.getText().length-1) return false
              return true
            }
        })
        // ALT+Right →

        this.bindings[39].unshift({
            key: 39,
            altKey: true,
            handler: function (range, context) {
                if (this.prompt != "") {
                    let index = this.quill.getText().length - 1
                    let offset = promptManuscript(index, this.quill, "skipWord")

                    if (offset == -1) {
                        this.prompt = false
                        return false
                    }
                    this.quill.setSelection(index + offset)
                    return false
                }
                return true
            }
        })

        // Down ↓
        this.addBinding({
            key: 40,
            handler: function (range, context) {
                if (this.prompt != "") {
                    let index = this.quill.getText().length - 1
                    let offset = promptManuscript(index, this.quill, "insertClause")

                    if (offset == -1) {
                        this.prompt = false
                        return false
                    }
                    this.quill.setSelection(index + offset)
                    return false
                }
                return true
            }
        })
        // ALT+Down ↓
        this.addBinding({
            key: 40,
            altKey: true,
            handler: function (range, context) {
                if (this.prompt != "") {
                    let index = this.quill.getText().length - 1
                    let offset = promptManuscript(index, this.quill, "skipClause")

                    if (offset == -1) {
                        this.prompt = false
                        return false
                    }
                    this.quill.setSelection(index + offset)
                    return false
                }
                return true
            }
        })
        //Left
        this.bindings[37].unshift({
            key: 37,
            handler: function (range, context) {
                if (this.prompt != "") {
                    let index = this.quill.getText().length - 1
                    let offset = promptManuscript(index, this.quill, "removeWord", context)
                    this.quill.setSelection(offset)
                    return false
                }
                return true
            }
        })
        this.addBinding({
            key: 86,
            ctrlKey: true,
            handler: function (range, context) {
                return this.manuscriptEditor ? true : false;
            }
        })
      //
      //CTRL+1
      this.addBinding({
        key: 49,
        ctrlKey: true,
        handler: function (range, context) {
          if(store.getters.getModalOpen == false) {
            EventBus.$emit("changeStandardList", 1)
            return false
          }
        }
      })

      //CTRL+2
      this.addBinding({
        key: 50,
        ctrlKey: true,
        handler: function (range, context) {
          if(store.getters.getModalOpen == false) {
            EventBus.$emit("changeStandardList",  2)
            return false
          }
        }
      })

      //CTRL+3
      this.addBinding({
        key: 51,
        ctrlKey: true,
        handler: function (range, context) {
          if(store.getters.getModalOpen == false) {
            EventBus.$emit("changeStandardList",  3)
            return false
          }
        }
      })

      //CTRL+4
      this.addBinding({
        key: 52,
        ctrlKey: true,
        handler: function (range, context) {
          if(store.getters.getModalOpen == false) {
            EventBus.$emit("changeStandardList",  4)
            return false
          }
        }
      })

      //CTRL+5
      this.addBinding({
        key: 53,
        ctrlKey: true,
        handler: function (range, context) {
          if(!store.getters.getModalOpen) {
            EventBus.$emit("changeStandardList",  5)
            return false
          }
        }
      })
    }
}

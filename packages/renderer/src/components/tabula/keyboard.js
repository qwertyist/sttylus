import Quill from 'quill'
import Delta from 'quill-delta'
import api from '../../api/api'
import db from '../../store/db'
import { store } from '../../store'
import EventBus from '../../eventbus'
import { promptManuscript, queryManuscript, removePreviews } from './manuscript'
var Keyboard = Quill.import('modules/keyboard')

let separators = [
  ' ',
  '.',
  '!',
  '?',
  ':',
  ';',
  '-',
  '(',
  ')',
  '[',
  ']',
  '{',
  '}',
  'Enter',
]

export default class keyboard extends Keyboard {
  static DEFAULTS = {
    ...Keyboard.DEFAULTS,
    bindings: {
      ...Keyboard.DEFAULTS.bindings,
      ['list autofill']: undefined,
    },
  }
  constructor(quill, options) {
    super(quill, options)
    this.instance = (Math.random() + 1).toString(36).substring(7)
    this.abbreviated = false
    this.querying = null
    this.capitalizeNext = true
    this.scrollIntoView = false
    this.manuscriptEditor = options.manuscriptEditor
    this.URL = false
    this.addKeybindings()
    this.scrollHandler(quill.root)
    this.prompt = ''
    this.lastKey = ''
    this.cache = null
    this.getAbbCache()
    this.timeout = null
    this.currentWord = ''
    this.zeroWidthSpace = false
    this.dontExpand = false
    this.dontExpandWord = ''
  }

  listen() {
    /*
    if (this.timeout) {
      clearTimeout(this.timeout)
      this.timeout = null
    }
    this.timeout = setTimeout(() => {
      EventBus.$emit('previewAbb', this.insertHint())
    }, 2500)
    */
    this.quill.root.addEventListener('keydown', (e) => {
      if (this.URL) {
        if (e.key != '.' && separators.indexOf(e.key) !== -1) {
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
            let match = e.key.match(/\p{Letter}|\p{Number}|"+/gu)
            if (match) {
              let letter = match[0].toUpperCase()
              e.preventDefault()
              let sel = this.quill.getSelection()
              this.quill.insertText(sel.index, letter)
              this.capitalizeNext = false
            }
          }
        } else if (e.shiftKey) {
          this.capitalizeNext = false
        }
        if (e.key == 'v' && e.ctrlKey) {
          console.log('paste!')
        }
        return
      }
    })
    super.listen()
  }

  /*
   * Alternativ 1:
   * Om listan är uppdaterad genom "TouchList" i backend
   * Om cachad listas längd är annorlunda än den sparade
   * Ladda ner ny version av förkortningar från lista och cacha
   *
   * Gör cachen i webbläsaren något mer permanent.
   * -- Vad händer när man byter lista med ctrl+1..5, då kommer den att cacha om varje gång
   * -- och det är onödigt resurskrävande...
   *
   * Alternativ 2:
   * Listor lagras i sin helhet i webbläsaren och uppdateras enligt reglerna ovan.
   *
   *
   */
  getAbbCache() {
    if (store.state.synced) {
      console.log('last sync ok')
      this.cache = new Map()
      db.getAbbCache()
        .then((cache) => {
          this.cache = cache

          console.log('set up db abb cache')
        })
        .catch((err) => console.error(err))
    } else {
      api
        .getAbbCache()
        .then((resp) => {
          this.cache = new Map(Object.entries(resp.data))
          //console.log(this.instance, this.cache)
          console.log('set up api abb cache')
        })
        .catch((err) => {
          console.error("couldn't get cached abbs", err)
        })
    }
  }

  /*unloadAbb(abb) {

    }*/

  wordBeforeCursor(prefix) {
    this.currentWord = ''
    if (this.dontExpandWord != '') {
      console.log('dontexpandword', this.dontExpandWord)
      this.currentWord = prefix.split(this.dontExpandWord).pop()
    }
    if (this.currentWord == '') {
      this.currentWord = prefix.split(/[\u200B\s-.,:;_!?\/"'()]/).pop()
    }
    this.dontExpandWord = ''
    this.zeroWidthSpace = prefix.indexOf('\u200B') != -1

    return this.currentWord
  }

  capitalize(word) {
    return word.charAt(0).toUpperCase() + word.slice(1)
  }

  abbreviate(index, abb, abbreviator, quill) {
    if (!this.manuscriptEditor && queryManuscript(abb)) {
      quill.deleteText(index - abb.length, abb.length)
      this.prompt = abb
    } else {
      const caps = abb.toUpperCase() == abb
      const title = abb[0].toUpperCase() === abb[0]
      let match = this.cache.get(abb.toLowerCase())
      api
        .abbreviate(abb)
        .then((resp) => {
          if (resp.status == 208) {
            store.commit('incrementMissedAbb', {
              word: abb,
              abb: resp.data,
            })
          }
        })
        .catch(() => { })
      if (match) {
        let word = match

        if (match.charAt(0) == '.') {
          word = match.slice(1)
        } else {
          if (title) {
            word = match.charAt(0).toUpperCase() + match.slice(1)
          }
          if (caps && abb.length > 1) {
            word = match.toUpperCase()
          }
          if (match.charAt(0) == '_') {
            word = match.slice(1).toLowerCase()
          }
        }
        word.replaceAll("\u200B", "")
        EventBus.$emit('sendCC', word + abbreviator)
        this.insertAbbreviation(index, abb, abbreviator, word, quill)
        setTimeout(() => {
          quill.setSelection(quill.getSelection().index, 0)
          this.zeroWidthSpace = false
          db.useAbb({
            abb: abb.toLowerCase(),
            word: this.cache.get(abb.toLowerCase()),
          })
        }, 20)
        return word
      }

      EventBus.$emit('sendCC', abb + abbreviator)
      quill.insertText(index, abbreviator)
      return abb
    }
  }

  insertAbbreviation(index, abb, abbreviator, word, quill) {
    let format = quill.getFormat(index)
    if (abb == word || this.dontExpand) {
      quill.insertText(index, abbreviator, format)
      this.dontExpand = false
      return
    }
    let delta = new Delta()
      .retain(index - (this.zeroWidthSpace ? 1 : 0) - abb.length)
      .delete(abb.length + (this.zeroWidthSpace ? 1 : 0))
      .insert(word + abbreviator, format)
    quill.updateContents(delta)
  }

  scrollHandler(node) {
    if (!this.manuscriptEditor) {
      node.addEventListener('keyup', ({ key: key, target: { lastChild } }) => {
        this.lastKey = key
        if (this.scrollIntoView) {
          lastChild.scrollIntoView()
          this.scrollIntoView = false
        }
      })
    }
  }
  paste() {
    const text = e.clipboardData
      ? (e.originalEvent || e).clipboardData.getData('text/plain')
      : // For IE
      window.clipboardData
        ? window.clipboardData.getData('Text')
        : ''

    if (document.queryCommandSupported('insertText')) {
      // eslint-disable-line
      document.execCommand('insertText', false, text)
    } else {
      // Insert text at the current position of caret
      const range = document.getSelection().getRangeAt(0)
      range.deleteContents()

      const textNode = document.createTextNode(text)
      range.insertNode(textNode)
      range.selectNodeContents(textNode)
      range.collapse(false)

      const selection = window.getSelection()
      selection.removeAllRanges()
      selection.addRange(range)
    }
  }
  addKeybindings() {
    // TAB
    this.bindings[9].unshift({
      key: 9,
      handler: function(range) {
        if (this.lastKey == 'Tab') {
          this.quill.insertText(range.index, ' ')
        }
        this.dontExpand = true
      },
    })
    //ESCAPE
    this.addBinding({
      key: 27,
      handler: function(range, _context) {
        EventBus.$emit('closeNav', true)
        let end = this.quill.getText().length - 1
        if (this.lastKey == 'Escape') {
          this.prompt = ''
          removePreviews(range.index, this.quill)
        }
        this.quill.setSelection(end)
        return true
      },
    })

    //F1
    this.addBinding({
      key: 112,
      shiftKey: false,
      handler: function(range, context) {
        let phrase = ''
        if (range.length > 0) {
          phrase = this.quill.getText(range.index, range.length)
        } else {
          phrase = context.prefix.split(' ').pop()
        }
        store.commit('setLookupPhrase', phrase)
        EventBus.$emit('lookupPhrase', phrase)
        return true
      },
    })
    //F2
    this.addBinding({
      key: 113,
      shiftKey: false,
      handler: function(range, context) {
        let phrase = ''
        if (range.length > 0) {
          phrase = this.quill.getText(range.index, range.length)
        } else {
          phrase = context.prefix.split(' ').pop()
        }
        EventBus.$emit('addAbbreviation', phrase)
      },
    })
    //Shift+F2
    this.addBinding({
      key: 113,
      shiftKey: true,
      handler: function(range, context) {
        let phrase = ''
        if (range.length > 0) {
          phrase = this.quill.getText(range.index, range.length)
        } else {
          phrase = context.prefix.split(' ').pop()
        }
        EventBus.$emit('addAbbreviation', phrase.toLowerCase())
      },
    })
    //F3
    this.addBinding({
      key: 114,
      shiftKey: false,
      handler: function() {
        //EventBus.$emit('toggleCollab')
      },
    })
    //F4
    this.addBinding({
      key: 115,
      handler: function(range, _context) {
        removePreviews(range.index, this.quill)
        this.URL = false
        this.quill.setText('')
        this.capitalizeNext = true
        this.abbreviated = false
        EventBus.$emit('clear')
        this.prompt = ''
      },
    })
    //F6 sizeDown
    this.addBinding({
      key: 117,
      ctrlKey: false,
      handler: function() {
        EventBus.$emit('sizeChange', { inc: false, send: false })
      },
    })
    //ctrl+F6 send sizeDown
    this.addBinding({
      key: 117,
      ctrlKey: true,
      handler: function() {
        EventBus.$emit('sizeChange', { inc: false, send: true })
      },
    })
    //F7 sizeUp
    this.addBinding({
      key: 118,
      ctrlKey: false,
      handler: function() {
        EventBus.$emit('sizeChange', { inc: true, send: false })
      },
    })
    //CTRL+F7 send sizeUp
    this.addBinding({
      key: 118,
      ctrlKey: true,
      handler: function() {
        //console.log('keybinding ctrl+f7')
        EventBus.$emit('sizeChange', { inc: true, send: true })
      },
    })
    //F8 colorChange
    this.addBinding({
      key: 119,
      ctrlKey: false,
      handler: function() {
        EventBus.$emit('colorChange', false)
      },
    })
    //ctrl+F8 send colorChange
    this.addBinding({
      key: 119,
      ctrlKey: true,
      handler: function() {
        EventBus.$emit('colorChange', true)
      },
    })
    //F9 (Create)
    this.addBinding({
      key: 120,
      handler: function() {
        //console.log('keybinding F9')
        EventBus.$emit('createSession', true)
      },
    })
    //F10 (Join)
    this.addBinding({
      key: 121,
      handler: function() {
        EventBus.$emit('joinSession')
      },
    })

    //Space
    this.addBinding({
      key: 32,
      shiftKey: null,
      handler: function(range, context) {
        if (this.abbreviated) {
          this.abbreviated = false
          return true
        }
        let abb = this.wordBeforeCursor(context.prefix)
        //console.log(abb.length)
        //console.log("word before cursor:", abb)
        if (!abb) {
          EventBus.$emit('sendCC', ' ')
          this.abbreviated = false
          return true
        }
        //abb = this.capitalize(abb)
        this.abbreviate(range.index, abb, ' ', this.quill)
      },
    })
    //Ctrl+Space
    this.addBinding({
      key: 32,
      ctrlKey: true,
      handler: function(range, context) {
        let abb = this.wordBeforeCursor(context.prefix)
        //console.log(abb.length)
        //console.log("word before cursor:", abb)
        this.abbreviated = false
        this.capitalize = false
        this.dontExpandWord = this.abbreviate(range.index, abb, '', this.quill)
      },
    })
    //Enter
    this.bindings[13].unshift({
      key: 13,
      handler: function(range, context) {
        this.url = false
        let scroll = range.index == this.quill.getLength() - 1
        let abb = this.wordBeforeCursor(context.prefix)
        if (abb != '') {
          this.abbreviate(range.index, abb, '\n', this.quill)
        } else {
          EventBus.$emit('sendCC', '\n')
          this.quill.insertText(range.index, '\n')
        }
        if (this.options.capitalizeOnNewLine) {
          this.capitalizeNext = true
        } else {
          this.capitalizeNext = false
        }
        this.abbreviated = false
        this.scrollIntoView = scroll
        EventBus.$emit('newLine', scroll)
      },
    })

    //Shift+Enter
    this.bindings[13].unshift({
      key: 13,
      shiftKey: true,
      handler: function(range, context) {
        let abb = this.wordBeforeCursor(context.prefix)
        this.url = false
        if (abb) {
          this.abbreviate(range.index, abb, '\n', this.quill)
          this.capitalizeNext = true
        } else {
          EventBus.$emit('sendCC', '\n')
          this.quill.insertText(range.index, '\n')
        }
        if (this.options.capitalizeOnNewLine) {
          this.capitalizeNext = false
        } else {
          this.capitalizeNext = true
        }
        this.scrollIntoView = true
      },
    })

    //Backspace
    this.addBinding({
      key: 8,
      handler: function() {
        this.capitalizeNext = false
        return true
      },
    })

    this.addBinding({
      key: 8,
      ctrlKey: true,
      handler: function() {
        this.capitalizeNext = false
        return true
      },
    })

    //Period
    this.addBinding({
      key: 190,
      handler: function(range, context) {
        let abb = this.wordBeforeCursor(context.prefix)
        if (abb == '..' || !abb.trim()) {
          EventBus.$emit('sendCC', '.')
          return true
        }
        this.capitalizeNext = true
        this.abbreviated = true
        this.URL = true
        this.abbreviate(range.index, abb, '.', this.quill)
      },
    })
    //Colon
    this.addBinding({
      key: 190,
      shiftKey: true,
      handler: function(range, context) {
        this.capitalizeNext = true
        context.prefix.split(' ').map((w) => {
          if (w[0] == undefined) {
          } else if (w[0] === w[0].toLowerCase()) {
            this.capitalizeNext = false
          }
        })

        let abb = this.wordBeforeCursor(context.prefix)
        if (abb.length == 0) {
          EventBus.$emit('sendCC', ':')
          return true
        }
        this.abbreviate(range.index, abb, ':', this.quill)
      },
    })

    //Comma
    this.addBinding({
      key: 188,
      handler: function(range, context) {
        let abb = this.wordBeforeCursor(context.prefix)
        if (abb.length == 0) {
          EventBus.$emit('sendCC', ',')
          return true
        }
        this.abbreviate(range.index, abb, ',', this.quill)
      },
    })
    //Semicolon
    this.addBinding({
      key: 188,
      shiftKey: true,
      handler: function(range, context) {
        let abb = this.wordBeforeCursor(context.prefix)
        if (abb.length == 0) {
          EventBus.$emit('sendCC', ';')
          return true
        }
        this.abbreviate(range.index, abb, ';', this.quill)
      },
    })
    //Dash
    this.addBinding({
      key: 189,
      handler: function(range, context) {
        let abb = this.wordBeforeCursor(context.prefix)
        if (abb.length == 0) {
          EventBus.$emit('sendCC', '-')
          return true
        }
        this.abbreviate(range.index, abb, '-', this.quill)
      },
    })
    //Exclamation
    this.addBinding({
      key: 49,
      shiftKey: true,
      handler: function(range, context) {
        let abb = this.wordBeforeCursor(context.prefix)
        if (abb.length == 0) {
          EventBus.$emit('sendCC', '!')
          return true
        }
        this.abbreviate(range.index, abb, '!', this.quill)
        this.capitalizeNext = true
        this.abbreviated = true
      },
    })
    //Question mark
    this.addBinding({
      key: 187,
      shiftKey: true,
      handler: function(range, context) {
        let abb = this.wordBeforeCursor(context.prefix)
        if (abb.length == 0) {
          EventBus.$emit('sendCC', '?')
          return true
        }
        this.abbreviate(range.index, abb, '?', this.quill)
        this.capitalizeNext = true
        this.abbreviated = true
      },
    })

    this.addBinding({
      key: 57,
      shiftKey: true,
      handler: function(range, context) {
        let abb = this.wordBeforeCursor(context.prefix)
        if (abb.length == 0) {
          EventBus.$emit('sendCC', ')')
          return true
        }
        this.abbreviate(range.index, abb, ')', this.quill)
        this.capitalizeNext = false
        this.abbreviated = true
      },
    })
    // Forwardslash
    this.addBinding({
      key: 55,
      shiftKey: true,
      handler: function(range, context) {
        let abb = this.wordBeforeCursor(context.prefix)
        if (abb.length == 0) {
          EventBus.$emit('sendCC', '/')
          return true
        }
        this.abbreviate(range.index, abb, '/', this.quill)
        this.capitalizeNext = false
        this.abbreviated = false
      },
    })
    // Double quotation mark "
    this.addBinding({
      key: 50,
      shiftKey: true,
      handler: function(range, context) {
        let abb = this.wordBeforeCursor(context.prefix)
        if (abb.length == 0) {
          EventBus.$emit('sendCC', '"')
          return true
        }
        this.abbreviate(range.index, abb, '"', this.quill)
        this.capitalizeNext = false
        this.abbreviated = true
      },
    })
    // Quotation mark '
    this.addBinding({
      key: 191,
      handler: function(range, context) {
        //console.log("' quotation mark")
        let abb = this.wordBeforeCursor(context.prefix)
        if (abb.length == 0) {
          EventBus.$emit('sendCC', "'")
          return true
        }
        this.abbreviate(range.index, abb, "'", this.quill)
        this.capitalizeNext = false
        this.abbreviated = true
      },
    })
    //Ctrl+Backspace
    this.bindings[8].unshift({
      key: 8,
      ctrlKey: true,
      handler: function(_range, context) {
        //console.log('prefix:', context.prefix, 'length.', context.prefix.length)
        if (context.prefix.trim().split(' ').length == 1) {
          this.capitalizeNext = true
        }
        return true
      },
    })
    // Right →
    this.bindings[39].unshift({
      key: 39,
      handler: function() {
        if (this.prompt != '') {
          let index = this.quill.getText().length - 1
          let offset = promptManuscript(index, this.quill, 'insertWord')

          if (offset == -1) {
            this.prompt = false
            return false
          }
          this.quill.setSelection(index + offset)
          return false
        }
        return true
      },
    })
    // Ctrl+Right →
    this.bindings[39].unshift({
      key: 39,
      ctrlKey: true,
      handler: function(range, _context) {
        if (range.index == this.quill.getText().length - 1) return false
        return true
      },
    })
    // ALT+Right →

    this.bindings[39].unshift({
      key: 39,
      altKey: true,
      handler: function() {
        if (this.prompt != '') {
          let index = this.quill.getText().length - 1
          let offset = promptManuscript(index, this.quill, 'skipWord')

          if (offset == -1) {
            this.prompt = false
            return false
          }
          this.quill.setSelection(index + offset)
          return false
        }
        return true
      },
    })

    // Down ↓
    this.addBinding({
      key: 40,
      handler: function() {
        if (this.prompt != '') {
          let index = this.quill.getText().length - 1
          let offset = promptManuscript(index, this.quill, 'insertClause')

          if (offset == -1) {
            this.prompt = false
            return false
          }
          this.quill.setSelection(index + offset)
          return false
        }
        return true
      },
    })
    // ALT+Down ↓
    this.addBinding({
      key: 40,
      altKey: true,
      handler: function() {
        if (this.prompt != '') {
          let index = this.quill.getText().length - 1
          let offset = promptManuscript(index, this.quill, 'skipClause')

          if (offset == -1) {
            this.prompt = false
            return false
          }
          this.quill.setSelection(index + offset)
          return false
        }
        return true
      },
    })
    //Left
    this.bindings[37].unshift({
      key: 37,
      handler: function(_range, context) {
        if (this.prompt != '') {
          let index = this.quill.getText().length - 1
          let offset = promptManuscript(
            index,
            this.quill,
            'removeWord',
            context
          )
          this.quill.setSelection(offset)
          return false
        }
        return true
      },
    })
    this.addBinding({
      key: 86,
      ctrlKey: true,
      handler: function() {
        return this.manuscriptEditor ? true : false
      },
    })
    // CTRL+K
    this.addBinding({
      key: 75,
      ctrlKey: true,
      handler: function(range, _context) {
        this.quill.insertEmbed(range.index, 'protype', 'A')
        this.quill.setSelection(range.index + 1)
        return false
      },
    })
    //CTRL+1
    this.addBinding({
      key: 49,
      ctrlKey: true,
      handler: function() {
        if (store.getters.getModalOpen == false) {
          EventBus.$emit('changeStandardList', 1)
          return false
        }
      },
    })

    //CTRL+2
    this.addBinding({
      key: 50,
      ctrlKey: true,
      handler: function() {
        if (store.getters.getModalOpen == false) {
          EventBus.$emit('changeStandardList', 2)
          return false
        }
      },
    })

    //CTRL+3
    this.addBinding({
      key: 51,
      ctrlKey: true,
      handler: function() {
        if (store.getters.getModalOpen == false) {
          EventBus.$emit('changeStandardList', 3)
          return false
        }
      },
    })

    //CTRL+4
    this.addBinding({
      key: 52,
      ctrlKey: true,
      handler: function() {
        if (store.getters.getModalOpen == false) {
          EventBus.$emit('changeStandardList', 4)
          return false
        }
      },
    })

    //CTRL+5
    this.addBinding({
      key: 53,
      ctrlKey: true,
      handler: function() {
        if (!store.getters.getModalOpen) {
          EventBus.$emit('changeStandardList', 5)
          return false
        }
      },
    })
  }
}

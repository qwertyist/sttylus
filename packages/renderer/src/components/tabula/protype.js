import Quill from 'quill'

let Embed = Quill.import('blots/embed')
class Protype extends Embed {
  static create(value) {
    let node = super.create(value)
    node.innerHTML = `${value}`
    node.setAttribute('data-proc', value)
    return node
  }
  static value(domNode) {
    return domNode.getAttribute('data-proc')
  }
}

Protype.blotName = 'protype'
Protype.className = 'protype'
Protype.tagName = 'span'

export default Protype

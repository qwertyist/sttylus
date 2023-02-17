import Quill from "quill"

let Embed = Quill.import("blots/embed")
class ZoomMessage extends Embed {
    static create(value) {
        let node = super.create(value)
        node.innerHTML = `${value}`
        node.setAttribute("data-proc", value)
        return node
    }
    static value(domNode) {
        return domNode.getAttribute("data-proc") 
    }
}

ZoomMessage.blotName = "ZoomMessage"
ZoomMessage.className = "ZoomMessage"
ZoomMessage.tagName = "span"

export default ZoomMessage;

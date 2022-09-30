import Quill from "quill"

let Embed = Quill.import("blots/embed")
class Preview extends Embed {
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

Preview.blotName = "preview"
Preview.className = "preRender"
Preview.tagName = "span"

export default Preview;
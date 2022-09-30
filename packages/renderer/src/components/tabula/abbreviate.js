import api from "../../api/api";

let cache = []
/*export function abbreviate(abb) {
    return api.abbreviate(JSON.stringify(abb))
}*/

export function abbreviate(index, abb, abbreviator, quill) {
    console.log("i:", index)
    console.log("abb:", abb)
    console.log("abbreviator:", abbreviator)
    api.abbreviate(abb)
        .then(r => {
            let word = r.data.word;
            let format = quill.getFormat(index)
            let deleted = quill.deleteText(index - abb.length, abb.length, "user")
            if (Object.keys(format).length !== 0) {
                quill.insertText(index - abb.length, word + abbreviator, format);
            } else {
                quill.insertText(index - abb.length, word + abbreviator, { bold: false, underscore: false, italic: false })
            }
        })

}
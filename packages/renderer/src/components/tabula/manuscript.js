import  Delta  from "quill-delta";
import api from "../../api/api";
let manuscripts = new Map();
let default_manuscript = {
    deletes: [],
    line: 0,
    op: 0,
    word: 0,
    prompt: false,
    lookahead: false,
    empty: false,
    text: [],
    preview: true
}

let manuscript = { ...default_manuscript }

const no_attrs = {
    bold: false,
    underline: false,
    italic: false,
}

var fontColor = "#000000"

export function loadManuscripts(docs) {
    manuscripts = new Map();
    docs.map(d => {
        manuscripts.set(d.abb, d.id)
    })
}

export function queryManuscript(abb) {
    manuscript.deltas = []
    if (manuscript.prompt) { return false }
    let id = manuscripts.get(abb.toLowerCase())
    if (id) {
        return api.getManuscript(id).then(resp => {
            let delta = new Delta(JSON.parse(resp.data.content))
            console.log("recv manuscript", resp.data.content)
            let content = []
            delta.eachLine((ops, attr) => {
                content.push(ops)
                let preview = ""
                ops.map(op => {
                    preview += op.insert
                })
                manuscript.text.push(preview)
            })
            manuscript.content = content
            manuscript.line = 0
            manuscript.word = 0
            manuscript.op = 0
            manuscript.prompt = true

            manuscript.preview = true
            return true
        })
    } else { return false }
}

function renderPreview(index, q, action = "paragraph") {
    let text = []
    switch (action) {
        case "paragraph":
            console.log("insert preview of paragraph")
            text = manuscript.text[manuscript.line].split(" ").slice(1);
            text.map((t, i) => {
                //console.log("preview index:", index + i, "t:", t)
                q.insertEmbed(index + i, "preview", t + " ")
            })
            break
        case "word":
            let word = manuscript.text[manuscript.line].split(" ")[manuscript.word]
            console.log(manuscript.text[manuscript.line])
            console.log("insert preview of word:", word)
            q.insertEmbed(index, "preview", word + " ")
            break
    }
}

export function removePreviews(index, q) {
    const length = q.getLength() - 1
    q.removeFormat(index, length)
    manuscript = { ...default_manuscript }
    manuscript.text = []
    console.log("after removing previews:", manuscript)
}

export function promptManuscript(index, q, action = "insertWord", context) {
    let skip = false
    switch (action) {
        case "skipWord":
            console.log("### ACTION - SKIP WORD")
            action = "insertWord"
            skip = true
            break;
        case "skipClause":
            console.log("### ACTION - SKIP CLAUSE")
            action = "insertClause"
            skip = true
            break;
    }
    if (manuscript.line == manuscript.content.length) {
        manuscript.prompt = false
    }

    let offset = 0
    if (!manuscript.prompt) {
        console.log("Got to end")
        removePreviews(index, q)
        return -1
    }

    if (manuscript.empty && action == "insertClause") {
        manuscript.op++
        manuscript.word = 0
        manuscript.empty = false
    }

    let newLine = false
    let curr_line = manuscript.content[manuscript.line]
    let line_attributes = {}
    if ("attributes" in curr_line) {
        line_attributes = curr_line.attributes || {}
    }
    let curr_op = curr_line.ops[manuscript.op] || null
    let op_attributes = {}
    let attr = {}
    let insert = ""
    if (curr_op) {
        if (curr_line.ops.length > 1 && manuscript.op < curr_line.ops.length) {
            manuscript.lookahead = true
        } else {
            manuscript.lookahead = false
        }
        op_attributes = curr_op.attributes || no_attrs
        insert = curr_op.insert
    }
    /*
    if (manuscript.op == curr_line.ops.length - 1 && manuscript.word == insert.split(" ").length) {
        console.log("go to next line?!")
        manuscript.line = 0
        manuscript.op = 0
        manuscript.word = 0
        return promptManuscript(index + 1, q, action)
    }
    /*  
        console.log("line:", manuscript.line)
        console.log("curr_line:", curr_line)
        console.log("line_attributes:", line_attributes)
        console.log("op_attributes:", op_attributes)
        console.log("insert:",insert)
        console.log("words:", insert.split(" "))
        */
    switch (action) {
        case "insertEdgeCase":
            q.deleteText(index, 1)
            q.insertText()
        case "insertWord":
            if (!skip) {
                console.log("### ACTION - INSERT WORD")
            }
            let chunk = insert.split(" ")[manuscript.word]
            console.log("insert chunk:", chunk)
            manuscript.empty = false
            if (manuscript.word > 0 || insert.split(" ").length == 1) {
                q.deleteText(index, 1)
            }
            if (manuscript.line > 0 && manuscript.word == 0 && manuscript.op == 0) {
                console.log("Add one to offset, because new line, first word, first op")
                offset += 1
                newLine = true
            }
            attr = { ...line_attributes, ...op_attributes, color: null }

            if (!skip) {
                q.insertText((newLine ? 1 : 0) + index, chunk, attr)
                offset += chunk.length + 1;
            } else {
            }
            console.log("word # in op", manuscript.word)
            console.log("insert:", insert.split(" "))
            if (manuscript.word + 1 < insert.split(" ").length) {
                console.log("still words left in operation before formatting")
                if (!skip) {
                    q.insertText((newLine ? 1 : 0) + index + chunk.length, " ", attr)
                }
                manuscript.word++;
                if (insert.split(" ")[manuscript.word] == "") {
                    console.log("handle empty chunk")
                    manuscript.empty = true
                }
            } else {
                if (manuscript.op < curr_line.ops.length - 1) {
                    console.log("go to next op")
                    manuscript.op++;
                    manuscript.word = 0;
                    return offset + promptManuscript(index + offset - 1, q, action)
                } else {
                    console.log("go to next line")
                    manuscript.line++
                    manuscript.word = 0
                    manuscript.op = 0
                    manuscript.preview = true
                }
            }
            break;
        case "insertClause":
            if (!skip) {
                console.log("### ACTION - INSERT CLAUSE")
            }
            if (manuscript.word == 0 && manuscript.line > 0 && manuscript.op == 0) {
                q.insertText(index, "\n")
                newLine = true
                offset++
            }
            let paragraph = {}
            let clauses = []
            let clause = []
            attr = { ...line_attributes, ...op_attributes, color: null}
            if (curr_op) {
                paragraph = curr_op.insert.split(" ").slice(manuscript.word).join(" ") || curr_op.insert
                clauses = paragraph.split(/(?<=\,|\.|\!|\?)/);
                clause = clauses[0].split(" ");
                clause[clause.length - 1] += " "
            } else {
                manuscript.line++;
                manuscript.op = 0;
                manuscript.word = 0;
                manuscript.lookahead = false;
                return offset
            }
            for (let i = 0; i < clause.length; i++) {
                if (!newLine) {
                    q.deleteText(index + offset, 1)
                } else {
                    newLine = false
                }
                if (!skip) {
                    q.insertText(index + offset, clause[i] + (i < clause.length - 1 ? " " : ""), attr)
                    offset += clause[i].length + (i < clause.length - 1 ? 1 : 0)
                }
            }
            manuscript.word += clause.length
            if (manuscript.word + 1 > curr_op.insert.split(" ").length) {
                if (manuscript.lookahead) {
                    console.log("should look ahead")
                    manuscript.op++;
                    manuscript.word = 0;
                    var matches = paragraph.match(/\,|\.|!|\?/)
                    if (matches == null) {
                        console.log("Look ahead")
                        return offset + promptManuscript(index + offset - 1, q, action)
                    } else {
                        console.log("Dont look ahead")
                        manuscript.line++
                        manuscript.word = 0
                        manuscript.op = 0
                        manuscript.preview = true
                        manuscript.lookahead = false
                        break;
                    }
                }

                manuscript.word = 0
                manuscript.line++
                manuscript.preview = true
            }

            break
        case "removeWord":
            console.log("### ACTION - REMOVE WORD")
            break
            let last_word = context.prefix.split(" ").splice(-2, 1)[0]
            console.log("prefix:", context.prefix.split(" "))
            console.log("last word in context:", last_word)
            const length = q.getLength()
            q.deleteText(index - last_word.length - 1, last_word.length + 1)
            q.removeFormat(index - last_word.length - 1, length)
            offset = index - last_word.length - 1
            renderPreview(offset, q, "paragraph")
            manuscript.word--
            if (manuscript.word < 0) {
                if (manuscript.op > 0) {
                    manuscript.op--
                    let curr_op = curr_line.ops[manuscript.op] || null
                    manuscript.word == curr_op.insert.split(" ").length - 1
                }
            }
            return offset

    }

    //    console.log("index and offset:", index + offset)
    console.log("line:", manuscript.line, "op:", manuscript.op, "word:", manuscript.word)
    if (manuscript.preview && manuscript.word == 1) {
        renderPreview(index + offset, q, "paragraph")
        manuscript.preview = false
    }
    return offset
}

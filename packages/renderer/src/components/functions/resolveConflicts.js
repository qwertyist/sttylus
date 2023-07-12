var regex = /([\p{L}\p{N}\s]*)=(.*)/iu
var duplicates = []

export function resetValidState(abbs) {
    return abbs.map((a) => {
        console
        a.validabb = null
        a.validword = null
        a.issue = ''
        a.notduplicate = null
        return a
    })
}

export function getNumberOfConflicts(abbs) {
    return abbs.filter((a) => a.issue != null).length
}

function findDuplicates(abbs) {
    const lookup = abbs.reduce((a, e) => {
        a[e.abb] = ++a[e.abb] || 0
        return a
    }, {})
    duplicates = abbs.filter((a) => lookup[a.abb] && a.abb != '')
    if (duplicates.length > 0) {
        duplicates.map((dup) => {
            abbs.map((a, i) => {
                if (dup.abb == a.abb) {
                    a.notduplicate = false
                    a.validabb = false
                    a.issue = '__importduplicate__'
                    console.log('found duplicate - overwrite?', a.overwrite)
                    console.log(a.abb)
                } else {
                    a.notduplicate = null
                }
                abbs[i] = { ...a }
            })
        })
    } else {
        abbs = abbs.map((a) => {
            a.notduplicate = null
            return a
        })
    }

    return abbs
}
export function resolveIssues(abbs) {
    let counter = 0

    var resolved = abbs.map((a) => {
        if (a.overwrite == undefined) {
            a.overwrite = false
        }
        if (a.abb == a.word) {
            a.validabb = false
            a.validword = false
            a.issue = 'Fälten är identiska'
            return a
        }
        if (a.validabb == false && a.validword == false) {
            if (a.issue == '__formaterror__') {
                if (a.abb != '') {
                    a.validabb = null
                    a.issue = null
                } else {
                    return a
                }
            } else {
                if (a.abb == a.word) {
                    a.issue = 'Fel format vid importering'
                    counter++
                    return a
                } else {
                    a.validabb = null
                    a.validword = null
                    a.issue = null
                }
            }
        }
        if (a.abb == '') {
            a.validabb = false
            a.issue = 'Förkortning saknas'
            counter++
            return a
        } else {
            a.issue = null
            a.validabb = null
        }
        if (a.word == '') {
            a.validword = false
            a.issue = 'Text/fras saknas'
            counter++
            return a
        } else {
            a.validword = null
            a.issue = null
        }
        return a
    })
    resolved = findDuplicates(resolved)
    return { abbs: resolved, issues: counter }
}

export function parseTxt(text) {
    let seen = []
    let duplicates = []
    let lines = text.split('\n')
    let parsed = lines
        .filter((a) => {
            if (a.trim() == '') {
                return false
            } else {
                return true
            }
        })
        .map((a, i) => {
            let fields = regex.exec(a)
            // Find single word or badly formatted lines
            if (fields == null) {
                return {
                    abb: '',
                    word: a.trim(),
                    notduplicate: null,
                    validabb: false,
                    validword: false,
                    issue: '__formaterror__',
                }
            }
            // Find empty abbreviations
            if (fields[1] == '') {
                return {
                    abb: '',
                    word: fields[2].trim(),
                    validabb: false,
                    validword: null,
                    notduplicate: null,
                    issue: 'Förkortning saknas',
                }
            }
            // Find empty abbreviations
            let abb = fields[1]
                .trim()
                .replaceAll(/[\p{P}\p{Z}]/gu, '')
                .toLowerCase()
            let obj = {
                abb: abb,
                word: fields[2].trim(),
                notduplicate: null,
                validword: null,
                validabb: null,
                issue: null,
            }
            if (seen.indexOf(abb) == -1) {
                seen.push(abb)
                obj.notduplicate = null
            } else {
                if (duplicates.indexOf(abb) == -1) {
                    duplicates.push(abb)
                }
                obj.notduplicate = false
                obj.issue = '__textduplicate'
                obj.validabb = false
                obj.validword = null
                obj.validword = null
            }
            if (fields[2] == '') {
                obj.validword = false
                obj.issue = 'Text/fras saknas'
                obj.validabb = null
                obj.notduplicate = null
            }
            return obj
        })
    duplicates.map((d) => {
        parsed = parsed.map((a) => {
            if (a.abb == d && a.notduplicate == null) {
                a.notduplicate = false
                a.issue = '__textduplicate__'
                a.validabb = false
            }
            return a
        })
    })
    parsed = parsed.map((a) => {
        a.comment = JSON.stringify({
            notduplicate: a.notduplicate,
            validabb: a.validabb,
            validword: a.validword,
            issue: a.issue,
        })
        return a
    })
    return parsed
}

function findWordBeforeCursor({ anchor = 0 }) {
    return ($position) => {
        if ($position.depth <= 0) {
            return false
        }

        const textFrom = $position.before()
        const text = $position.doc.textBetween(textFrom, anchor, '\0', '\0')
        const wordsToAnchor = text.split(' ')
        const query = wordsToAnchor[wordsToAnchor.length - 1]
        let position
        position = {
            range: {
                from: anchor - query.length,
                to: anchor,
            },
            query: query,
            text: text,
        }
        return position
    }
}

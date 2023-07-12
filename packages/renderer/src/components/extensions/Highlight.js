export default class HighlightMark extends Mark {
    get name() {
        return 'mark'
    }

    get schema() {
        return {
            attrs: {
                color: {
                    default: 'rgba(240,87,100,0.7)',
                },
            },
            parseDOM: [
                {
                    tag: 'mark',
                },
            ],
            toDOM: (mark) => [
                'mark',
                {
                    style: `background:${mark.attrs.color}`,
                },
                0,
            ],
        }
    }

    commands({ type }) {
        return (attrs) => toggleMark(type, attrs)
    }
}

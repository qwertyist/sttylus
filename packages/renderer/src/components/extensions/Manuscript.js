import { Mark } from 'tiptap';

export default class Manuscript extends Mark {
  constructor (options = {}) {
    super(options);
  }

  // eslint-disable-next-line class-methods-use-this
  get name () {
    return 'manuscript';
  }

  // eslint-disable-next-line class-methods-use-this
  get defaultOptions () {
    return {
      manuscriptClass: 'manuscript'
    };
  }

  // eslint-disable-next-line class-methods-use-this
  get schema () {
    return {
      attrs: {
        id: {},
        text: {},
        newLine: false
      },
      group: 'inline',
      inline: true,
      selectable: false,
      atom: true,
      toDOM: node => [
        'span',
        {
          class: this.options.manuscriptClass,
          'manuscript-node-id': node.attrs.id,
          'manuscript-newLine': node.attrs.newLine
        },
        `${node.attrs.text}`
      ],
      parseDOM: [
        {
          tag: 'span[manuscript-node-id]',
          getAttrs: dom => {
            const id = dom.getAttribute('manuscript-node-id');
            const text = dom.innerText;
            const newLine = dom.getAttribute('manuscript-newLine');
            return { id, text, newLine };
          }
        }
      ]
    };
  }
}
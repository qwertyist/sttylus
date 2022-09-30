/* eslint-disable class-methods-use-this */
import { Extension } from 'tiptap';
import KeyboardShredder from '../plugins/KeyboardShredder';

export default class ExpandWord extends Extension {
  constructor (options = {}) {
    super(options);
  }

  get name () {
    return 'expandWord';
  }

  get defaultOptions () {
    return {
      prompt: false
    };
  }

  get update () {
    return view => {
      view.updateState(view.state);
    };
  }
  get schema () {
    return {
      attrs: {
        id: {},
        label: {}
      },
      group: 'inline',
      inline: true,
      selectable: false,
      atom: true
    };
  }

  get plugins () {
    return [
      KeyboardShredder({
        command: ({ query, range }) => this.expand(query, range),
        appendText: ' ',
        userId: this.options.userId,
        capitalizeOnNewLine: this.options.capitalizeOnNewLine,
        manuscripts: this.options.manuscripts,
        onChange: this.options.onChange,
        onKeyDown: this.options.onKeyDown,
        prompt: this.options.prompt
      })
    ];
  }

}
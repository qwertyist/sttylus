/* eslint-disable complexity */
/* eslint-disable max-depth */
import { Plugin, PluginKey } from 'prosemirror-state';
import axios from 'axios';
import api from "../../api/api"
import { TextSelection } from 'tiptap';
import EventBus from "../../eventbus.js"
import { store } from "../../store/index"

const backend = import.meta.env.VITE_STTYLUS_BACKEND

let querying = false
let prompt = false;
let promptOnNewLine = false;
let manuscript = {
  abb: '',
  content: '',
  line: 0,
  index: 0,
  beginning: false,
};

function createManuscriptMark(schema) {
  return schema.marks.manuscript.create({
    id: (Math.random() + 1).toString(36).substr(2, 5),
    class: 'manuscript',
    text: '',
    newLine: false,
  });
}

function findWordBeforeCursor({ anchor = 0 }) {
  return $position => {
    if ($position.depth <= 0) {
      return false;
    }

    const textFrom = $position.before();
    const text = $position.doc.textBetween(textFrom, anchor, '\0', '\0');


    const wordsToAnchor = text.split(/[ \"\':;\/\-()[\]_{}\uFEFF]+/);
    const query = wordsToAnchor[wordsToAnchor.length - 1];
    let position;
    position = {
      range: {
        from: anchor - query.length,
        to: anchor,
      },
      query,
      text,
    };
    return position;
  };
}

function queryAbbreviator(query) {
  return api.abbreviate(query)
}

function queryManuscript(abb, id) {
  axios({
    method: 'GET',
    url: backend + '/api/docs/' + id,
  }).then(response => {
    /* let parser = new DOMParser();
     const htmldoc = parser.parseFromString(response.data.content, 'text/html');
     const paragraphs = htmldoc.getElementsByTagName('p');
     let wordsPerParagraph = [];
     for (let p of paragraphs) {
       const lines = p.innerHTML.split('<br>');
       for (let l of lines) {
         wordsPerParagraph.push(l.split(/ |\r/));
       }
     }
     */
    let parsed = JSON.parse(response.data.content)
    let content = parsed.content
    let paragraph = []
    let wordsPerParagraph = [];
    //console.log("reading content from parsed", content)
    content.forEach((paragraph, line) => {
      //console.log("paragraph:", paragraph)
      if (paragraph.content !== undefined) {
        wordsPerParagraph.push(paragraph.content[0].text.split(/ /))
      } else {
        wordsPerParagraph.push(" ")
      }
    })
    /*for(let paragraph in content) {
       console.log("paragraph:", paragraph)
       for(let word in paragraph.content[0]) {
         wordsPerParagraph.push(word)
       }
     }*/
    //console.log(wordsPerParagraph)

    manuscript = {
      abb,
      content: wordsPerParagraph,
      line: 0,
      index: 0,
      beginning: true,
    };
  });
}

function traverseManuscript(action) {
  if (manuscript.line + 1 > manuscript.content.length) {
    console.log('stop prompting');
    prompt = false;
    return false;
  }
  let newLine = false;
  let chunk = [];
  let chunks = [];
  let text = ""
  switch (action) {
    case 'promptWord':
      if (promptOnNewLine || (manuscript.index == 0 && manuscript.line > 0)) {
        newLine = true;
        promptOnNewLine = false;
      } else {
        newLine = false;
      }
      //console.log("Prompt word:", manuscript.line, manuscript.index)
      if (manuscript.index == manuscript.content[manuscript.line].length) {
        manuscript.line += 1;
        manuscript.index = 0;
        newLine = true;
      }

      if (manuscript.content[manuscript.line] == undefined) {
        prompt = false
        console.log("stop prompting")
        return
      }

      text = manuscript.content[manuscript.line][manuscript.index] + ' ';
      manuscript.index += 1;


      chunk = [
        {
          id: (Math.random() + 1).toString(36).substr(2, 5),
          class: 'manuscript',
          text,
          newLine,
        },
      ];

      return chunk;
    case 'removeWord':
      if (manuscript.index == 1 && manuscript.line == 0) {
        manuscript.index -= 1;

        manuscript.beginning = true;
        return;
      } else if (manuscript.index == 0) {
        if (manuscript.line == 0) {
          manuscript.beginning = true;
          return;
        }
        manuscript.line -= 1;
        manuscript.index = manuscript.content[manuscript.line].length - 1;
        return;
      }
      manuscript.index -= 1;
      break;
    case 'deterWord':
    case 'promptParagraph':
      if (manuscript.index > 0) {
        newLine = false;
        text = manuscript.content[manuscript.line].slice(manuscript.index);
        manuscript.index = 0;
      } else {
        newLine = true;
        text = manuscript.content[manuscript.line];
        manuscript.index = 0;
      }
      manuscript.line += 1;
      if (manuscript.line > manuscript.content.length - 1) {
        console.log("Line exceeding content")
        prompt = false;
      }
      console.log("text:", text)
      for (let i = 0; i < text.length; i++) {
        if (i > 0) {
          newLine = false;
        }
        chunks.push({
          id: (Math.random() + 1).toString(36).substr(2, 5),
          class: 'manuscript',
          text: text[i] + ' ',
          newLine,
        });
      }
      return chunks;
    case 'removeParagraph':
      if (manuscript.line == 0) {
        manuscript.line = 0;
        manuscript.index = 0;
        return;
      }
      if (manuscript.index > 0) {
        manuscript.index = 0;
        return;
      }

      manuscript.line -= 1;
      manuscript.index = 0;
      break;
    case 'deterParagraph':
    case 'promptClause':
      if (promptOnNewLine || (manuscript.index == 0 && manuscript.line > 0)) {
        newLine = true;
        promptOnNewLine = false;
      } else {
        newLine = false;
      }
      if (manuscript.index == manuscript.content[manuscript.line].length) {
        newLine = true;
        manuscript.index = 0;
        manuscript.line += 1;
        if (manuscript.line == manuscript.content.length) {
          prompt = false
          console.log("stop prompting")
          return
        }
      }
      let paragraph = " "
      let clause = " "
      if (manuscript.content[manuscript.line].length > 1) {
        paragraph = manuscript.content[manuscript.line]
          .slice(manuscript.index)
          .join(' ');
        let clauses = paragraph.split(/(?<=[\,|\.|\!|\?]+)/);
        clause = clauses[0].split(" ");
        /*        console.log(clauses.slice(1).length) 
                console.log("remaining clauses:", clauses.slice(1))
                console.log(clauses.slice(1).join("").length)
                */
        if (clauses.slice(1).length == clauses.slice(1).join("").length) {
          clause[clause.length - 1] += clauses.slice(1).join("")
        }
        manuscript.index += clause.length;
      } else {
        clause = manuscript.content[manuscript.line]
        manuscript.index++
      }
      for (let i = 0; i < clause.length; i++) {
        if (i > 0) {
          newLine = false;
        }
        chunks.push({
          id: (Math.random() + 1).toString(36).substr(2, 5),
          class: 'manuscript',
          text: clause[i] + ' ',
          newLine,
        });
      }
      return chunks;
    case 'removeClause':
    case 'deterClause':
  }
}

function promptManuscriptNodes(view, chunks, clause) {
  const { doc, schema } = view.state;
  const content = doc.content;
  let offset = 0;
  //console.log(chunks)
  if (!chunks) {
    return
  }
  for (let i = 0; i < chunks.length; i++) {
    if (chunks[i].newLine) {
      const newLine = schema.text(chunks[i].text, [
        schema.marks.manuscript.create(chunks[i]),
      ]);
      view.dispatch(view.state.tr.insert(content.size, newLine));
      // view.dispatch(view.state.tr.addMark(content.size, content.size, schema.marks.manuscript.create(chunks[i])))
      // view.dispatch(view.state.tr.removeStoredMark(schema.marks.manuscript))
      offset += 2;
    } else {
      view.dispatch(
        view.state.tr.insertText(chunks[i].text, content.size - 1 + offset)
      );
      view.dispatch(
        view.state.tr.addMark(
          content.size - 1 + offset,
          content.size + offset + chunks[i].text.length,
          schema.marks.manuscript.create(chunks[i])
        )
      );
      // view.dispatch(view.state.tr.removeStoredMark(schema.marks.manuscript))
    }

    view.dispatch(
      view.state.tr.setSelection(view.state.tr.selection).scrollIntoView()
    );

    offset += chunks[i].text.length;
  }
}

function removeManuscriptNodes(view, removeLine) {
  const { doc } = view.state;
  let nodes = [];
  let lineStart = 0;
  doc.nodesBetween(0, doc.content.size - 1, (node, startpos) => {
    node.marks.forEach(mark => {
      if (mark.attrs.newLine == true) {
        if (removeLine) {
          lineStart = startpos - 3;
        }
        nodes.push([startpos - 2, node.text.length + 2]);
      } else {
        nodes.push([startpos, node.text.length]);
      }
    });
  });
  if (nodes.length == 0) {
    prompt = false;
    return;
  }
  let from = 0;
  let to = 0;
  if (removeLine) {
    from = lineStart;
    to = nodes[nodes.length - 1][0] - 1 + nodes[nodes.length - 1][1];
    view.dispatch(view.state.tr.deleteRange(from, to));
  } else {
    from = nodes[nodes.length - 1][0];
    to = nodes[nodes.length - 1][1];
    view.dispatch(view.state.tr.deleteRange(from, from + to));
  }
}

export default function KeyboardShredder({
  matcher = {
    char: '@',
    allowSpaces: false,
    startOfLine: false,
  },
  chars = [
    'Enter',
    '\n',
    ' ',
    '.',
    ',',
    '-',
    ':',
    ';',
    '!',
    '?',
    '/',
    '\\',
    '(',
    ')',
    '[',
    ']',
    '{',
    '}',
    '"',
    '\uFEFF',
  ],
  terminalChars = ['. ', '! ', '? '],
  capitalizeOnNewLine = true,
  hardbreak = false,
  manuscriptClass = 'manuscript',
  onChange = () => false,
  onKeyDown = () => false,
  directions = ['ArrowRight', 'ArrowLeft', 'ArrowDown', 'ArrowUp', 'Escape'],
  userId = '',
  manuscripts = {},
}) {
  return new Plugin({
    key: new PluginKey('suggestions'),

    view() {
      return {
        update: (view, prevState) => {
          const prev = this.key.getState(prevState);
          const next = this.key.getState(view.state);

          // See how the state changed
          const moved =
            prev.active && next.active && prev.range.from !== next.range.from;
          const started = !prev.active && next.active;
          const stopped = prev.active && !next.active;
          const changed = !started && !stopped && prev.query !== next.query;
          const handleChange = changed && !moved;
          const handleExit = stopped || moved;

          const state = handleExit ? prev : next;

          const props = {
            view,
            range: state.range,
            query: state.query,
            text: state.text,
            manuscript: state.manuscript,
            prompt: state.prompt,
            append: state.append,
            capitalize: state.capitalize,
          };

          if (handleChange) {
            onChange(props);
          }
        },
      };
    },

    state: {
      // Initialize the plugin's internal state.
      init() {
        EventBus.$on("setSelectedManuscripts", () => {
          manuscripts = store.getters.selectedManuscripts
        })
        return {
          active: false,
          range: {},
          query: null,
          text: null,
          abbreviate: false,
          manuscript: false,
          manuscripts: [],
          capitalize: true,
          step: 0,
        };
      },

      // Apply changes to the plugin state from a view transaction.
      apply(tr, prev) {
        const { selection } = tr;
        const next = { ...prev };
        if (tr.steps.length === 0) {
          return prev;
        }
        if (selection.from === selection.to) {
          if (
            selection.from < prev.range.from ||
            selection.from > prev.range.to
          ) {
            next.active = false;
          }

          const $position = selection.$from;

          const anchor = tr.selection.anchor;
          const match = findWordBeforeCursor({ chars, anchor })($position);
          const decorationId = (Math.random() + 1).toString(36).substr(2, 5);

          const lastKey = $position.doc
            .textBetween($position.before(), anchor, '\0', '\0')
            .slice(-2);

          if (
            terminalChars.indexOf(lastKey) !== -1 ||
            ($position.parentOffset == 0 && capitalizeOnNewLine && !hardbreak) || ($position.parentOffset == 0 && !capitalizeOnNewLine && hardbreak)
          ) {
            next.capitalize = true;
          } else {
            next.capitalize = false;
          }
          next.lastKey = lastKey;
          if (manuscripts !== null) {
            const result = manuscripts.filter(function (m) {
              return m.abb === match.query.toLowerCase();
            });
            if (result.length > 0 && prompt == false) {
              prompt = true;
              console.log('start new prompt');
              queryManuscript(match.query, result[0].id);
            }
          }
          if (match) {
            next.active = true;
            next.decorationId = prev.decorationId
              ? prev.decorationId
              : decorationId;
            next.range = match.range;
            next.text = match.text;
            next.query = match.query;
          } else {
            next.active = false;
          }
          if (prev.abbreviate == false) {
            next.abbreviate = true;
          }
        } else {
          next.active = false;
        }

        // Make sure to empty the range if suggestion is inactive
        if (!next.active) {
          next.decorationId = null;
          next.range = {};
          next.query = null;
          next.text = null;
        }

        return next;
      },
    },

    props: {
      // Call the keydown hook if suggestion is active.
      handleKeyDown(view, event) {
        if (manuscript.abb != '' /*&& (event.key == " " || event.key == "Enter")*/) {
          if (event.key == "Enter") {
            promptOnNewLine = true;
          }
          const { range } = this.getState(view.state);
          console.log("range:", range)
          view.dispatch(
            view.state.tr.deleteRange(
              range.from,
              range.to + 1
            )
          );

          manuscript.abb = '';
          return true
        }
        const { active, range, query, abbreviate } = this.getState(view.state);
        let { capitalize, lastKey } = this.getState(view.state);
        let { schema } = view.state;

        let send = false;
        let newLine = false;
        if (!active) {
          return false;
        }
        /*        
        if (prompt && manuscript.line < manuscript.content.length) {
                  console.log(
                    'line:',
                    manuscript.line,
                    'of:',
                    manuscript.content.length,
                    '\nindex:',
                    manuscript.index,
                    'of:',
                    manuscript.content[manuscript.line].length
                  );
                }
                */
        if (prompt && !event.ctrlKey && directions.indexOf(event.key) !== -1) {
          let chunks = [];
          event.preventDefault();

          switch (event.key) {
            case 'Escape':
              prompt = false;
              manuscript.line = 0
              manuscript.index = 0
              break;
            case 'ArrowRight':

              chunks = traverseManuscript('promptWord');
              promptManuscriptNodes(view, chunks);
              break;
            case 'ArrowDown':
              chunks = traverseManuscript('promptClause');
              promptManuscriptNodes(view, chunks);
              capitalize = true;
              break;
            case 'ArrowLeft':
              /*  if (event.shiftKey) {
              } else {
              }
              */
              traverseManuscript('removeWord');
              if (manuscript.content[manuscript.line] == undefined) {
                prompt = false
                console.log("stop prompting")
                return
              }

              removeManuscriptNodes(view, false);

              break;
            case 'ArrowUp':
            /*
              traverseManuscript('removeParagraph');
              removeManuscriptNodes(view, true);
              capitalize = true;
              break;
              */
          }


          send = true;
          return onKeyDown({
            view,
            event,
            collapse: true,
            focusLast: true,
            missed: false,
            send,
          });
        }

        if (event.ctrlKey && event.key === ' ') {
          view.dispatch(view.state.tr.insertText('\uFEFF'));
        }

        if (event.key == 'Tab') {
          if (prompt) {
            const tab = schema.text(' ', [createManuscriptMark(schema)]);
            view.dispatch(view.state.tr.insert(range.to, tab));
          } else {
            view.dispatch(view.state.tr.insertText(' ', range.to));
          }
          return;
        }
        let missed = '';
        if (event.key == "Enter") {
          newLine = true
          send = true

          let lastIndex = TextSelection.atEnd(view.state.doc).$anchor.pos
          let cur = view.state.tr.curSelection.$head.pos
          if (lastIndex - cur < 3) {
            return onKeyDown({
              view,
              event,
              range,
              focusLast: true,
              missed,
              send,
              newLine
            });
          }
          if (event.shiftKey) {
            hardbreak = true
            view.dispatch(
              view.state.tr
                .replaceSelectionWith(schema.nodes.paragraph.create())
                .scrollIntoView()
            );
          } else {
            hardbreak = false
          }


        }
        if (chars.indexOf(event.key) !== -1 && abbreviate && query !== '') {
          setTimeout(() => { 
            querying = false }, 125)
          if (!querying) {
            querying = true
            const response = queryAbbreviator(query);
            send = true;

            response.then(res => {
              if (res.data.missed == true) {
                missed = res.data;
                return onKeyDown({
                  view,
                  event,
                  range,
                  focusLast: false,
                  missed,
                });
              }
              const word = res.data.word;
              view.dispatch(view.state.tr.insertText(word, range.from, range.to));
              querying = false
              if (prompt) {
                view.dispatch(
                  view.state.tr.addMark(
                    range.from,
                    range.from + word.length + 1,
                    createManuscriptMark(schema)
                  )
                );
              }
            }).catch(err => {
              console.log("Couldn't abbreviate:", err)
            });
          }
        }

        if (capitalize) {
          if (!event.ctrlKey) {
            if (event.key.length < 2) {
              event.preventDefault();
              view.dispatch(view.state.tr.insertText(event.key.toUpperCase()));
              return;
            }
          }
        }
        return onKeyDown({
          view,
          event,
          range,
          focusLast: false,
          missed,
          send,
          newLine
        });
      },

      // Setup decorator on the currently active suggestion.
      /* decorations(editorState) {
        const { active, range, decorationId } = this.getState(editorState)

        if (!active) return null

        return DecorationSet.create(editorState.doc, [
          Decoration.inline(range.from, range.to, {
            nodeName: 'span',
            class: suggestionClass,
            'data-decoration-id': decorationId,
          }),
        ])
      },*/
    },
  });
}

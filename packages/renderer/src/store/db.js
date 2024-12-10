import Dexie from 'dexie'
const db = new Dexie('sttylus')
db.version(1).stores({
  lists: '++id, name, type, updated',
  abbreviations: '++id, abb, word, updated, remind, uses, lastUse, listId',
})

export default {
  setLists(lists) {
    lists.forEach((list) => {
      const stripped_list = {
        id: list.id,
        name: list.name,
        type: list.type,
        updated: list.updated,
      }
      db.lists.add(stripped_list)
    })
  },
  getLists() {
    return db.table('lists').toArray()
  },
  setAbbs(abbs) {
    console.log('dexie:', abbs)
    abbs.forEach((abb) => {
      const stripped_abb = {
        id: abb.id,
        abb: abb.abb,
        word: abb.word,
        updated: abb.updated,
        remind: abb.remind,
        listId: abb.listId,
      }
      db.abbreviations
        .bulkAdd(stripped_abb)
        .then()
        .catch((err) => console.error(err))
    })
  },
}

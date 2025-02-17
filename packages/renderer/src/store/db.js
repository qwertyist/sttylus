import Dexie from 'dexie'
import api from '../api/api.js'
import EventBus from '../eventbus.js'
import { store } from '../store/index.js'
const db = new Dexie('sttylus')
db.version(1).stores({
  lists: '++id, name, type, updated',
  abbreviations: '++id, abb, word, updated, remind, uses, lastUse, listId',
  syncInfo: 'key, lastSync',
})

export default {
  syncData: async (hard = false) => {
    console.log(hard ? 'Performing full sync' : 'Performing partial sync')
    try {
      let syncStatus = {
        newLists: 0,
        newAbbs: 0,
        errors: [],
      }
      const localListsCount = await db.lists.count()
      const localAbbsCount = await db.abbreviations.count()
      const lastSync = await db.syncInfo.get('lastFullSync')
      const lastSyncTime = lastSync?.lastSync || 0
      if (hard || localListsCount === 0) {
        db.lists
          .clear()
          .then()
          .catch((err) => console.error(err))
        db.abbreviations
          .clear()
          .then()
          .catch((err) => console.error(err))
        db.syncInfo.put({
          key: 'lastFullSync',
          lastSync: 0,
        })
        ////Utför fullständig sync
        await api
          .getUserLists()
          .then(async (resp) => {
            const stripped_lists = resp.data.map((list) => {
              return {
                id: list.id,
                name: list.name,
                type: list.type,
                updated: list.updated,
              }
            })
            db.lists.bulkPut(stripped_lists)
            db.syncInfo.put({
              key: 'lastListSync',
              lastSync: Date.now(),
            })
            await stripped_lists.forEach(async (list, i) => {
              await api
                .getAbbs(list.id)
                .then((resp) => {
                  const stripped_abbs = resp.data.map((abb) => {
                    return {
                      id: abb.id,
                      abb: abb.abb,
                      word: abb.word,
                      updated: abb.updated,
                      remind: abb.remind,
                      listId: list.id,
                      lastUse: 0,
                      uses: 0,
                    }
                  })
                  db.transaction('rw', db.abbreviations, async () => {
                    await db.abbreviations.bulkPut(stripped_abbs)
                  })
                    .then(() => {
                      db.syncInfo.put({
                        key: 'lastFullSync',
                        lastSync: Date.now(),
                      })
                      if (i + 1 == stripped_lists.length) {
                        store.commit('setDbSynced')
                        console.log(
                          'put abbs from last list, updating synced state'
                        )
                        EventBus.$emit('syncedAbbs')
                      }
                    })
                    .catch((err) => console.error(err))
                })
                .catch((err) => {
                  console.error(
                    'db storage could not retrieve user abbs form api:',
                    err,
                    list.name
                  )
                })
            })
          })
          .catch((err) => {
            console.error(
              'db storage could not retrieve user lists from api:',
              err
            )
          })
      } else {
        console.log('check for updates')
        db.table('lists')
          .toArray()
          .then((lists) => {
            lists.forEach((list) => {
              api.getList(list.id).then((resp) => {
                if (resp.data.updated > list.updated) {
                  api.getAbbs(list.id).then((resp) => {
                    const stripped_abbs = resp.data.map((abb) => {
                      return {
                        id: abb.id,
                        abb: abb.abb,
                        word: abb.word,
                        updated: abb.updated,
                        remind: abb.remind,
                        listId: list.id,
                        lastUse: 0,
                        uses: 0,
                      }
                    })
                    db.transaction('rw', db.abbreviations, async () => {
                      await db.abbreviations.bulkPut(stripped_abbs)
                    })
                    db.lists.put(list)
                  })
                }
              })
            })
          })
      }
    } catch (err) {
      console.error('sync failed', err)
    }
  },
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
  addAbb(abb, listId) {
    const stripped_abb = {
      id: abb.id,
      abb: abb.abb,
      word: abb.word,
      updated: abb.updated,
      remind: abb.remind,
      listId: listId,
    }

    db.abbreviations
      .add(stripped_abb)
      .then(() => {
        console.log("created abb:", stripped_abb)
        EventBus.$emit('getAbbCache')
      })
      .catch((err) => console.error(err))
  },
  updateAbb(abb, listId) {

    db.abbreviations
      .where({ abb: abb, listId: listId })
      .delete()
      .then(() => {
        addAbb(abb, listId)
      })
  },
  deleteAbb(abb, listId) {
    db.abbreviations
      .where({ abb: abb, listId: listId })
      .delete()
      .then(() => {
        EventBus.$emit('getAbbCache')
      })
  },
  useAbb(abb) {
    db.abbreviations.where(abb).modify((usedAbb) => {
      ++usedAbb.uses
      usedAbb.lastUse = Date.now()
    })
  },
  setAbbs(abbs) {
    const stripped_abbs = abbs.map((abb) => {
      return {
        id: abb.id,
        abb: abb.abb,
        word: abb.word,
        updated: abb.updated,
        remind: abb.remind,
        listId: abb.listId,
      }
    })
    db.abbreviations
      .bulkAdd(stripped_abbs)
      .then()
      .catch((err) => console.error(err))
  },
  getAbbCache() {
    return db
      .table('abbreviations')
      .toArray()
      .then((abbs) => {
        let cache = new Map()
        const listIds = [store.state.settings.selectedLists.standard].concat(
          store.state.settings.selectedLists.addon
        )
        listIds.forEach((listId) => {
          abbs.forEach((abb) => {
            if (abb.listId == listId) {
              cache.set(abb.abb, abb.word)
            }
          })
        })

        return cache
      })
  },
  lastSyncOk: async () => {
    const lastSync = await db.syncInfo
      .get('lastFullSync')
      .catch((err) => console.error(err))
    const lastSyncTime = lastSync?.lastSync || 0
    return lastSyncTime > 0
  },
}

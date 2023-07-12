import api from './api.js'
var cachedLists = {}
var cachedAbbs = {}
var viewedList = {}
export default {
    initListCache(listIds) {
        cachedLists = new Map()
        cachedAbbs = new Map()
        api.getUserLists()
            .then((resp) => {
                resp.data.forEach((l) => {
                    cachedLists.set(l.id, l)
                })
                cachedLists.forEach((l) => {
                    api.getAbbs(l.id).then((resp) => {
                        cachedAbbs.set(l.id, resp.data)
                    })
                })
            })
            .catch((err) => {})
    },
    cacheViewedAbbs(abbs) {},
    viewList(list) {
        viewedList = list
    },
    getCachedAbbs(listId) {
        return cachedAbbs.get(listId)
    },
    provider(ctx, callback) {
        callback(cachedAbbs.get(viewedList.id))
    },
}

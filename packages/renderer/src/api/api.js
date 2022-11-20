import axios from "axios";
import { store } from "../store"
import { setupCache } from "axios-cache-adapter"

const cache = setupCache({
  maxAge: 28800,
  key: req => {
    return req.identifier
  },
  invalidate: async (cfg, request) => {
    if (!cfg) {
      if (request.createAbb) {
        await cfg.store.removeItem(request.identifier);
      }
      if (request.clearCache) {
        await cfg.store.clear();
      }
    }
  }
})

function createAxiosInstance() {
  if (import.meta.env.VITE_STTYLUS_MODE != "desktop") {
//    console.log("Axios instance points to production server:\n", import.meta.env.VITE_STTYLUS_BACKEND + "/api2")
    let api = "/api2"
    if (import.meta.env.VITE_STTYLUS_DEBUG == "true"){
      api = "/"
    }
    return axios.create({
      baseURL: import.meta.env.VITE_STTYLUS_BACKEND + api,
      headers: {
        'Cache-Control': 'no-cache',
        common: {
          "X-Id-Token": localStorage.getItem("user-token")
        }
      }
    })
  } else {
 //   console.log("Axios instance points to local server:\n", import.meta.env.VITE_STTYLUS_LOCAL_BACKEND + "/api")
    return axios.create({
      baseURL: import.meta.env.VITE_STTYLUS_LOCAL_BACKEND + "/",
      headers: {
        'Cache-Control': 'no-cache',
        common: {
          "X-Id-Token": localStorage.getItem("user-token")
        }
      }
    })
  }
}
const HTTP = createAxiosInstance();

export default {
  checkForToken() {
    return HTTP.defaults.headers.common["X-Id-Token"]
  },
  setIdToken(token) {
    HTTP.defaults.headers.common["X-Id-Token"] = store.state.userData.id
    return
  },
  getLocalIP() {
    return HTTP.get("/ip", { cache: { ignoreCache: true } })
  },
  login(credentials) {
    return HTTP.post("/login", credentials)
  },
  register(credentials) {
    return HTTP.post("/register", credentials)
  },
  isRegistered(email) {
    return HTTP.get("/registered/" + email, { cache: { ignoreCache: true } })
  },
  resetPassword(email, password) {
    return HTTP.put("/password", { email: email, password: password })
  },
  syncUser(credentials) {
    console.log("sync with:", credentials)
    return HTTP.post("/sync", credentials)
  },
  getList(id) {
    return HTTP.get("/abbs/list/" + id, { cache: { ignoreCache: true } })
  },
  getLists(listIDs) {
    return HTTP.post("/abbs/lists", { list_ids: listIDs })
  },
  filterAbbs(ctx) {
    return HTTP.post("/abbs/filter", ctx)
  },
  getUserLists() {
    return HTTP.get("/abbs/lists", { cache: { ignoreCache: true, identifer: "user_lists" } })
  },
  getUserListsByID(id) {
    return HTTP.get("/abbs/lists", {
      headers: {
        common: {
          "X-Id-Token": id
        }
      },
      cache: { ignoreCache: true }
    })
  },
  getAbb(listID, abb) {
    return HTTP.get("/abbs/abbreviation/" + listID + "/" + abb, { cache: { ignoreCache: true } })
  },
  getAbbs(listID) {
    return HTTP.post("/abbs/abbreviations/" + listID)
  },
  createList(list) {
    return HTTP.post("/abbs/list", list)
  },
  updateList(list) {
    return HTTP.put("/abbs/list", list)
  },
  deleteList(id) {
    return HTTP.delete("/abbs/list/" + id)
  },
  createAbb(listID, abb) {
    return HTTP.post("/abbs/abbreviation/" + listID, {
      abb: abb.abb,
      word: abb.word,
      creator: store.state.user.id
    }, { createAbb: true, identifier: abb.abb }
    )
  },
  updateAbb(listID, abb) {
    return HTTP.put("/abbs/abbreviation/" + listID, abb)
  },
  deleteAbb(listID, abb) {
    return HTTP.delete("/abbs/abbreviation/" + listID + "/" + abb.abb)
  },
  cacheAbbs() {
    return HTTP.post("/abbs/cache", { standard: store.state.settings.selectedLists.standard, addon: store.state.settings.selectedLists.addon }, { clearCache: true })
  },
  abbreviate(abb) {
    return HTTP.get("/abbs/abbreviate/" + encodeURIComponent(abb), { identifier: abb })
  },
  initSharedList(baseListId = "") {

    console.log("make post request", baseListId)
    if (baseListId) {
      return HTTP.post("/abbs/shared", { id: baseListId }, { cache: { ignoreCache: true } })
    } else {
      return HTTP.get("/abbs/shared", { cache: { ignoreCache: true } })
    }
  },
  joinSharedList(listid) {
    return HTTP.put("/abbs/shared/" + listid)
  },
  leaveSharedList() {
    return HTTP.put("/abbs/shared/leave")
  },
  createSharedAbb(listid, abb) {
    return HTTP.post("/abbs/shared/" + listid, abb)
  },
  deleteSharedAbb(listid, abb) {
    return HTTP.delete("/abbs/shared/" + listid + "/" + abb)
  },
  getSharedAbbs(listid) {
    return HTTP.get("/abbs/shared/" + listid, { cache: { ignoreCache: true } })
  },
  lookup(phrase) {
    return HTTP.get("/abbs/lookup/" + phrase, { cache: { ignoreCache: true } })
  },
  dontRemind(abb) {
    return HTTP.get("/abbs/learned/" + abb, { cache: { ignoreCache: true } })
  },
  getSuggestions() {
    return HTTP.get("/abbs/suggestions", { identifier: "suggestions", cache: { ignoreCache: true } })
  },
  ignoreSuggestion(word) {
    return HTTP.delete("/abbs/suggestions/" + word)
  },
  ignoreAllSuggestions() {
    return HTTP.delete("/abbs/suggestions")
  },
  importAbbs(abbs) {
    return HTTP.post("/abbs/import", abbs)
  },
  importList(listID, abbs) {
    return HTTP.post("/abbs/import/" + listID, abbs)
  },
  uploadProtype(data) {
    return HTTP.post("/abbs/upload/protype", data)
  },
  uploadTextOnTop(data) {
    return HTTP.post("/abbs/upload/textontop", data)
  },
  uploadIllumiType(data) {
    return HTTP.post("/abbs/upload/illumitype", data)
  },
  uploadTxt(abbs) {
    return HTTP.post("/abbs/upload/txt", abbs)
  },
  importSTTylus(short_id) {
    return HTTP.get("/abbs/public/" + short_id)
  },
  exportSTTylus(list) {
    console.log("api got list:", list)
    return HTTP.post("/abbs/public/" + list.id, list )
  },
  exportLists(target, lists) {
    return HTTP.post("/abbs/export/" + target, { user: store.state.userData.name, standard: lists.standard, addon: lists.addons }, { responseType: "blob", cache: { ignoreCache: true } })
  },
  resolveConflicts(listID, abbs) {
    return HTTP.post("/abbs/conflicts/" + listID, abbs)
  },
  createManuscript(doc) {
    return HTTP.post("/docs", doc)
  },
  updateManuscript(doc) {
    return HTTP.put("/docs", doc)
  },
  deleteManuscript(id) {
    return HTTP.delete("/docs/" + id)
  },
  getManuscript(id) {
    return HTTP.get("/docs/" + id, { cache: { ignoreCache: true } })
  },
  getManuscripts() {
    return HTTP.get("/docs", { cache: { ignoreCache: true } })
  },
  importManuscript(data) {
    return HTTP.post("/docs/import", data)
  },
  saveSettings(settings) {
    if (settings == null) {
      settings = store.state.settings
    }
    return HTTP.post("/settings/", settings)
  },
  getSettings() {
    return HTTP.get("/settings/", { cache: { ignoreCache: true } })
  },
  getUsers() {
    return HTTP.get("/users", { cache: { ignoreCache: true } })
  },
  updateUser(user) {
    return HTTP.put("/user/" + user.id, user)
  },
  createUser(user) {
    return HTTP.post("/user", user)
  },
  createStandard(userID) {
    return HTTP.get("/abbs/standardlist/" + userID, { cache: { ignoreCache: true } })
  },
  deleteUser(user) {
    return HTTP.delete("/user/" + user.id)
  },
  startLocal(options) {

  },
  stopLocal() {

  }


}

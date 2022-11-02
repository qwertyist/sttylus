const version = require("./package.json").version
const axios = require("axios");
const request = require("request");
const compare = require("compare-versions");
const AdmZip = require("adm-zip");
const fs = require("fs");

mainWindowOptions = {
    width: 1200,
    height: 700,
}

updateWindowOptions = {
    width: 800,
    height: 600,
    frame: false,
}

let pid = 0;

chrome.developerPrivate.openDevTools({
    renderViewId: -1,
    renderProcessId: -1,
    extensionId: chrome.runtime.id
})


function checkForUpdates() {
    console.log("Installed version:", version)    
    console.log("Checking for updates...")
    axios.get("https://sttylus.se/latest")
    .then(resp => {
        const latest = resp.data
        console.log("latest version:", latest)
        if(compare(latest.version, version)) {
            nw.Window.open("update.html", updateWindowOptions, function(update_win) {
                global.current = version
                global.latest = latest
                request.get({ url:"https://sttylus.se/get/latest", encoding: null}, (err, res, body) => {
                    var zip = new AdmZip(body);
                    var entries = zip.getEntries();
                    console.log(entries.length)
                    var path = require('path');
                    var nwDir = path.dirname(process.execPath);
                    var dir = nwDir + "/versions/" + latest.version
                    if (!fs.existsSync(dir)) { 
                        fs.mkdirSync(dir) 
                        zip.extractAllTo(dir, true)
                    } 
                    else { 
                        console.log("Version", latest.version, "already downloaded") 
                    }
                })
            })
        }
    })
    .catch(err => {
        console.error("couldnt check latest updates:", err);
    })
    return false
}

function patchApplication() {

}

/*
if (import.meta.env.MODE == "production") {
    path = "./tabula.exe"
} else {
   path = "./backend"
}
*/
var spawn = require("child_process").spawn;
var path = "./backend.exe"

//checkForUpdates()
pid = -1
function loadBackend() {
    const backend = spawn(path)
    backend.stdout.on("data", (data) => {
        console.log(`stdout: ${data}`)
    })
    backend.stderr.on("data", (data) => {
        console.log(`stdout: ${data}`)
    })
    backend.on("close", (code) => {
        console.log(`child process exited with code ${code}`)
        pid = -1
    })
    return backend.pid
}

pid = loadBackend()
console.log("Backend loaded with pid:", pid)

function main() {
    nw.Window.open("index.html", mainWindowOptions, function(win) {
        win.on("close", () => {
            console.log("Close backend first?")
            win.hide()
            nw.Window.getAll(windowList => { windowList[1].close(true) })

            if(pid != -1) {
                process.kill(pid)
            }
            win.close(true)
        })
    })
}

nw.App.registerGlobalHotKey(new nw.Shortcut({
  key: "F11",
  active: function () {
    nw.Window.get().toggleFullscreen();
  },
  failed: function() {
    alert("failed!")
  }
}))
main()

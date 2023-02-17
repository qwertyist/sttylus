/*chrome.developerPrivate.openDevTools({
    renderViewId: -1,
    renderProcessId: -1,
    extensionId: chrome.runtime.id
})
*/
const version = require("./package.json").version
const axios = require("axios");
const request = require("request");
const compare = require("compare-versions");
const AdmZip = require("adm-zip");
const fs = require("fs");



let appData = process.env.APPDATA + "\\STTylus"

console.log(process.platform)
if (process.platform != "win32") {
    appData = process.env.HOME + "/.config/STTylus"
}


console.log("Dina inställningar sparas i:", appData)


if (!fs.existsSync(appData)) {
    fs.mkdirSync(appData)
}




mainWindowOptions = {
    width: 1200,
    height: 700,
}

updateWindowOptions = {
    width: 800,
    height: 600,
    frame: false,
}
locateWindowOptions = {
    width: 400,
    height: 150,
    frame: false,
}
var initial = 0;
let pid = -1;
var spawn = require("child_process").spawn;
var path = "./latest/backend.exe"
if (process.platform != "win32") {
    path = "./backend" 
}
var dbFile = {}
global.dbFile = (f) => {
    dbFile = f
}


function createUser() {
    fs.copyFile("user.db", appData + "\\sttylus.db", fs.constants.COPYFILE_EXCL, (err) => {
        if (err) {
            console.error(err)
        }
    })
}

/*
function terminate() {
    let promises = [execution, timeout]
    let result = promises
    return Promise.race(result).then(v => {
        console.log("race finished", v)
    })
}
*/
if (!fs.existsSync(appData + "\\sttylus.db")) {
    createUser()
    pid = loadBackend() 
    main();
} else {
    pid = loadBackend() 
    main();
}
console.log("done")

function checkForUpdates() {
    console.log("Installed version:", version)
    console.log("Checking for updates...")
    axios.get("https://sttylus.se/latest")
        .then(resp => {
            const latest = resp.data
            console.log("latest version:", latest)
            if (compare(latest.version, version)) {
                nw.Window.open("update.html", updateWindowOptions, function (update_win) {
                    global.current = version
                    global.latest = latest
                    request.get({ url: "https://sttylus.se/get/latest", encoding: null }, (err, res, body) => {
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

window.isFullScreen = false;
console.log("Backend loaded with pid:", pid)

function main() {
    nw.Window.open("index.html", mainWindowOptions, function (win) {
        win.on("close", () => {
            console.log("Close backend first?")
            win.hide()
            nw.Window.getAll(windowList => { windowList[1].close(true) })

            if (pid != -1) {
                process.kill(pid)
            }
            win.close(true)
        })
    })
    nw.App.registerGlobalHotKey(new nw.Shortcut({
        key: "F11",
        active: function () {
            if (window.isFullScreen) {
                nw.Window.get().leaveFullscreen();
                window.isFullScreen = false;
            } else {
                nw.Window.get().enterFullscreen();
                window.isFullScreen = true;
            }
        }
    }))
}

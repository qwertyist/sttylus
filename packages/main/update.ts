import { autoUpdater } from "electron"

export function update(): void {
    autoUpdater.setFeedURL({ url: "http://localhost:3000/"})
    autoUpdater.on("error", (err) => {
        console.log("autoupdater failed:", err)
    })
    autoUpdater.on("checking-for-update", () => {
        console.log("Checking for update:")
    })
    setInterval(() => {
        autoUpdater.checkForUpdates()
    }, 60000)
}
var spawn = require("child_process").spawn;
var path = ""

if (import.meta.env.MODE == "production") {
    path = "./tabula.exe"
} else {
   path = "./backend"
}

export function loadBackend(): number {
    const backend = spawn(path)
    backend.stdout.on("data", (data : String) => {
        console.log(`stdout: ${data}`)
    })
    backend.stderr.on("data", (data : String) => {
        console.log(`stdout: ${data}`)
    })
    backend.on("close", (code : String) => {
        console.log(`child process exited with code ${code}`)
    })
    return backend.pid
}

const { app, BrowserWindow } = require('electron')

require('electron-reload')(__dirname);

function runApp() {
  const backend = require('child_process').execFile('backend/bin/main')
  createWindow(backend);
}

function createWindow (backend) {
  let win = new BrowserWindow({
    width: 800,
    height: 600,
    webPreferences: {
      nodeIntegration: false,
    }
  })

  win.loadFile('public/index.html')
  win.on('closed', function () {
    // Dereference the window object, usually you would store windows
    // in an array if your app supports multi windows, this is the time
    // when you should delete the corresponding element.
    mainWindow = null

    backend.kill('SIGINT')
  })
}

app.on('ready', runApp)
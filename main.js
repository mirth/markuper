const { app, BrowserWindow } = require('electron')

require('electron-reload')(__dirname);

function createWindow () {
  let win = new BrowserWindow({
    width: 800,
    height: 600,
    webPreferences: {
      nodeIntegration: true
    }
  })

  win.loadFile('public/index.html')
}

app.on('ready', createWindow)
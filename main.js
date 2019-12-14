const { app, BrowserWindow } = require('electron');
const { execFile } = require('child_process');

require('electron-reload')(__dirname);

function createWindow(backend) {
  let win = new BrowserWindow({
    width: 800,
    height: 600,
    webPreferences: {
      nodeIntegration: false,
    },
  });

  win.loadFile('public/index.html');
  win.on('closed', () => {
    // Dereference the window object, usually you would store windows
    // in an array if your app supports multi windows, this is the time
    // when you should delete the corresponding element.
    win = null;

    backend.kill('SIGINT');
  });
}

function runApp() {
  const backend = execFile('backend/bin/main');
  backend.stdout.on('data', (chunk) => {
    // eslint-disable-next-line no-console
    console.log(chunk);
  });
  backend.stderr.on('data', (chunk) => {
    // eslint-disable-next-line no-console
    console.log(chunk);
  });

  createWindow(backend);
}


app.on('ready', runApp);

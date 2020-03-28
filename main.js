const { app, BrowserWindow } = require('electron');
const { autoUpdater } = require('electron-updater');
const runBackend = require('./run_backend');


autoUpdater.setFeedURL({
  provider: 'github',
  repo: 'markuper',
  owner: 'mirth',
  private: true,
  token: process.env.GH_TOKEN,
});

if (process.env.ENV === 'dev') {
  // eslint-disable-next-line global-require
  require('electron-reload')(__dirname);
}

function createWindow(backend) {
  let win = new BrowserWindow({
    width: 1024,
    height: 768,
    webPreferences: {
      nodeIntegration: false,
    },
  });

  win.loadFile('public/index.html', { hash: '#' });
  win.on('closed', () => {
    // Dereference the window object, usually you would store windows
    // in an array if your app supports multi windows, this is the time
    // when you should delete the corresponding element.
    win = null;

    backend.kill('SIGINT');
  });
}


function runApp() {
  const backend = runBackend();

  if (process.env.ENV === 'dev') {
  //  autoUpdater.checkForUpdates();
  } else {
   autoUpdater.checkForUpdatesAndNotify();
  }

  createWindow(backend);
}


app.on('ready', runApp);

const { app, BrowserWindow } = require('electron');
const { autoUpdater } = require('electron-updater');
const log = require('electron-log');
const ua = require('universal-analytics');
const { v4: uuid } = require('uuid');
const { JSONStorage } = require('node-localstorage');
const runBackend = require('./run_backend');


function setupAnalytics() {
  const nodeStorage = new JSONStorage(app.getPath('userData'));
  const userId = nodeStorage.getItem('userid') || uuid();
  nodeStorage.setItem('userid', userId);

  const visitor = ua('UA-73579103-13', userId);
  visitor.event('application', 'app_start').send();
}

setupAnalytics();

autoUpdater.setFeedURL({
  provider: 's3',
  bucket: 'markuper-release-dev',
  channel: 'alpha',
});

function sendStatusToWindow(message) {
  // eslint-disable-next-line no-console
  console.log(message);
}

log.transports.file.level = 'debug';
autoUpdater.logger = log;

autoUpdater.on('download-progress', (progressObj) => {
  let logMessage = `Download speed: ${progressObj.bytesPerSecond}`;
  logMessage = `${logMessage} - Downloaded ${parseInt(progressObj.percent, 10)}%`;
  logMessage = `${logMessage} (${progressObj.transferred}/${progressObj.total})`;
  sendStatusToWindow(logMessage);
});

autoUpdater.on('update-downloaded', () => {

});


if (process.env.NODE_ENV === 'dev') {
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
  // win.maximize();
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

  if (process.env.NODE_ENV === 'dev') {
    autoUpdater.checkForUpdates();
  } else {
    autoUpdater.checkForUpdatesAndNotify();
  }

  createWindow(backend);
}


app.on('ready', runApp);

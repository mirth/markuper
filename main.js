const { app, BrowserWindow, autoUpdater, dialog } = require('electron');
const runBackend = require('./run_backend');

const server = 'https://hazel.mirth.now.sh';
const feed = `${server}/update/${process.platform}/${app.getVersion()}`;

autoUpdater.on('update-downloaded', (event, releaseNotes, releaseName) => {
  const dialogOpts = {
    type: 'info',
    buttons: ['Restart', 'Later'],
    title: 'Application Update',
    message: process.platform === 'win32' ? releaseNotes : releaseName,
    detail: 'A new version has been downloaded. Restart the application to apply the updates.'
  }

  dialog.showMessageBox(dialogOpts).then((returnValue) => {
    if (returnValue.response === 0) autoUpdater.quitAndInstall()
  })
})

console.log('feed: ', feed)
autoUpdater.setFeedURL(feed);

autoUpdater.checkForUpdates();

if (process.env.ENV === 'dev') {
  // eslint-disable-next-line global-require
  require('electron-reload')(__dirname);
}

function createWindow(backend) {
  let win = new BrowserWindow({
    width: 800,
    height: 600,
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

  createWindow(backend);
}


app.on('ready', runApp);

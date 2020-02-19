const { app, BrowserWindow, autoUpdater } = require('electron');
const runBackend = require('./run_backend');

const server = 'https://zeit.ink/5F';
const feed = `${server}/update/${process.platform}/${app.getVersion()}`;

autoUpdater.setFeedURL(feed);

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

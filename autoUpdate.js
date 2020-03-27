/* eslint-disable no-console */
const { autoUpdater } = require('electron-updater');
const log = require('electron-log');

function sendStatusToWindow(message) {
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
  const msg = 'After app relauch you will have had a new verion';
  const notify = new Notification('Update downloaded', {
    body: msg,
  });
  notify.onclick = () => {
    console.log('Notification clicked');
  };
  // sendStatusToWindow('');
});

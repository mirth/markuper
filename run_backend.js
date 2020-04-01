const { app } = require('electron');
const { execFile } = require('child_process');
const appRootDir = require('app-root-dir').get();
const path = require('path');

function runBackend() {
  const binaryFilename = 'main.exe';
  const backendPath = (process.env.NODE_ENV === 'dev' || process.env.NODE_ENV === 'test') ? `public/${binaryFilename}`
    : path.join(appRootDir, 'public', binaryFilename);

  const backend = execFile(backendPath, [`-appversion=${app.getVersion()}`], {
    env: {
      NODE_ENV: process.env.NODE_ENV,
    },
  });
  backend.stdout.on('data', (chunk) => {
    // eslint-disable-next-line no-console
    console.log(chunk);
  });
  backend.stderr.on('data', (chunk) => {
    // eslint-disable-next-line no-console
    console.log(chunk);
  });

  return backend;
}

module.exports = runBackend;

const { execFile } = require('child_process');
const appRootDir = require('app-root-dir').get();
const path = require('path');
const os = require('os');

function getBackendBinaryFilename() {
  if (os.platform() === 'win32') {
    return 'main.exe';
  }

  return 'main';
}

function runBackend() {
  const binaryFilename = getBackendBinaryFilename();
  const backendPath = (process.env.ENV === 'dev' || process.env.ENV === 'test') ? `backend/bin/${binaryFilename}`
    : path.join(appRootDir, 'backend', 'bin', binaryFilename);

  const backend = execFile(backendPath, {
    env: {
      ENV: process.env.ENV,
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

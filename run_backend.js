const { execFile } = require('child_process');
const appRootDir = require('app-root-dir').get();
const path = require('path');

function runBackend() {
  const backendPath = (process.env.ENV === 'dev' || process.env.ENV === 'test') ? 'backend/bin/main'
    : path.join(appRootDir, 'backend', 'bin', 'main');

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

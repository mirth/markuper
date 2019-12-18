const { execFile } = require('child_process');

function runBackend() {
  const backend = execFile('backend/bin/main', {
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

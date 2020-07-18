*Markuper*

`yarn start:dev`
`yarn test`
`yarn lint`
`yarn release`

*Publish macOS*
AWS_ACCESS_KEY_ID=XXXXXX AWS_SECRET_ACCESS_KEY=XXXX yarn publish:macos -p always

*Run single test file*
yarn cross-env-shell NODE_ENV=test yarn mocha --require @babel/register --require babel-polyfill ./src/e2e/BoundingBox.e2e.js
# Markuper - dala labeling tool

https://markuper.app/

`yarn start:dev`

`yarn test`

`yarn lint`


__Publish macOS__

`AWS_ACCESS_KEY_ID=XXXXXX AWS_SECRET_ACCESS_KEY=XXXX yarn publish:macos -p always`

__Run single test file__

`yarn cross-env-shell NODE_ENV=test yarn mocha --require @babel/register --require babel-polyfill ./src/e2e/BoundingBox.e2e.js`

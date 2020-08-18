/* eslint-disable no-eval */
/* eslint-disable no-restricted-syntax */
/* eslint-disable prefer-template */
/* eslint-disable consistent-return */
/* eslint-disable func-names */
import { Application } from 'spectron';
import electronPath from 'electron';
import path from 'path';
import {
  itNavigatesToProject, clickButton, expectSampleMarkupToBeEq, assertSampleUri,
} from './test_common';

const appPath = path.join(__dirname, '../..');
const app = new Application({
  path: electronPath,
  args: [appPath],
  webdriverOptions: {
    deprecationWarnings: false,
  },
});


describe('Skip button test', function () {
  this.timeout(20000);
  before(() => app.start());
  after(() => {
    if (app && app.isRunning()) {
      return app.stop();
    }
  });

  const xml = `
  <content>
    <bounding_box group="bbox">
      <radio group="animal" value="cat" visual="Cat" />
      <radio group="animal" value="dog" visual="Dog" />
    </bounding_box>
  </content>
  `;

  itNavigatesToProject(app, appPath, xml);

  it('begins assess', async () => {
    await app.client.waitUntilTextExists('span', 'Begin assess');
    await clickButton(app, 'span', 'Begin assess');
  });

  const imgDir = path.join(appPath, 'src', 'e2e', 'test_data', 'proj0');

  it('shows first sample', async () => {
    await assertSampleUri(app, imgDir, 'kek0.jpg');
  });

  it('skips sample', async () => {
    await app.client.waitForExist('//*[@id="skip_sample_button"]/button');
    await app.client.element('//*[@id="skip_sample_button"]/button').click();
  });

  it('shows second sample', async () => {
    await assertSampleUri(app, imgDir, 'kek1.jpg');
  });

  expectSampleMarkupToBeEq(app, appPath, '[SKIPPED]');
});

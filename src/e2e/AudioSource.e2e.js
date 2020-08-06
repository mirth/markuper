/* eslint-disable no-eval */
/* eslint-disable no-restricted-syntax */
/* eslint-disable prefer-template */
/* eslint-disable consistent-return */
/* eslint-disable func-names */
import { Application } from 'spectron';
import electronPath from 'electron';
import path from 'path';
import { expect } from 'chai';
import {
  itNavigatesToProject, clickButton, createProjectWithAudio, makeSampleUri,
} from './test_common';

const appPath = path.join(__dirname, '../..');
const app = new Application({
  path: electronPath,
  args: [appPath],
  webdriverOptions: {
    deprecationWarnings: false,
  },
});

describe('Audio sample simple test', function () {
  this.timeout(20000);
  before(() => app.start());
  after(() => {
    if (app && app.isRunning()) {
      return app.stop();
    }
  });

  const xml = `
  <content>
    <radio group="animal" value="cat" vizual="Cat" />
    <radio group="animal" value="dog" vizual="Dog" />
  </content>
  `;

  itNavigatesToProject(app, appPath, xml, createProjectWithAudio);

  it('begins assess', async () => {
    await app.client.waitUntilTextExists('span', 'Begin assess');
    await clickButton(app, 'span', 'Begin assess');
  });

  it('contains correct audio source', async () => {
    const imgDir = path.join(appPath, 'src', 'e2e', 'test_data', 'proj0');

    await app.client.waitForExist('//*[@id="audio-source"]/source');
    const src = await app.client.element('//*[@id="audio-source"]/source').getAttribute('src');
    const expectedSampleUri = makeSampleUri(imgDir, 'track0.mp3');
    expect(path.normalize(src)).to.be.eq(path.normalize('file://' + expectedSampleUri));
  });
});

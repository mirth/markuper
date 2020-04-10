/* eslint-disable no-restricted-syntax */
/* eslint-disable prefer-template */
/* eslint-disable consistent-return */
/* eslint-disable func-names */
import { Application } from 'spectron';
import electronPath from 'electron';
import _ from 'lodash';
import path from 'path';
import { expect } from 'chai';
import {
  getPath, getRadio, assertRadioLabels, itNavigatesToProject, getSamplePath, getSampleClass, sleep,
  clickButton, getRadioState, getChecked, clickLink,
} from './test_common';

const appPath = path.join(__dirname, '../..');
const app = new Application({
  path: electronPath,
  args: [appPath],
});


async function selectRect(upperLeft, downRight) {
  await app.client.moveToObject('#image-container', upperLeft[0], upperLeft[1]);
  await app.client.buttonDown(0);
  await app.client.moveToObject('#image-container', downRight[0], downRight[1]);
}

describe('Kek', function () {
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
      <radio group="animal" value="cat" vizual="Cat" />
      <radio group="animal" value="dog" vizual="Dog" />

      <checkbox group="color" value="black" vizual="Black" />
      <checkbox group="color" value="white" vizual="White" />
    </bounding_box>
  </content>
  `;

  itNavigatesToProject(app, appPath, xml);

  it('begins assess', async () => {
    await app.client.waitUntilTextExists('span', 'Begin assess');
    await clickButton(app, 'span', 'Begin assess');
  });

  it('lulz', async () => {
    // await app.client.leftClick('#image-container', 0, 0);

    selectRect([10, 10], [100, 100]);
    // await app.client.buttonUp(0);
    await sleep(10000)
  })

})

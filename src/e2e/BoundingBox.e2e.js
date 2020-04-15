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
  itNavigatesToProject, sleep, clickButton, getRadioState, getRadio, expectSampleMarkupToBeEq,
  expectBoxEqual,
} from './test_common';

const appPath = path.join(__dirname, '../..');
const app = new Application({
  path: electronPath,
  args: [appPath],
  webdriverOptions: {
    deprecationWarnings: false,
  },
});


async function selectRect(upperLeft, downRight) {
  await app.client.moveToObject('#image-container', upperLeft[0], upperLeft[1]);
  await app.client.buttonDown(0);
  await app.client.moveToObject('#image-container', downRight[0], downRight[1]);
  await app.client.buttonUp(0);
}

async function getBoxesMarkup() {
  const boxes = '//*[@id="boxes"]/table/tbody/tr/td/div/div/div[1]/div/div/span/small';
  await app.client.waitForExist(boxes);
  const elements = await app.client.elements(boxes);
  const markup = await Promise.all(elements.value.map(async (el) => {
    const txt = await app.client.elementIdText(el.ELEMENT);
    return txt.value;
  }));

  return markup;
}

// <checkbox group="color" value="black" vizual="Black" />
// <checkbox group="color" value="white" vizual="White" />

describe('Simple bounding box test', function () {
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
    </bounding_box>
  </content>
  `;

  itNavigatesToProject(app, appPath, xml);

  it('begins assess', async () => {
    await app.client.waitUntilTextExists('span', 'Begin assess');
    await clickButton(app, 'span', 'Begin assess');
  });

  it('lulz', async () => {
    selectRect([10, 10], [100, 100]);

    await getRadio(app, 'bbox/0', 2).click();
    await sleep(500);

    const selected = await getRadioState(app, 'bbox/0');
    expect(selected).to.be.deep.eq([false, true]);
    await app.client.keys('Enter');

    const mark = await getBoxesMarkup();
    expectBoxEqual(eval('({' + mark + '})'), {
      left: 0, top: 0, width: 4, height: 4,
    });
    await sleep(1500);
    await app.client.keys('Enter');
    await sleep(1500);
    await app.client.keys('Enter');
  });

  it('goes on project page', async () => {
    await clickButton(app, 'span', 'testproj0');
    await app.client.waitForText('span', 'Begin assess');
  });

  expectSampleMarkupToBeEq(app, appPath, {
    bbox: [{
      box: {
        x: 0, y: 0, width: 4, height: 4,
      },
      animal: 'dog',
    }],
  });
});

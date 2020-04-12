/* eslint-disable no-restricted-syntax */
/* eslint-disable prefer-template */
/* eslint-disable consistent-return */
/* eslint-disable func-names */
import { Application } from 'spectron';
import electronPath from 'electron';
import path from 'path';
import { expect, assert } from 'chai';
import {
  itNavigatesToProject, sleep, clickButton, getRadioState,
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

function almostEqual(rawActual, etalon) {
  // eslint-disable-next-line no-eval
  const actual = eval('({' + rawActual + '})');
  const thresh = 2;
  assert.approximately(actual.left, etalon.left, thresh, 'numbers are close');
  assert.approximately(actual.top, etalon.top, thresh, 'numbers are close');
  assert.approximately(actual.width, etalon.width, thresh, 'numbers are close');
  assert.approximately(actual.height, etalon.height, thresh, 'numbers are close');
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

    // fixme this not working
    // await getRadio(app, 'bbox/0', 2).click();
    await sleep(500);
    await app.client.keys('2');

    const selected = await getRadioState(app, 'bbox/0');
    expect(selected).to.be.deep.eq([false, true]);
    await app.client.keys('Enter');

    const mark = await getBoxesMarkup();

    almostEqual(mark, {
      left: 10, top: 10, width: 90, height: 90,
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

  // fixme
  // -{"bbox":[{"animal":"dog","box":{"height":90,"width":90,"x":8.34375,"y":9}}]}
  // +{"bbox":[{"animal":"dog","box":{"height":90,"width":90,"x":10,"y":10}}]}
  // expectSampleMarkupToBeEq(app, appPath, {
  //   bbox: [{
  //     box: {x: 10, y: 10, width: 90, height: 90},
  //     animal: 'dog',
  //   }],
  // });
});

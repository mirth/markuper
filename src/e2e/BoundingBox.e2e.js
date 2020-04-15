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
  let markup = await Promise.all(elements.value.map(async (el) => {
    const txt = await app.client.elementIdText(el.ELEMENT);
    return txt.value;
  }));
  markup = markup.map((mark) => (eval('({' + mark + '})')));

  return markup;
}

const assertBoxMarkup = (expected, actual) => {
  expect(expected.length).to.be.eq(actual.length);
  actual.forEach((act, i) => {
    expect(act).to.be.deep.eq(expected[i]);
  });
};

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

  it('makes correct bounding box', async () => {
    selectRect([10, 10], [100, 100]);

    await getRadio(app, 'bbox/0', 2).click();
    await sleep(500);

    const selected = await getRadioState(app, 'bbox/0');
    expect(selected).to.be.deep.eq([false, true]);
    await app.client.keys('Enter');

    const actual = await getBoxesMarkup();
    assertBoxMarkup([
      {
        left: 0, top: 0, width: 4, height: 4,
      },
    ], actual);

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


describe('Bounding boxes manipulation', function () {
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

  const drawAndCheckBBox = async (upperLeft, downRight) => {
    selectRect(upperLeft, downRight);

    await getRadio(app, 'bbox/0', 2).click();
    await sleep(500);

    const selected = await getRadioState(app, 'bbox/0');
    expect(selected).to.be.deep.eq([false, true]);
    await app.client.keys('Enter');
    await sleep(500);
  };

  it('makes 1st bounding box', async () => {
    await drawAndCheckBBox([10, 10], [100, 100], [{
      left: 0, top: 0, width: 4, height: 4,
    }]);
  });

  it('makes 2st bounding box', async () => {
    await drawAndCheckBBox([100, 100], [375, 375]);
    const actual = await getBoxesMarkup();
    assertBoxMarkup([
      {
        left: 0, top: 0, width: 4, height: 4,
      },
      {
        left: 4, top: 4, width: 12, height: 12,
      },
    ], actual);
  });

  it('makes 3rd bounding box', async () => {
    await drawAndCheckBBox([200, 200], [400, 400]);
    const actual = await getBoxesMarkup();
    assertBoxMarkup([
      {
        left: 0, top: 0, width: 4, height: 4,
      },
      {
        left: 4, top: 4, width: 12, height: 12,
      },
      {
        left: 8, top: 8, width: 8, height: 8,
      },
    ], actual);
  });

  it('removes 2nd box', async () => {
    await app.client.element('//*[@id="boxes"]/table/tbody/tr[2]/td/div/div/div[2]/button').click();

    const actual = await getBoxesMarkup();
    assertBoxMarkup([
      {
        left: 0, top: 0, width: 4, height: 4,
      },
      {
        left: 8, top: 8, width: 8, height: 8,
      },
    ], actual);
  });

  it('removes 1st box', async () => {
    await app.client.element('//*[@id="boxes"]/table/tbody/tr[1]/td/div/div/div[2]/button').click();

    const actual = await getBoxesMarkup();
    assertBoxMarkup([
      {
        left: 8, top: 8, width: 8, height: 8,
      },
    ], actual);
  });

  it('makes 2d bounding box', async () => {
    await drawAndCheckBBox([10, 10], [100, 100]);
    await sleep(1000);
    const actual = await getBoxesMarkup();
    assertBoxMarkup([
      {
        left: 8, top: 8, width: 8, height: 8,
      },
      {
        left: 0, top: 0, width: 4, height: 4,
      },
    ], actual);

    await app.client.keys('Enter');
    await sleep(1500);
    await app.client.keys('Enter');
    await sleep(1500);
  });

  it('goes on project page', async () => {
    await clickButton(app, 'span', 'testproj0');
    await app.client.waitForText('span', 'Begin assess');
  });

  expectSampleMarkupToBeEq(app, appPath, {
    bbox: [
      {
        box: {
          x: 8, y: 8, width: 8, height: 8,
        },
        animal: 'dog',
      },
      {
        box: {
          x: 0, y: 0, width: 4, height: 4,
        },
        animal: 'dog',
      },
    ],
  });
});

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
  getElement, itNavigatesToProject, sleep, clickButton, getRadioState, getRadio,
  expectSampleMarkupToBeEq,
} from './test_common';
import { TestCheckboxRadioRadio, TestRadioCheckbox } from './classification_common';

const appPath = path.join(__dirname, '../..');
const app = new Application({
  path: electronPath,
  args: [appPath],
  webdriverOptions: {
    deprecationWarnings: false,
  },
});

async function makeToImageCoords() {
  const TEST_IMAGE_SIZE = 32;
  const img = await getElement(app, '//*[@id="image-container"]/img');
  let width = await app.client.elementIdCssProperty(img, 'width');
  width = parseInt(width.value.slice(0, -2), 10);

  return (point) => {
    const [x, y] = point;
    return [Math.round((width / TEST_IMAGE_SIZE) * x), Math.round((width / TEST_IMAGE_SIZE) * y)];
  };
}

async function selectRect(upperLeft, downRight) {
  const convPoint = await makeToImageCoords();
  const img = '#image-container';
  const upperLeftConv = convPoint(upperLeft);
  const downRightConv = convPoint(downRight);
  await app.client.moveToObject(img, upperLeftConv[0], upperLeftConv[1]);
  await app.client.buttonDown(0);
  await app.client.moveToObject(img, downRightConv[0], downRightConv[1]);
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

const clickSubmit = async () => app.client.element('//*[@id="device_submit"]/div/button').click();

// fixme test for checking markup for sample stays

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
    await selectRect([1, 1], [4, 4]);

    await getRadio(app, 'bbox/0', 2).click();
    await sleep(500);

    const selected = await getRadioState(app, 'bbox/0');
    expect(selected).to.be.deep.eq([false, true]);
    await app.client.keys('Enter');

    const actual = await getBoxesMarkup();
    assertBoxMarkup([
      {
        left: 1, top: 1, width: 3, height: 3,
      },
    ], actual);

    await sleep(1500);
    await app.client.keys('Enter');
    await clickSubmit();
  });

  it('goes on project page', async () => {
    await clickButton(app, 'span', 'testproj0');
    await app.client.waitForText('span', 'Begin assess');
  });

  expectSampleMarkupToBeEq(app, appPath, {
    bbox: [{
      box: {
        x: 1, y: 1, width: 3, height: 3,
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
    await selectRect(upperLeft, downRight);

    await getRadio(app, 'bbox/0', 2).click();
    await sleep(500);

    const selected = await getRadioState(app, 'bbox/0');
    expect(selected).to.be.deep.eq([false, true]);
    await app.client.keys('Enter');
    await sleep(500);
  };

  it('makes 1st bounding box', async () => {
    await drawAndCheckBBox([1, 1], [4, 4], [{
      left: 1, top: 1, width: 3, height: 3,
    }]);
  });

  it('makes 2st bounding box', async () => {
    await drawAndCheckBBox([4, 4], [15, 15]);
    const actual = await getBoxesMarkup();
    assertBoxMarkup([
      {
        left: 1, top: 1, width: 3, height: 3,
      },
      {
        left: 4, top: 4, width: 11, height: 11,
      },
    ], actual);
  });

  it('makes 3rd bounding box', async () => {
    await drawAndCheckBBox([8, 8], [16, 16]);
    const actual = await getBoxesMarkup();
    assertBoxMarkup([
      {
        left: 1, top: 1, width: 3, height: 3,
      },
      {
        left: 4, top: 4, width: 11, height: 11,
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
        left: 1, top: 1, width: 3, height: 3,
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
    await drawAndCheckBBox([1, 1], [4, 4]);
    await sleep(1000);
    const actual = await getBoxesMarkup();
    assertBoxMarkup([
      {
        left: 8, top: 8, width: 8, height: 8,
      },
      {
        left: 1, top: 1, width: 3, height: 3,
      },
    ], actual);

    await sleep(1500);
    await app.client.keys('Enter');
    await clickSubmit();
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
          x: 1, y: 1, width: 3, height: 3,
        },
        animal: 'dog',
      },
    ],
  });
});


describe('Bounding box with Focus and state [Checkbox, Radio, Radio]', function () {
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
      <checkbox group="color" value="black" vizual="Black" />
      <checkbox group="color" value="white" vizual="White" />
      <checkbox group="color" value="pink" vizual="Pink" />

      <radio group="animal" value="cat" vizual="Cat" />
      <radio group="animal" value="dog" vizual="Dog" />
      <radio group="animal" value="chuk" vizual="Chuk" />
      <radio group="animal" value="gek" vizual="Gek" />

      <radio group="size" value="smoll" vizual="Smoll" />
      <radio group="size" value="big" vizual="Big" />
    </bounding_box>
  </content>
  `;

  itNavigatesToProject(app, appPath, xml);

  it('begins assess', async () => {
    await app.client.waitUntilTextExists('span', 'Begin assess');
    await clickButton(app, 'span', 'Begin assess');
  });

  it('draw bounding box', async () => {
    await selectRect([4, 4], [15, 15]);
  });

  TestCheckboxRadioRadio(app, 'bbox');

  it('submits the sample', async () => {
    await clickSubmit();
    await sleep(1500);
  });

  it('goes on project page', async () => {
    await clickButton(app, 'span', 'testproj0');
    await app.client.waitForText('span', 'Begin assess');
  });

  expectSampleMarkupToBeEq(app, appPath, {
    bbox: [{
      box: {
        x: 4, y: 4, width: 11, height: 11,
      },
      animal: 'gek',
      color: ['black'],
      size: 'big',
    }],
  });
});

describe('Bounding box with Focus and state [Checkbox, Radio]', function () {
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
      <radio group="animal" value="chuk" vizual="Chuk" />
      <radio group="animal" value="gek" vizual="Gek" />

      <checkbox group="color" value="black" vizual="Black" />
      <checkbox group="color" value="white" vizual="White" />
      <checkbox group="color" value="pink" vizual="Pink" />
    </bounding_box>
  </content>
  `;

  itNavigatesToProject(app, appPath, xml);

  it('begins assess', async () => {
    await app.client.waitUntilTextExists('span', 'Begin assess');
    await clickButton(app, 'span', 'Begin assess');
  });

  it('draw bounding box', async () => {
    await selectRect([4, 4], [15, 15]);
  });

  TestRadioCheckbox(app, 'bbox');

  it('submits the sample', async () => {
    await clickSubmit();
    await sleep(1500);
  });

  it('goes on project page', async () => {
    await clickButton(app, 'span', 'testproj0');
    await app.client.waitForText('span', 'Begin assess');
  });

  expectSampleMarkupToBeEq(app, appPath, {
    bbox: [{
      box: {
        x: 4, y: 4, width: 11, height: 11,
      },
      animal: 'cat',
      color: ['black', 'pink'],
    }],
  });
});

describe('Display correct box labels', function () {
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
      <radio group="animal" value="cat" vizual="Cat" color="#100000" />
      <radio group="animal" value="dog" vizual="Dog" color="#200000" />
      <radio group="animal" value="chuk" vizual="Chuk" color="#300000" />
      <radio group="animal" value="gek" vizual="Gek" color="#400000" />

      <checkbox group="color" value="black" vizual="Black" color="#500000" />
      <checkbox group="color" value="white" vizual="White" color="#600000" />
      <checkbox group="color" value="pink" vizual="Pink" color="#700000" />
    </bounding_box>
  </content>
  `;

  itNavigatesToProject(app, appPath, xml);

  it('begins assess', async () => {
    await app.client.waitUntilTextExists('span', 'Begin assess');
    await clickButton(app, 'span', 'Begin assess');
  });

  it('draw bounding box', async () => {
    await selectRect([4, 4], [15, 15]);
  });

  TestRadioCheckbox(app, 'bbox');

  it('displays box labels correctly', async () => {
    const elements = await app.client.elements('//*[@id="image-container"]/div/ul/li/span');
    const labels = await Promise.all(elements.value.map(async (el) => {
      const txt = await app.client.elementIdText(el.ELEMENT);
      const color = await app.client.elementIdCssProperty(el.ELEMENT, 'background-color');
      return [txt.value, color.value];
    }));
    expect(labels).to.be.deep.eq([
      ['Cat', 'rgba(16, 0, 0, 1)'],
      ['Black', 'rgba(80, 0, 0, 1)'],
      ['Pink', 'rgba(112, 0, 0, 1)'],
    ]);
  });
});

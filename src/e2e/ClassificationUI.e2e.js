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


function expectFocusIsOn(devices, device) {
  it(`focused on device ${device}`, async () => {
    const cls = await Promise.all(devices.map(async (dev) => {
      const cl = await app.client.element(`//*[@id="${dev}"]/div`).getAttribute('class');
      return cl;
    }));

    const pairs = _.zip(devices, cls);

    for (const [iterName, iterCl] of pairs) {
      if (iterName === device) {
        expect(iterCl).to.have.string('selected');
      } else {
        expect(iterCl).to.not.have.string('selected');
      }
    }

    const devSubmitCl = await app.client.element('//*[@id="device_submit"]/div/button').getAttribute('class');

    if (device === 'devSubmit') {
      expect(devSubmitCl).to.have.string('filled');
    }
  });
}

async function assertCheckboxLabels(device) {
  const inputs = `//*[@id="${device}"]/div/label/ul/li/label/input`;
  await app.client.waitForExist(inputs);
  const elements = await app.client.elements(inputs);
  const labels = await Promise.all(elements.value.map(async (el) => {
    const pth = await getPath(app, el, '..');
    const txt = await app.client.elementIdText(pth.value.ELEMENT);
    return txt.value;
  }));
  expect(labels).to.be.deep.eq(['Black', 'White', 'Pink']);
}

function expectSampleMarkupToBeEq(markup) {
  it('displays sample markup on project page', async () => {
    await clickLink(app, 'span', 'testproj0');
    await app.client.waitForExist('ul');

    const imgDir = path.join(appPath, 'src', 'e2e', 'test_data', 'proj0');

    const pathText = await getSamplePath(app, 'kek0.jpg').getText();
    const cl = await getSampleClass(app, 'kek0.jpg').getText();
    expect(pathText).to.be.eq(path.join(imgDir, 'kek0.jpg') + ':');
    expect(cl).to.be.eq(markup);
  });
}

function itSubmitsSample() {
  it('submits the sample', async () => {
    await app.client.keys('Enter');
    await sleep(1500);
    await clickButton(app, 'span', 'testproj0');
    await sleep(1500);
  });
}

describe('Device state keep for assessed samples', function () {
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

    <checkbox group="color" value="black" vizual="Black" />
    <checkbox group="color" value="white" vizual="White" />
  </content>
  `;

  itNavigatesToProject(app, appPath, xml);

  it('begins assess', async () => {
    await app.client.waitUntilTextExists('span', 'Begin assess');
    await clickButton(app, 'span', 'Begin assess');
  });

  it('displays radio 2 selected', async () => {
    await getRadio(app, 'root/0', 2).click();

    const selected = await getRadioState(app, 'root/0');
    expect(selected).to.be.deep.eq([false, true]);
    await app.client.keys('Enter');
  });

  it('displays checkbox 1 as checked', async () => {
    await app.client.keys('1');

    const checked = await getChecked(app, 'root/1');
    expect(checked).to.be.deep.eq([true, false]);
    await app.client.keys('Enter');
  });

  itSubmitsSample();

  it('goes to next new sample', async () => {
    await app.client.waitUntilTextExists('span', 'Begin assess');
    await clickButton(app, 'span', 'Begin assess');

    const selected = await getRadioState(app, 'root/0');
    expect(selected).to.be.deep.eq([false, false]);
    const checked = await getChecked(app, 'root/1');
    expect(checked).to.be.deep.eq([false, false]);
  });

  it('goes back to first assessed sample', async () => {
    await clickButton(app, 'span', 'testproj0');
    await app.client.waitForText('span', 'Begin assess');
    await getSamplePath(app, 'kek0.jpg').element('..').click();
    await app.client.waitForVisible("button/*[@innertext='Cat']");

    const selected = await getRadioState(app, 'root/0');
    expect(selected).to.be.deep.eq([false, true]);
    const checked = await getChecked(app, 'root/1');
    expect(checked).to.be.deep.eq([true, false]);
  });
});

describe('Focus and state [Checkbox, Radio, Radio]', function () {
  this.timeout(20000);
  before(() => app.start());
  after(() => {
    if (app && app.isRunning()) {
      return app.stop();
    }
  });

  const xml = `
  <content>
    <checkbox group="color" value="black" vizual="Black" />
    <checkbox group="color" value="white" vizual="White" />
    <checkbox group="color" value="pink" vizual="Pink" />

    <radio group="animal" value="cat" vizual="Cat" />
    <radio group="animal" value="dog" vizual="Dog" />
    <radio group="animal" value="chuk" vizual="Chuk" />
    <radio group="animal" value="gek" vizual="Gek" />

    <radio group="size" value="smoll" vizual="Smoll" />
    <radio group="size" value="big" vizual="Big" />
  </content>
  `;

  itNavigatesToProject(app, appPath, xml);

  it('begins assess', async () => {
    await app.client.waitUntilTextExists('span', 'Begin assess');
    await clickButton(app, 'span', 'Begin assess');
  });

  const focusIsOn = (device) => {
    expectFocusIsOn(['root/0', 'root/1', 'root/2'], device);
  };

  it('displays correct devices labels', async () => {
    await assertCheckboxLabels('root/0');
    await assertRadioLabels(app, 'root/1', ['Cat', 'Dog', 'Chuk', 'Gek']);
    await assertRadioLabels(app, 'root/2', ['Smoll', 'Big']);
  });

  focusIsOn('root/0');

  it('displays checkbox 1 as checked', async () => {
    await app.client.keys('1');
    await app.client.keys('2');
    await app.client.keys('1');
    await app.client.keys('1');
    await app.client.keys('2');

    const checked = await getChecked(app, 'root/0');
    expect(checked).to.be.deep.eq([true, false, false]);
    await app.client.keys('Enter');
  });

  focusIsOn('root/1');

  it('displays radio 4 selected', async () => {
    await getRadio(app, 'root/1', 1).click();
    await getRadio(app, 'root/1', 1).click();
    await getRadio(app, 'root/1', 4).click();

    const selected = await getRadioState(app, 'root/1');
    expect(selected).to.be.deep.eq([false, false, false, true]);
    await app.client.keys('Enter');
  });

  focusIsOn('root/2');

  it('displays radio 2 selected', async () => {
    await getRadio(app, 'root/2', 2).click();
    const selected = await getRadioState(app, 'root/2');
    expect(selected).to.be.deep.eq([false, true]);
    await app.client.keys('Enter');
  });

  focusIsOn('device_submit');

  itSubmitsSample();
  expectSampleMarkupToBeEq('animal: gek, color: black, size: big');
});

describe('Focus and state [Radio, Checkbox]', function () {
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
    <radio group="animal" value="chuk" vizual="Chuk" />
    <radio group="animal" value="gek" vizual="Gek" />

    <checkbox group="color" value="black" vizual="Black" />
    <checkbox group="color" value="white" vizual="White" />
    <checkbox group="color" value="pink" vizual="Pink" />
  </content>
  `;

  const focusIsOn = (device) => {
    expectFocusIsOn(['root/0', 'root/1'], device);
  };
  itNavigatesToProject(app, appPath, xml);

  it('begins assess', async () => {
    await app.client.waitUntilTextExists('span', 'Begin assess');
    await clickButton(app, 'span', 'Begin assess');
  });


  it('displays correct devices labels', async () => {
    await sleep(2000);
    await assertRadioLabels(app, 'root/0', ['Cat', 'Dog', 'Chuk', 'Gek']);
    await assertCheckboxLabels('root/1');
  });

  focusIsOn('root/0');

  it('tries to submit empty radio field', async () => {
    await app.client.keys('Enter');
    await sleep(1500);
  });

  focusIsOn('root/0');

  it('displays radio 3 selected', async () => {
    await app.client.keys('3');
    const selected = await getRadioState(app, 'root/0');
    expect(selected).to.be.deep.eq([false, false, true, false]);
  });

  focusIsOn('root/0');

  it('submits radio field', async () => {
    await app.client.keys('Enter');
  });

  focusIsOn('root/1');

  it('displays radio 1 selected', async () => {
    await getRadio(app, 'root/0', 1).click();
    const selected = await getRadioState(app, 'root/0');
    expect(selected).to.be.deep.eq([true, false, false, false]);
  });

  focusIsOn('root/1');

  it('displays checkbox 2 as checked', async () => {
    await app.client.keys('2');

    const checked = await getChecked(app, 'root/1');
    expect(checked).to.be.deep.eq([false, true, false]);
  });

  focusIsOn('root/1');

  it('displays checkboxes 1,2 as checked', async () => {
    await app.client.keys('1');

    const checked = await getChecked(app, 'root/1');
    expect(checked).to.be.deep.eq([true, true, false]);
  });

  focusIsOn('root/1');

  it('displays all checkboxes as checked', async () => {
    await app.client.keys('3');

    const checked = await getChecked(app, 'root/1');
    expect(checked).to.be.deep.eq([true, true, true]);
  });

  focusIsOn('root/1');

  it('displays checkbox 2 as unchecked', async () => {
    await app.client.keys('2');

    const checked = await getChecked(app, 'root/1');
    expect(checked).to.be.deep.eq([true, false, true]);
    await app.client.keys('Enter');
  });

  focusIsOn('device_submit');

  itSubmitsSample();
  expectSampleMarkupToBeEq('animal: cat, color: black,pink');
});

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
  getPath, getBtn, assertRadioLabels, itNavigatesToProject, getSamplePath, getSampleClass, sleep,
} from './test_common';

const appPath = path.join(__dirname, '../..');
const app = new Application({
  path: electronPath,
  args: [appPath],
});


function expectFocusIsOn(devices, device) {
  it(`focused on device ${device}`, async () => {
    const cls = await Promise.all(devices.map(async (dev) => {
      const cl = await app.client.element(`//*[@id="${dev}"]`).getAttribute('class');
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
  });
}

const getDisabled = async (device) => {
  const elements = await app.client.elements(`//*[@id="${device}"]/ul/li/div/button`);
  const disabled = await Promise.all(elements.value.map(async (el) => {
    const cl = await app.client.elementIdAttribute(el.ELEMENT, 'class');
    return cl.value.includes('disabled');
  }));

  return disabled;
};
const getChecked = async (device) => {
  const elements = await app.client.elements(`//*[@id="${device}"]/div/ul/li/label/input`);
  const checked = await Promise.all(elements.value.map(async (el) => {
    const ch = await app.client.elementIdSelected(el.ELEMENT, 'checked');
    return ch.value;
  }));

  return checked;
};


async function assertCheckBoxLabels(device) {
  const elements = await app.client.elements(`//*[@id="${device}"]/div/ul/li/label/input`);
  const labels = await Promise.all(elements.value.map(async (el) => {
    const kek = await getPath(app, el, '..');
    const txt = await app.client.elementIdText(kek.value.ELEMENT);
    return txt.value;
  }));
  expect(labels).to.be.deep.eq(['Black', 'White', 'Pink']);
}

function expectSampleMarkupToBeEq(markup) {
  it('displays sample markup on project page', async () => {
    await app.client.element("button/*[@innertext='testproj0']").click();
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
    await sleep(2500);
    await app.client.element("button/*[@innertext='testproj0']").click();
    await sleep(2500);
  });
}

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

    // fixme
    const beginAssess = app.client.element('//*[@id="grid"]/div/div[2]/div/div/div/div/div[1]/div/div[3]/div/div/button[1]');
    await beginAssess.click();
  });

  const focusIsOn = (device) => {
    expectFocusIsOn(['device0', 'device1', 'device2', 'device_submit'], device);
  };

  it('displays correct devices labels', async () => {
    await sleep(2000);
    await assertCheckBoxLabels('device0');
    await assertRadioLabels(app, 'device1', ['Cat', 'Dog', 'Chuk', 'Gek']);
    await assertRadioLabels(app, 'device2', ['Smoll', 'Big']);
  });

  focusIsOn('device0');

  it('displays checkbox 1 as checked', async () => {
    await app.client.keys('1');
    await app.client.keys('2');
    await app.client.keys('1');
    await app.client.keys('1');
    await app.client.keys('2');

    const checked = await getChecked('device0');
    expect(checked).to.be.deep.eq([true, false, false]);
    await app.client.keys('Enter');
  });

  focusIsOn('device1');

  it('displays button 4 pressed', async () => {
    await getBtn(app, 'device1', 1).click();
    await getBtn(app, 'device1', 1).click();
    await getBtn(app, 'device1', 4).click();
    const disabled = await getDisabled('device1');
    expect(disabled).to.be.deep.eq([false, false, false, true]);
    await app.client.keys('Enter');
  });

  focusIsOn('device2');

  it('displays button 2 pressed', async () => {
    await getBtn(app, 'device2', 2).click();
    const disabled = await getDisabled('device2');
    expect(disabled).to.be.deep.eq([false, true]);
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
    expectFocusIsOn(['device0', 'device1', 'device_submit'], device);
  };
  itNavigatesToProject(app, appPath, xml);

  it('begins assess', async () => {
    await app.client.waitUntilTextExists('span', 'Begin assess');

    // fixme
    const beginAssess = app.client.element('//*[@id="grid"]/div/div[2]/div/div/div/div/div[1]/div/div[3]/div/div/button[1]');
    await beginAssess.click();
  });


  it('displays correct devices labels', async () => {
    await assertRadioLabels(app, 'device0', ['Cat', 'Dog', 'Chuk', 'Gek']);
    await assertCheckBoxLabels('device1');
  });

  it('selects first device', async () => {
    const device1 = await app.client.element('//*[@id="device0"]').getAttribute('class');
    expect(device1).to.have.string('selected');
    const device2 = await app.client.element('//*[@id="device1"]').getAttribute('class');
    expect(device2).to.not.have.string('selected');
  });

  focusIsOn('device0');

  it('displays button 3 pressed', async () => {
    await app.client.keys('3');
    const disabled = await getDisabled('device0');
    expect(disabled).to.be.deep.eq([false, false, true, false]);
  });

  focusIsOn('device0');

  it('submits radio field', async () => {
    await app.client.keys('Enter');
  });

  focusIsOn('device1');

  it('displays button 1 pressed', async () => {
    await getBtn(app, 'device0', 1).click();
    const disabled = await getDisabled('device0');
    expect(disabled).to.be.deep.eq([true, false, false, false]);
  });

  focusIsOn('device1');

  it('displays checkbox 2 as checked', async () => {
    await app.client.keys('2');

    const checked = await getChecked('device1');
    expect(checked).to.be.deep.eq([false, true, false]);
  });

  focusIsOn('device1');

  it('displays checkboxes 1,2 as checked', async () => {
    await app.client.keys('1');

    const checked = await getChecked('device1');
    expect(checked).to.be.deep.eq([true, true, false]);
  });

  focusIsOn('device1');

  it('displays all checkboxes as checked', async () => {
    await app.client.keys('3');

    const checked = await getChecked('device1');
    expect(checked).to.be.deep.eq([true, true, true]);
  });

  focusIsOn('device1');

  it('displays checkbox 2 as unchecked', async () => {
    await app.client.keys('2');

    const checked = await getChecked('device1');
    expect(checked).to.be.deep.eq([true, false, true]);
    await app.client.keys('Enter');
  });

  focusIsOn('device_submit');

  itSubmitsSample();
  expectSampleMarkupToBeEq('animal: cat, color: black,pink');
});

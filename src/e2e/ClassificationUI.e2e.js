/* eslint-disable no-restricted-syntax */
/* eslint-disable prefer-template */
/* eslint-disable consistent-return */
/* eslint-disable func-names */
import { Application } from 'spectron';
import electronPath from 'electron';
import path from 'path';
import { expect } from 'chai';
import { getBtn, assertButtonLabels } from './test_common';
import api from '../api';

const appPath = path.join(__dirname, '../..');
const app = new Application({
  path: electronPath,
  args: [appPath],
});

function sleep(ms) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

async function createProject(xml) {
  const imgDir = path.join(appPath, 'src', 'e2e', 'test_data', 'proj0');
  const glob0 = path.join(imgDir, '*.jpg');
  const glob1 = path.join(imgDir, '*.png');

  await api.post('/project', {
    description: {
      name: 'testproj0',
    },
    template: {
      task: 'classification',
      xml,
    },
    data_sources: [
      { type: 'local_directory', source_uri: glob0 },
      { type: 'local_directory', source_uri: glob1 },
    ],
  });
}

function expectFocusIsOn(device) {
  it(`focused on device ${device}`, async () => {
    const device0 = await app.client.element('//*[@id="device0"]').getAttribute('class');
    const device1 = await app.client.element('//*[@id="device1"]').getAttribute('class');
    const deviceSubmit = await app.client.element('//*[@id="device_submit"]').getAttribute('class');

    const pairs = [
      ['device0', device0],
      ['device1', device1],
      ['device_submit', deviceSubmit],
    ];

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

const getPath = (el, pth) => app.client.elementIdElement(el.ELEMENT, pth);

async function assertCheckBoxLabels(device) {
  const elements = await app.client.elements(`//*[@id="${device}"]/div/ul/li/label/input`);
  const labels = await Promise.all(elements.value.map(async (el) => {
    const kek = await getPath(el, '..');
    const txt = await app.client.elementIdText(kek.value.ELEMENT);
    return txt.value;
  }));
  expect(labels).to.be.deep.eq(['Black', 'White', 'Pink']);
}

describe('Radio + Checkbox fields', function () {
  this.timeout(30000);
  before(() => app.start());
  after(() => {
    if (app && app.isRunning()) {
      return app.stop();
    }
  });

  const xml1 = `
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

  it('navigates to project page', async () => {
    await createProject(xml1);
    await app.client.refresh();
    await app.client.waitUntilTextExists('span', 'testproj0');
    await sleep(1500);
    await app.client.element("button/*[@innertext='testproj0']").click();
  });

  it('begins assess', async () => {
    await app.client.waitUntilTextExists('span', 'Begin assess');

    // fixme
    const beginAssess = app.client.element('//*[@id="grid"]/div/div[2]/div/div/div/div/div[1]/div/div[3]/div/div/button[1]');
    await beginAssess.click();
  });


  it('displays correct devices labels', async () => {
    await assertButtonLabels(app);
    await assertCheckBoxLabels('device1');
  });

  it('selects first device', async () => {
    const device1 = await app.client.element('//*[@id="device0"]').getAttribute('class');
    expect(device1).to.have.string('selected');
    const device2 = await app.client.element('//*[@id="device1"]').getAttribute('class');
    expect(device2).to.not.have.string('selected');
  });

  expectFocusIsOn('device0');

  it('displays button 3 pressed', async () => {
    await app.client.keys('3');
    const disabled = await getDisabled('device0');
    expect(disabled).to.be.deep.eq([false, false, true, false]);
  });

  expectFocusIsOn('device0');

  it('submits radio field', async () => {
    await app.client.keys('Enter');
  });

  expectFocusIsOn('device1');

  it('displays button 1 pressed', async () => {
    await getBtn(app, 1).click();
    const disabled = await getDisabled('device0');
    expect(disabled).to.be.deep.eq([true, false, false, false]);
  });

  expectFocusIsOn('device1');

  it('displays checkbox 2 as checked', async () => {
    await app.client.keys('2');

    const checked = await getChecked('device1');
    expect(checked).to.be.deep.eq([false, true, false]);
  });

  expectFocusIsOn('device1');

  it('displays checkboxes 1,2 as checked', async () => {
    await app.client.keys('1');

    const checked = await getChecked('device1');
    expect(checked).to.be.deep.eq([true, true, false]);
  });

  expectFocusIsOn('device1');

  it('displays all checkboxes as checked', async () => {
    await app.client.keys('3');

    const checked = await getChecked('device1');
    expect(checked).to.be.deep.eq([true, true, true]);
  });

  expectFocusIsOn('device1');

  it('displays checkbox 2 as unchecked', async () => {
    await app.client.keys('2');

    const checked = await getChecked('device1');
    expect(checked).to.be.deep.eq([true, false, true]);
    await app.client.keys('Enter');
  });

  expectFocusIsOn('device_submit');
});

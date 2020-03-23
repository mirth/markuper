/* eslint-disable prefer-template */
/* eslint-disable consistent-return */
/* eslint-disable func-names */
import { Application } from 'spectron';
import electronPath from 'electron';
import path from 'path';
import { expect } from 'chai';
import { assertButtonLabels } from './test_common';

const appPath = path.join(__dirname, '../..');
const app = new Application({
  path: electronPath,
  args: [appPath],
});

function sleep(ms) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}


function createProjectWithTemplate(xml) {
  it('opens Create New Project popup', async () => {
    await app.client.waitForVisible('button');
    await app.client.element('button').click();
    await app.client.waitForVisible('input');
  });

  it('inputs new project name', async () => {
    await sleep(2000);
    await app.client.element('input').setValue('testproj0');
  });

  it('set project task', async () => {
    const template = "input[placeholder='Select task']";
    await app.client.waitForExist(template);
    await app.client.element(template).setValue('classification');
  });

  it('set xml for template', async () => {
    await app.client.elements('textarea').setValue(xml);
  });

  const imgDir = path.join(appPath, 'src', 'e2e', 'test_data', 'proj0');
  const glob0 = path.join(imgDir, '*.jpg');
  const glob1 = path.join(imgDir, '*.png');

  const srcInput = "input[placeholder='/some/path or /some/glob/*.jpg']";

  it('adds first data source', async () => {
    await app.client.element(srcInput).setValue(glob0);
    await app.client.element('button=Add source').click();
    await app.client.waitForVisible('//ul/li/div/input');
    const inputValue = await app.client.element('//ul/li/div/input').getValue();
    expect(inputValue).to.be.eq(glob0);
  });

  it('adds second data source', async () => {
    await app.client.waitForVisible(srcInput);
    const input = app.client.element(srcInput);
    await input.setValue(glob1);
    await app.client.element('button=Add source').click();
    await app.client.waitForVisible('//ul/li[2]/div/input');
    const inputValue = await app.client.element('//ul/li[2]/div/input').getValue();
    expect(inputValue).to.be.eq(glob1);
  });

  it('creates project', async () => {
    await app.client.element('button=Create').click();
  });
}

describe('Application launch', function () {
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

  createProjectWithTemplate(xml1);

  it('navigates to project page', async () => {
    await sleep(500);
    await app.client.waitUntilTextExists('span', 'testproj0');
    await app.client.element("button/*[@innertext='testproj0']").click();
  });

  it('begins assess', async () => {
    await app.client.waitUntilTextExists('span', 'Begin assess');

    // fixme
    const beginAssess = app.client.element('//*[@id="grid"]/div/div[2]/div/div/div/div/div[1]/div/div[3]/div/div/button[1]');
    await beginAssess.click();
  });

  const getPath = (el, pth) => app.client.elementIdElement(el.ELEMENT, pth);

  it('displays correct devices labels', async () => {
    await assertButtonLabels(app);


    const elements = await app.client.elements('//*[@id="device1"]/div/ul/li/label/input');
    const labels = await Promise.all(elements.value.map(async (el) => {
      const kek = await getPath(el, '..');
      const txt = await app.client.elementIdText(kek.value.ELEMENT);
      return txt.value;
    }));
    expect(labels).to.be.deep.eq(['Black', 'White', 'Pink']);
  });

  it('displays correct selected device', async () => {
    const device1 = await app.client.element('//*[@id="device0"]').getAttribute('class');
    expect(device1).to.have.string('selected');
    const device2 = await app.client.element('//*[@id="device1"]').getAttribute('class');
    expect(device2).to.not.have.string('selected');
  });
});

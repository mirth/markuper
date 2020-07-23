/* eslint-disable no-restricted-syntax */
/* eslint-disable prefer-template */
/* eslint-disable consistent-return */
/* eslint-disable func-names */
import { Application } from 'spectron';
import electronPath from 'electron';
import path from 'path';
import { expect } from 'chai';
import {
  radioClick, itNavigatesToProject, getSamplePath,
  clickButton, getRadioState, getChecked, expectSampleMarkupToBeEq,
} from './test_common';
import { TestCheckboxRadioRadio, TestRadioCheckbox, itSubmitsSample } from './classification_common';

const appPath = path.join(__dirname, '../..');
const app = new Application({
  path: electronPath,
  args: [appPath],
});

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
    await radioClick(app, 'root/0', 2);

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

  itSubmitsSample(app);

  it('goes to next new sample', async () => {
    await app.client.waitUntilTextExists('span', 'Begin assess');
    await clickButton(app, 'span', 'Begin assess');

    const selected = await getRadioState(app, 'root/0');
    expect(selected).to.be.deep.eq([false, false]);
    const checked = await getChecked(app, 'root/1');
    expect(checked).to.be.deep.eq([false, false]);
  });

  it('goes back to first assessed sample and checks it has valid markup', async () => {
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

  TestCheckboxRadioRadio(app, 'root');

  itSubmitsSample(app);
  expectSampleMarkupToBeEq(app, appPath, { animal: 'gek', color: ['black'], size: 'big' });
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

  itNavigatesToProject(app, appPath, xml);

  it('begins assess', async () => {
    await app.client.waitUntilTextExists('span', 'Begin assess');
    await clickButton(app, 'span', 'Begin assess');
  });


  TestRadioCheckbox(app, 'root');

  itSubmitsSample(app);
  expectSampleMarkupToBeEq(app, appPath, { animal: 'cat', color: ['black', 'pink'] });
});

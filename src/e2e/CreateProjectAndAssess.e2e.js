/* eslint-disable prefer-template */
/* eslint-disable consistent-return */
/* eslint-disable func-names */
import { Application } from 'spectron';
import electronPath from 'electron';
import path from 'path';
import { expect } from 'chai';

const appPath = path.join(__dirname, '../..');
const app = new Application({
  path: electronPath,
  args: [appPath],
});

function sleep(ms) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}
// pause
describe('Application launch', function () {
  this.timeout(20000);
  before(() => app.start());
  after(() => {
    if (app && app.isRunning()) {
      return app.stop();
    }
  });

  describe('Creates project and assesses samples', () => {
    it('opens Create New Project popup', async () => {
      await app.client.waitForVisible('button');
      await app.client.element('button').click();
      await app.client.waitForVisible('input');
    });

    it('inputs new project name', async () => {
      await sleep(300);
      await app.client.element('input').setValue('testproj0');
    });

    it('set project task', async () => {
      const template = "input[placeholder='Select task']";
      await app.client.waitForExist(template);
      await app.client.element(template).setValue('classification');
    });

    it('adds two new labels for classification template', async () => {
      const newLabelInputSelector = "input[placeholder='Label goes here...'";
      const getNewLabelInput = () => app.client.element(newLabelInputSelector);
      const addLabelSelector = '//div[2]/button';
      const addLabel = () => app.client.element(addLabelSelector);
      await getNewLabelInput().setValue('chuk');
      await app.client.waitForVisible(addLabelSelector);
      await addLabel().click();
      await app.client.waitForExist(newLabelInputSelector);
      await getNewLabelInput().setValue('gek');
      await addLabel().click();
    });

    const imgDir = path.join(appPath, 'src', 'e2e', 'test_data', 'proj0');
    const glob0 = path.join(imgDir, '*.jpg');
    const glob1 = path.join(imgDir, '*.png');

    const srcInput = "input[placeholder='/some/path']";

    it('adds first data source', async () => {
      await app.client.element(srcInput).setValue(glob0);
      await app.client.element('button=Add source').click();
      await app.client.waitForVisible('//ul[2]/li/div/input');
      const inputValue = await app.client.element('//ul[2]/li/div/input').getValue();
      expect(inputValue).to.be.eq(glob0);
    });

    it('adds second data source', async () => {
      await app.client.waitForVisible(srcInput);
      const input = app.client.element(srcInput);
      await input.setValue(glob1);
      await app.client.element('button=Add source').click();
      await app.client.waitForVisible('//ul[2]/li[2]/div/input');
      const inputValue = await app.client.element('//ul[2]/li[2]/div/input').getValue();
      expect(inputValue).to.be.eq(glob1);
    });

    it('creates project', async () => {
      await app.client.element('button=Create').click();
    });

    it('navigates to project page', async () => {
      await app.client.waitUntilTextExists('a', 'testproj0');
      app.client.refresh(); // fixme remove when router will be fixed
      await app.client.waitUntilTextExists('a', 'testproj0');
      await app.client.element('=testproj0').click();
    });

    it('displays project description', async () => {
      await app.client.waitUntilTextExists('span', 'testproj0');
      await app.client.waitUntilTextExists('b', 'classification');
      await app.client.waitUntilTextExists('span', glob0);
      await app.client.waitUntilTextExists('span', glob1);
      await app.client.waitUntilTextExists('span', 'cat, dog, chuk, gek');
    });

    it('begins assess', async () => {
      await app.client.waitUntilTextExists('span', 'Begin assess');
      const beginAssess = app.client.element("button/*[@innertext='Begin assess']");
      await beginAssess.click();
    });

    const makeUrl = (filename) => path.normalize('file://' + path.join(imgDir, filename));

    it('assesses 1st jpg sample', async () => {
      await app.client.waitForExist('img');
      const src = await app.client.element('img').getAttribute('src');
      expect(path.normalize(src)).to.be.eq(makeUrl('kek0.jpg'));
      await app.client.element('button=chuk').click();
    });

    it('assesses 2nd jpg sample', async () => {
      await app.client.waitForExist('img');
      const src = await app.client.element('img').getAttribute('src');
      expect(path.normalize(src)).to.be.eq(makeUrl('kek1.jpg'));
      await app.client.element('button=dog').click();
    });

    it('assesses 3rd jpg sample', async () => {
      await app.client.waitForExist('img');
      const src = await app.client.element('img').getAttribute('src');
      expect(path.normalize(src)).to.be.eq(makeUrl('kek2.jpg'));
      await app.client.element('button=dog').click();
    });

    it('assesses 1st png sample', async () => {
      await app.client.waitForExist('img');
      const src = await app.client.element('img').getAttribute('src');
      expect(path.normalize(src)).to.be.eq(makeUrl('kek3.png'));
      await app.client.element('button=cat').click();
    });

    it('assesses 2nd png sample', async () => {
      await app.client.waitForExist('img');
      const src = await app.client.element('img').getAttribute('src');
      expect(path.normalize(src)).to.be.eq(makeUrl('kek4.png'));
    });

    it('displays sample markup on project page', async () => {
      await app.client.element('a=testproj0').click();
      await app.client.waitForExist('ul');
      const getMarkupFor = (sampleID) => app.client.element(`p*=Sample ID: ${sampleID}`);
      let text = await getMarkupFor(0).getText();
      expect(text).to.be.eq('Sample ID: 0|Value: class:chuk');
      text = await getMarkupFor(1).getText();
      expect(text).to.be.eq('Sample ID: 1|Value: class:dog');
      text = await getMarkupFor(2).getText();
      expect(text).to.be.eq('Sample ID: 2|Value: class:dog');
      text = await getMarkupFor(3).getText();
      expect(text).to.be.eq('Sample ID: 3|Value: class:cat');
    });
  });
});

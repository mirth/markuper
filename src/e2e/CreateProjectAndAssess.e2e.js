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
  this.timeout(30000);
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
      await sleep(2000);
      await app.client.element('input').setValue('testproj0');
    });

    it('set project task', async () => {
      const template = "input[placeholder='Select task']";
      await app.client.waitForExist(template);
      await app.client.element(template).setValue('classification');
    });

    it('adds four new labels for classification template', async () => {
      const newLabelInputSelector = "input[placeholder='Label goes here...'";
      const getNewLabelInput = () => app.client.element(newLabelInputSelector);
      const addLabelSelector = '//div[2]/button';
      const addLabel = () => app.client.element(addLabelSelector);
      await getNewLabelInput().setValue('chuk');
      await app.client.waitForVisible(addLabelSelector);
      await addLabel().click();
      await app.client.waitForExist(newLabelInputSelector);
      await getNewLabelInput().setValue('azaza');
      await addLabel().click();
      await app.client.waitForExist(newLabelInputSelector);
      await getNewLabelInput().setValue('gek');
      await addLabel().click();
      await app.client.waitForExist(newLabelInputSelector);
      await getNewLabelInput().setValue('lul');
      await addLabel().click();
    });

    it("remove 'azaza' and 'lul' lables", async () => {
      const getNthCross = (nth) => `//li[${nth}]/div/div/button`;
      await app.client.waitForVisible(getNthCross(4));
      await app.client.element(getNthCross(4)).click();
      await app.client.element(getNthCross(5)).click();
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

    it('navigates to project page', async () => {
      await app.client.waitUntilTextExists('span', 'testproj0');
      await app.client.element("button/*[@innertext='testproj0'").click();
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

      // fixme
      const beginAssess = app.client.element('//*[@id="grid"]/div/div[2]/div/div/div/div/div[1]/div/div[3]/div/div/button[1]');
      await beginAssess.click();
    });

    const makeUrl = (filename) => path.normalize('file://' + path.join(imgDir, filename));

    it('assesses 1st jpg sample', async () => {
      await app.client.waitForExist('img');
      const src = await app.client.element('img').getAttribute('src');
      expect(path.normalize(src)).to.be.eq(makeUrl('kek0.jpg'));
      await app.client.keys('3');
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
      await app.client.keys('2');
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

    const getPath = (filename) => app.client.element(`small*=${filename}`);
    const getClass = (filename) => getPath(filename).element('../..').element('./span');

    // fixme test sample order
    it('displays sample markup on project page', async () => {
      await app.client.element("button/*[@innertext='testproj0'").click();
      await app.client.waitForExist('ul');

      {
        const pathText = await getPath('kek0.jpg').getText();
        const cl = await getClass('kek0.jpg').getText();
        expect(pathText).to.be.eq(path.join(imgDir, 'kek0.jpg') + ':');
        expect(cl).to.be.eq('class: chuk');
      }

      {
        const pathText = await getPath('kek1.jpg').getText();
        const cl = await getClass('kek1.jpg').getText();
        expect(pathText).to.be.eq(path.join(imgDir, 'kek1.jpg') + ':');
        expect(cl).to.be.eq('class: dog');
      }

      {
        const pathText = await getPath('kek2.jpg').getText();
        const cl = await getClass('kek2.jpg').getText();
        expect(pathText).to.be.eq(path.join(imgDir, 'kek2.jpg') + ':');
        expect(cl).to.be.eq('class: dog');
      }

      {
        const pathText = await getPath('kek3.png').getText();
        const cl = await getClass('kek3.png').getText();
        expect(pathText).to.be.eq(path.join(imgDir, 'kek3.png') + ':');
        expect(cl).to.be.eq('class: cat');
      }
    });

    it('contains correct pressed button', async () => {
      await getPath('kek1.jpg').element('../..').click();
      await app.client.waitForVisible("button/*[@innertext='cat'");
      await app.client.refresh();
      await app.client.waitForVisible("button/*[@innertext='cat'");

      const getBtn = (i) => app.client.element('//*[@id="grid"]').element(`./div/div[2]/div/div/div/ul/li[${i}]/div/button`);
      const catCl = await getBtn(1).getAttribute('class');
      const dogCl = await getBtn(2).getAttribute('class');
      const chukCl = await getBtn(3).getAttribute('class');
      const gekCl = await getBtn(4).getAttribute('class');

      expect(catCl).not.to.include('disabled');
      expect(dogCl).to.include('disabled');
      expect(chukCl).not.to.include('disabled');
      expect(gekCl).not.to.include('disabled');
    });
  });
});

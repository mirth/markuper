/* eslint-disable prefer-template */
/* eslint-disable consistent-return */
/* eslint-disable func-names */
import { Application } from 'spectron';
import electronPath from 'electron';
import path from 'path';
import { expect } from 'chai';
import {
  makeUrl, getBtn, assertRadioLabels, getSamplePath, getSampleClass, createProjectWithTemplate,
  sleep,
} from './test_common';


const appPath = path.join(__dirname, '../..');
const app = new Application({
  path: electronPath,
  args: [appPath],
});

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
    const xml = `
    <content>
      <radio group="animal" value="cat" vizual="Cat" />
      <radio group="animal" value="dog" vizual="Dog" />
      <radio group="animal" value="chuk" vizual="Chuk" />
      <radio group="animal" value="gek" vizual="Gek" />
    </content>
    `;
    const [imgDir, glob0, glob1] = createProjectWithTemplate(app, appPath, xml);

    it('navigates to project page', async () => {
      await sleep(500);
      await app.client.waitUntilTextExists('span', 'testproj0');
      await app.client.element("button/*[@innertext='testproj0']").click();
    });

    it('displays project description', async () => {
      await app.client.waitUntilTextExists('span', 'testproj0');
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

    it('assesses 1st jpg sample', async () => {
      await app.client.waitForExist('img');
      await assertRadioLabels(app, 'device0', ['Cat', 'Dog', 'Chuk', 'Gek']);
      const src = await app.client.element('img').getAttribute('src');
      expect(path.normalize(src)).to.be.eq(makeUrl(imgDir, 'kek0.jpg'));
      await app.client.keys('3');
      await app.client.keys('Enter');
      await app.client.keys('Enter');
    });

    it('assesses 2nd jpg sample', async () => {
      await app.client.waitForExist('img');
      await assertRadioLabels(app, 'device0', ['Cat', 'Dog', 'Chuk', 'Gek']);
      const src = await app.client.element('img').getAttribute('src');
      expect(path.normalize(src)).to.be.eq(makeUrl(imgDir, 'kek1.jpg'));
      await getBtn(app, 'device0', 2).click();
      await app.client.keys('Enter');
      await app.client.keys('Enter');
    });

    it('assesses 3rd jpg sample', async () => {
      await app.client.waitForExist('img');
      await assertRadioLabels(app, 'device0', ['Cat', 'Dog', 'Chuk', 'Gek']);
      const src = await app.client.element('img').getAttribute('src');
      expect(path.normalize(src)).to.be.eq(makeUrl(imgDir, 'kek2.jpg'));
      await app.client.keys('2');
      await app.client.keys('Enter');
      await app.client.keys('Enter');
    });

    it('assesses 1st png sample', async () => {
      await app.client.waitForExist('img');
      await assertRadioLabels(app, 'device0', ['Cat', 'Dog', 'Chuk', 'Gek']);
      const src = await app.client.element('img').getAttribute('src');
      expect(path.normalize(src)).to.be.eq(makeUrl(imgDir, 'kek3.png'));
      await getBtn(app, 'device0', 1).click();
      await app.client.keys('Enter');
      await app.client.keys('Enter');
    });

    it('assesses 2nd png sample', async () => {
      await app.client.waitForExist('img');
      await assertRadioLabels(app, 'device0', ['Cat', 'Dog', 'Chuk', 'Gek']);
      const src = await app.client.element('img').getAttribute('src');
      expect(path.normalize(src)).to.be.eq(makeUrl(imgDir, 'kek4.png'));
      await getBtn(app, 'device0', 1).click();
      await app.client.keys('Enter');
      await app.client.keys('Enter');
    });

    it('displays that there are no samples left', async () => {
      await app.client.waitUntilTextExists('span', 'No samples left');
    });

    // fixme test sample order
    it('displays sample markup on project page', async () => {
      await app.client.element("button/*[@innertext='testproj0']").click();
      await app.client.waitForExist('ul');

      {
        const pathText = await getSamplePath(app, 'kek0.jpg').getText();
        const cl = await getSampleClass(app, 'kek0.jpg').getText();
        expect(pathText).to.be.eq(path.join(imgDir, 'kek0.jpg') + ':');
        expect(cl).to.be.eq('animal: chuk');
      }

      {
        const pathText = await getSamplePath(app, 'kek1.jpg').getText();
        const cl = await getSampleClass(app, 'kek1.jpg').getText();
        expect(pathText).to.be.eq(path.join(imgDir, 'kek1.jpg') + ':');
        expect(cl).to.be.eq('animal: dog');
      }

      {
        const pathText = await getSamplePath(app, 'kek2.jpg').getText();
        const cl = await getSampleClass(app, 'kek2.jpg').getText();
        expect(pathText).to.be.eq(path.join(imgDir, 'kek2.jpg') + ':');
        expect(cl).to.be.eq('animal: dog');
      }

      {
        const pathText = await getSamplePath(app, 'kek3.png').getText();
        const cl = await getSampleClass(app, 'kek3.png').getText();
        expect(pathText).to.be.eq(path.join(imgDir, 'kek3.png') + ':');
        expect(cl).to.be.eq('animal: cat');
      }
      {
        const pathText = await getSamplePath(app, 'kek4.png').getText();
        const cl = await getSampleClass(app, 'kek4.png').getText();
        expect(pathText).to.be.eq(path.join(imgDir, 'kek4.png') + ':');
        expect(cl).to.be.eq('animal: cat');
      }
    });

    it('contains correct pressed button', async () => {
      await getSamplePath(app, 'kek1.jpg').element('..').click();

      await app.client.waitForVisible("button/*[@innertext='Cat']");
      const catCl = await getBtn(app, 'device0', 1).getAttribute('class');
      const dogCl = await getBtn(app, 'device0', 2).getAttribute('class');
      const chukCl = await getBtn(app, 'device0', 3).getAttribute('class');
      const gekCl = await getBtn(app, 'device0', 4).getAttribute('class');

      expect(catCl).not.to.include('disabled');
      expect(dogCl).to.include('disabled');
      expect(chukCl).not.to.include('disabled');
      expect(gekCl).not.to.include('disabled');
    });

    it("changes class to 'gek'", async () => {
      await getBtn(app, 'device0', 4).click();
      await sleep(1500);
      await app.client.keys('Enter');
      await app.client.keys('Enter');
      await sleep(1500);
      await app.client.element("button/*[@innertext='testproj0']").click();
      await app.client.waitForExist('ul');
      await sleep(1500);
      {
        const cl = await getSampleClass(app, 'kek1.jpg').getText();
        expect(cl).to.be.eq('animal: gek');
      }
    });
  });
});

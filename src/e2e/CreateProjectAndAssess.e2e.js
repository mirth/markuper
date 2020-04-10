/* eslint-disable prefer-template */
/* eslint-disable consistent-return */
/* eslint-disable func-names */
import { Application } from 'spectron';
import electronPath from 'electron';
import path from 'path';
import { expect } from 'chai';
import {
  makeUrl, getRadio, assertRadioLabels, getSamplePath, getSampleClass, createProjectWithTemplate,
  sleep, clickButton, getRadioState,
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

    it('displays project description', async () => {
      await app.client.waitUntilTextExists('span', 'testproj0');
      await app.client.waitUntilTextExists('span', glob0);
      await app.client.waitUntilTextExists('span', glob1);
      await app.client.waitUntilTextExists('span', 'cat, dog, chuk, gek');
    });

    it('begins assess', async () => {
      await app.client.waitUntilTextExists('span', 'Begin assess');
      await clickButton(app, 'span', 'Begin assess');
    });

    it('assesses 1st jpg sample', async () => {
      await app.client.waitForExist('img');
      await assertRadioLabels(app, 'root/0', ['Cat', 'Dog', 'Chuk', 'Gek']);
      const src = await app.client.element('img').getAttribute('src');
      expect(path.normalize(src)).to.be.eq(makeUrl(imgDir, 'kek0.jpg'));
      await app.client.keys('3');
      await app.client.keys('Enter');
      await app.client.keys('Enter');
    });

    it('assesses 2nd jpg sample', async () => {
      await app.client.waitForExist('img');
      await assertRadioLabels(app, 'root/0', ['Cat', 'Dog', 'Chuk', 'Gek']);
      const src = await app.client.element('img').getAttribute('src');
      expect(path.normalize(src)).to.be.eq(makeUrl(imgDir, 'kek1.jpg'));
      await getRadio(app, 'root/0', 2).click();
      await app.client.keys('Enter');
      await app.client.keys('Enter');
    });

    it('assesses 3rd jpg sample', async () => {
      await app.client.waitForExist('img');
      await assertRadioLabels(app, 'root/0', ['Cat', 'Dog', 'Chuk', 'Gek']);
      const src = await app.client.element('img').getAttribute('src');
      expect(path.normalize(src)).to.be.eq(makeUrl(imgDir, 'kek2.jpg'));
      await app.client.keys('2');
      await app.client.keys('Enter');
      await app.client.keys('Enter');
    });

    it('assesses 1st png sample', async () => {
      await app.client.waitForExist('img');
      await assertRadioLabels(app, 'root/0', ['Cat', 'Dog', 'Chuk', 'Gek']);
      const src = await app.client.element('img').getAttribute('src');
      expect(path.normalize(src)).to.be.eq(makeUrl(imgDir, 'kek3.png'));
      await getRadio(app, 'root/0', 1).click();
      await app.client.keys('Enter');
      await app.client.keys('Enter');
    });

    it('assesses 2nd png sample', async () => {
      await app.client.waitForExist('img');
      await assertRadioLabels(app, 'root/0', ['Cat', 'Dog', 'Chuk', 'Gek']);
      const src = await app.client.element('img').getAttribute('src');
      expect(path.normalize(src)).to.be.eq(makeUrl(imgDir, 'kek4.png'));
      await getRadio(app, 'root/0', 1).click();
      await app.client.keys('Enter');
      await app.client.keys('Enter');
    });

    it('displays that there are no samples left', async () => {
      await app.client.waitUntilTextExists('span', 'No samples left');
    });

    // fixme test sample order
    it('displays sample markup on project page', async () => {
      await clickButton(app, 'span', 'testproj0');
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

    it('contains correct selected radios', async () => {
      await getSamplePath(app, 'kek1.jpg').element('..').click();
      await app.client.waitForVisible("button/*[@innertext='Cat']");

      const selected = await getRadioState(app, 'root/0');
      expect(selected).to.be.deep.eq([false, true, false, false]);
    });

    it("changes class to 'gek'", async () => {
      await getRadio(app, 'root/0', 4).click();
      await sleep(1500);
      await app.client.keys('Enter');
      await app.client.keys('Enter');
      await sleep(1500);
      await clickButton(app, 'span', 'testproj0');
      await app.client.waitForExist('ul');
      await sleep(1500);
      {
        const cl = await getSampleClass(app, 'kek1.jpg').getText();
        expect(cl).to.be.eq('animal: gek');
      }
    });
  });
});

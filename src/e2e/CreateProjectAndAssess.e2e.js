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

    it('set xml for template', async () => {
      await app.client.elements('textarea').setValue(`
      <content>
        <radio group="animal" value="cat" vizual="Cat" />
        <radio group="animal" value="dog" vizual="Dog" />
        <radio group="animal" value="chuk" vizual="Chuk" />
        <radio group="animal" value="gek" vizual="Gek" />
      </content>
      `);
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

    const makeUrl = (filename) => path.normalize('file://' + path.join(imgDir, filename));
    const getBtn = (i) => app.client.element('//*[@id="grid"]').element(`./div/div[2]/div/div/div/div[1]/ul/li[${i}]/div/button`);// `./div/div[2]/div/div/div/ul/li[${i}]/div/button`);

    const assertButtonLabels = async () => {
      let btnTxt = await getBtn(1).element('.//span').getText();
      expect(btnTxt).to.be.eq('Cat');
      btnTxt = await getBtn(2).element('.//span').getText();
      expect(btnTxt).to.be.eq('Dog');
      btnTxt = await getBtn(3).element('.//span').getText();
      expect(btnTxt).to.be.eq('Chuk');
      btnTxt = await getBtn(4).element('.//span').getText();
      expect(btnTxt).to.be.eq('Gek');
    };

    it('assesses 1st jpg sample', async () => {
      await app.client.waitForExist('img');
      await assertButtonLabels();
      const src = await app.client.element('img').getAttribute('src');
      expect(path.normalize(src)).to.be.eq(makeUrl('kek0.jpg'));
      await app.client.keys('3');
      await app.client.keys('Enter');
    });

    it('assesses 2nd jpg sample', async () => {
      await app.client.waitForExist('img');
      await assertButtonLabels();
      const src = await app.client.element('img').getAttribute('src');
      expect(path.normalize(src)).to.be.eq(makeUrl('kek1.jpg'));
      await getBtn(2).click();
      await app.client.keys('Enter');
    });

    it('assesses 3rd jpg sample', async () => {
      await app.client.waitForExist('img');
      await assertButtonLabels();
      const src = await app.client.element('img').getAttribute('src');
      expect(path.normalize(src)).to.be.eq(makeUrl('kek2.jpg'));
      await app.client.keys('2');
      await app.client.keys('Enter');
    });

    it('assesses 1st png sample', async () => {
      await app.client.waitForExist('img');
      await assertButtonLabels();
      const src = await app.client.element('img').getAttribute('src');
      expect(path.normalize(src)).to.be.eq(makeUrl('kek3.png'));
      await getBtn(1).click();
      await app.client.keys('Enter');
    });

    it('assesses 2nd png sample', async () => {
      await app.client.waitForExist('img');
      await assertButtonLabels();
      const src = await app.client.element('img').getAttribute('src');
      expect(path.normalize(src)).to.be.eq(makeUrl('kek4.png'));
      await getBtn(1).click();
      await app.client.keys('Enter');
    });

    it('displays that there are no samples left', async () => {
      await app.client.waitUntilTextExists('span', 'No samples left');
    });

    const getPath = (filename) => app.client.element(`small*=${filename}`);
    const getClass = (filename) => getPath(filename).element('../..').element('./span');

    // fixme test sample order
    it('displays sample markup on project page', async () => {
      await app.client.element("button/*[@innertext='testproj0']").click();
      await app.client.waitForExist('ul');

      {
        // console.log("getPath('kek0.jpg'): ", getPath('kek0.jpg'))
        const pathText = await getPath('kek0.jpg').getText();
        const cl = await getClass('kek0.jpg').getText();
        expect(pathText).to.be.eq(path.join(imgDir, 'kek0.jpg') + ':');
        expect(cl).to.be.eq('animal: chuk');
      }

      {
        const pathText = await getPath('kek1.jpg').getText();
        const cl = await getClass('kek1.jpg').getText();
        expect(pathText).to.be.eq(path.join(imgDir, 'kek1.jpg') + ':');
        expect(cl).to.be.eq('animal: dog');
      }

      {
        const pathText = await getPath('kek2.jpg').getText();
        const cl = await getClass('kek2.jpg').getText();
        expect(pathText).to.be.eq(path.join(imgDir, 'kek2.jpg') + ':');
        expect(cl).to.be.eq('animal: dog');
      }

      {
        const pathText = await getPath('kek3.png').getText();
        const cl = await getClass('kek3.png').getText();
        expect(pathText).to.be.eq(path.join(imgDir, 'kek3.png') + ':');
        expect(cl).to.be.eq('animal: cat');
      }
      {
        const pathText = await getPath('kek4.png').getText();
        const cl = await getClass('kek4.png').getText();
        expect(pathText).to.be.eq(path.join(imgDir, 'kek4.png') + ':');
        expect(cl).to.be.eq('animal: cat');
      }
    });

    it('contains correct pressed button', async () => {
      await getPath('kek1.jpg').element('..').click();

      await app.client.waitForVisible("button/*[@innertext='Cat']");
      const catCl = await getBtn(1).getAttribute('class');
      const dogCl = await getBtn(2).getAttribute('class');
      const chukCl = await getBtn(3).getAttribute('class');
      const gekCl = await getBtn(4).getAttribute('class');

      expect(catCl).not.to.include('disabled');
      expect(dogCl).to.include('disabled');
      expect(chukCl).not.to.include('disabled');
      expect(gekCl).not.to.include('disabled');
    });

    it("changes class to 'gek'", async () => {
      await getBtn(4).click();
      await app.client.keys('Enter');
      await app.client.element("button/*[@innertext='testproj0']").click();
      await app.client.waitForExist('ul');
      {
        const cl = await getClass('kek1.jpg').getText();
        expect(cl).to.be.eq('animal: gek');
      }
    });
  });
});

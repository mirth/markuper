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


describe('Application launch', function () {
  this.timeout(10000);
  beforeEach(() => app.start());
  afterEach(() => {
    if (app && app.isRunning()) {
      return app.stop();
    }
  });

  it('Creates project and assesses samples', async () => {
    await app.client.waitForVisible('button');
    await app.client.element('button').click();
    await app.client.waitForVisible('input')
    await app.client.element('input').setValue('testproj0');

    await app.client.element("input[value='classification'").click();

    const newLabelInputSelector = "input[placeholder='Label goes here...'";
    const getNewLabelInput = () => app.client.element(newLabelInputSelector);
    const addLabel = () => app.client.element('button=+');
    await getNewLabelInput().setValue('chuk');
    await addLabel().click();
    await app.client.waitForExist(newLabelInputSelector);
    await getNewLabelInput().setValue('gek');
    await addLabel().click();

    const imgDir = path.join(appPath, 'src', 'e2e', 'test_data', 'proj0');
    const glob = path.join(imgDir, '*.jpg');
    await app.client.element("input[placeholder='/some/path']").setValue(glob);
    await app.client.element('button=Create').click();

    await app.client.waitUntilTextExists('a', 'testproj0');
    app.client.refresh(); // fixme remove when router will be fixed
    await app.client.waitUntilTextExists('a', 'testproj0');
    await app.client.element('=testproj0').click();

    await app.client.waitUntilTextExists('span', 'testproj0');
    await app.client.waitUntilTextExists('b', 'classification');
    await app.client.waitUntilTextExists('span', path.join(imgDir, '*.jpg'));
    await app.client.waitUntilTextExists('span', 'cat, dog, chuk, gek');

    await app.client.waitUntilTextExists('span', 'Begin assess');
    const beginAssess = app.client.element("button/*[@innertext='Begin assess']");
    await beginAssess.click();

    await app.client.waitForExist('img');
    const before = await app.client.element('img').getAttribute('src');
    const makeUrl = (filename) => path.normalize('file://' + path.join(imgDir, filename));

    expect(path.normalize(before)).to.be.eq(makeUrl('kek0.jpg'));

    await app.client.element('button=chuk').click();

    await app.client.waitForExist('img');
    const after = await app.client.element('img').getAttribute('src');
    expect(path.normalize(after)).to.be.eq(makeUrl('kek1.jpg'));
    expect(before).to.be.not.eq(after);

    await app.client.element('button=dog').click();

    await app.client.element('a=testproj0').click();

    await app.client.waitForExist('ul');

    const getMarkupFor = (sampleID) => app.client.element(`p*=Sample ID: ${sampleID}`);
    let text = await getMarkupFor(0).getText();
    expect(text).to.be.eq('Sample ID: 0|Value: chuk');
    text = await getMarkupFor(1).getText();
    expect(text).to.be.eq('Sample ID: 1|Value: dog');
  });
});

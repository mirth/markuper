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
    await app.client.waitForExist('button');
    await app.client.element('button=Create new project').click();
    await app.client.waitUntilTextExists('button', 'Create');
    await app.client.element('input').setValue('testproj0');

    await app.client.element("input[value='classification'").click();

    const imgDir = path.join(appPath, 'src', 'e2e', 'test_data', 'proj0');
    const glob = path.join(imgDir, '*.jpg');
    await app.client.element("input[placeholder='/some/path']").setValue(glob);
    await app.client.element('button=Create').click();

    await app.client.waitUntilTextExists('a', 'testproj0');
    app.client.refresh(); // fixme remove when router will be fixed
    await app.client.waitUntilTextExists('a', 'testproj0');

    await app.client.element('=testproj0').click();
    await app.client.waitUntilTextExists('a', 'Begin assess');
    await app.client.element('=Begin assess').click();

    await app.client.waitForExist('img');
    const before = await app.client.element('img').getAttribute('src');
    const makeUrl = (filename) => 'file://' + path.join(imgDir, filename);

    expect(before).to.be.eq(makeUrl('kek0.jpg'));

    await app.client.element('button').click();

    await app.client.waitForExist('img');
    const after = await app.client.element('img').getAttribute('src');
    expect(after).to.be.eq(makeUrl('kek1.jpg'));
    expect(before).to.be.not.eq(after);
  });
});

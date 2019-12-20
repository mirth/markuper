/* eslint-disable consistent-return */
/* eslint-disable func-names */
import { Application } from 'spectron';
import electronPath from 'electron';
import path from 'path';
import { expect } from 'chai';

const appPath = [path.join(__dirname, '../..')];

const app = new Application({
  path: electronPath,
  args: appPath,
});

describe('Application launch', function () {
  this.timeout(10000);

  beforeEach(() => app.start());

  afterEach(() => {
    if (app && app.isRunning()) {
      return app.stop();
    }
  });

  it('shows next sample after we assess current', async () => {
    await app.client.waitForExist('img');

    const before = await app.client.element('img').getAttribute('src');
    expect(before).to.be.eq('file://img0/');

    await app.client.element('button').click();

    await app.client.waitForExist('img');
    const after = await app.client.element('img').getAttribute('src');
    expect(after).to.be.eq('file://img1/');
    expect(before).to.be.not.eq(after);
  });
});

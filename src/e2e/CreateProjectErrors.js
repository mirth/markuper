/* eslint-disable prefer-template */
/* eslint-disable consistent-return */
/* eslint-disable func-names */
import { Application } from 'spectron';
import electronPath from 'electron';
import path from 'path';
import { expect } from 'chai';
import { createProjectWithTemplate } from './test_common';

const appPath = path.join(__dirname, '../..');
const app = new Application({
  path: electronPath,
  args: [appPath],
});

describe('Application launch', function () {
  this.timeout(30000);
  before(() => app.start());
  after(() => {
    if (app && app.isRunning()) {
      return app.stop();
    }
  });

  describe('Unable to create project because of duplicated groups', () => {
    const xml = `
    <content>
      <bounding_box group="box">
        <checkbox group="color" value="black" vizual="Black" />
        <checkbox group="color" value="white" vizual="White" />
        <checkbox group="color" value="pink" vizual="Pink" />

        <checkbox group="box" value="cat" vizual="Cat" />
        <checkbox group="box" value="dog" vizual="Dog" />
      </bounding_box>
    </content>
    `;
    createProjectWithTemplate(app, appPath, xml);

    it('display duplicate group error', async () => {
      await app.client.waitForVisible('//*[@id="create_project_error"]');
      const el = app.client.element('//*[@id="create_project_error"]');
      const err = await el.getText();
      expect(err).to.be.eq('Template has duplicate groups: box');
    });
  });
});

describe('Application launch', function () {
  this.timeout(30000);
  before(() => app.start());
  after(() => {
    if (app && app.isRunning()) {
      return app.stop();
    }
  });

  describe('Unable to create project because of some attribute is empty', () => {
    const xml = `
    <content>
      <radio group="animal" value="cat" vizual="Cat" />
      <checkbox group="animal" value="" vizual="Dog" />
    </content>
    `;
    createProjectWithTemplate(app, appPath, xml);

    it('display duplicate group error', async () => {
      await app.client.waitForVisible('//*[@id="create_project_error"]');
      const el = app.client.element('//*[@id="create_project_error"]');
      const err = await el.getText();
      expect(err).to.be.eq('Element [checkbox] has an empty attribute [value]');
    });
  });
});


describe('Application launch', function () {
  this.timeout(30000);
  before(() => app.start());
  after(() => {
    if (app && app.isRunning()) {
      return app.stop();
    }
  });

  describe('Unable to create project because of duplicated labels', () => {
    const xml = `
    <content>
      <radio group="animal" value="cat" vizual="Cat" />
      <radio group="animal" value="cat" vizual="Cot" />
    </content>
    `;
    createProjectWithTemplate(app, appPath, xml);

    it('display duplicate labels error', async () => {
      await app.client.waitForVisible('//*[@id="create_project_error"]');
      const el = app.client.element('//*[@id="create_project_error"]');
      const err = await el.getText();
      expect(err).to.be.eq('Template has duplicate labels: group [animal] labels [cat]');
    });
  });
});

describe('Application launch', function () {
  this.timeout(30000);
  before(() => app.start());
  after(() => {
    if (app && app.isRunning()) {
      return app.stop();
    }
  });

  describe('Unable to create project because of duplicated groups', () => {
    const xml = `
    <content>
      <radio group="animal" value="cat" vizual="Cat" />
      <checkbox group="animal" value="dog" vizual="Dog" />
    </content>
    `;
    createProjectWithTemplate(app, appPath, xml);

    it('display duplicate group error', async () => {
      await app.client.waitForVisible('//*[@id="create_project_error"]');
      const el = app.client.element('//*[@id="create_project_error"]');
      const err = await el.getText();
      expect(err).to.be.eq('Template has duplicate groups: animal');
    });
  });
});

describe('Application launch', function () {
  this.timeout(30000);
  before(() => app.start());
  after(() => {
    if (app && app.isRunning()) {
      return app.stop();
    }
  });

  describe('Unable to create project because of missing the attribute', () => {
    const xml = `
    <content>
      <radio group="animal" value="cat" vizual="Cat" />
      <checkbox group="animal" vizual="Dog" />
    </content>
    `;
    createProjectWithTemplate(app, appPath, xml);

    it('display duplicate group error', async () => {
      await app.client.waitForVisible('//*[@id="create_project_error"]');
      const el = app.client.element('//*[@id="create_project_error"]');
      const err = await el.getText();
      expect(err).to.be.eq('Element [checkbox] missing the attribute [value]');
    });
  });
});

/* eslint-disable no-unused-expressions */
import path from 'path';
import { expect } from 'chai';
import api from '../api';

export const makeUrl = (imgDir, filename) => path.normalize(`file://${path.join(imgDir, filename)}`);
export const getBtn = (app, device, i) => app.client.element(`//*[@id="${device}"]/ul/li[${i}]/div/button`);
export const getChbox = (app, device, i) => {
  const el = app.client.element(`//*[@id="${device}"]/div/ul/li[${i}]/label/input`);
  return el;
};

export const getPath = (app, el, pth) => app.client.elementIdElement(el.ELEMENT, pth);

export const assertRadioLabels = async (app, device, expectedLabels) => {
  const elements = await app.client.elements(`//*[@id="${device}"]/ul/li/div/button`);
  const actualLabels = await Promise.all(elements.value.map(async (el) => {
    const txt = await app.client.elementIdText(el.ELEMENT);
    return txt.value;
  }));
  expect(actualLabels).to.be.deep.eq(expectedLabels);
};

export function sleep(ms) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

export const createProject = async (appPath, xml) => {
  const imgDir = path.join(appPath, 'src', 'e2e', 'test_data', 'proj0');
  const glob0 = path.join(imgDir, '*.jpg');
  const glob1 = path.join(imgDir, '*.png');

  await api.post('/project', {
    description: {
      name: 'testproj0',
    },
    template: {
      task: 'classification',
      xml,
    },
    data_sources: [
      { type: 'local_directory', source_uri: glob0 },
      { type: 'local_directory', source_uri: glob1 },
    ],
  });
};

export const itNavigatesToProject = (app, appPath, xml) => {
  it('navigates to project page', async () => {
    await sleep(2000);
    await createProject(appPath, xml);
    await app.client.refresh();
    await app.client.waitUntilTextExists('span', 'testproj0');
    await sleep(2000);
    await app.client.element("button/*[@innertext='testproj0']").click();
  });
};

export const getSamplePath = (app, filename) => app.client.element(`small*=${filename}`);
export const getSampleClass = (app, filename) => getSamplePath(app, filename).element('../..').element('./span');

export function createProjectWithTemplate(app, appPath, xml) {
  it('opens Create New Project popup', async () => {
    await app.client.waitForVisible('button');
    await app.client.element('button').click();
    await app.client.waitForVisible('input');
  });

  it('inputs new project name', async () => {
    await sleep(2000);
    await app.client.element('input').setValue('testproj0');
  });

  it('show no project create error', async () => {
    const errorExists = await app.client.isExisting('//*[@id="create_project_error"]');
    expect(errorExists).to.be.false;
  });

  it('set project task', async () => {
    const template = "input[placeholder='Select task']";
    await app.client.waitForExist(template);
    await app.client.element(template).setValue('classification');
  });

  it('set xml for template', async () => {
    await app.client.elements('textarea').setValue(xml);
  });

  const srcInput = "input[placeholder='/some/path or /some/glob/*.jpg']";
  const imgDir = path.join(appPath, 'src', 'e2e', 'test_data', 'proj0');
  const glob0 = path.join(imgDir, '*.jpg');
  const glob1 = path.join(imgDir, '*.png');

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

  return [imgDir, glob0, glob1];
}

/* eslint-disable prefer-template */
/* eslint-disable no-unused-expressions */
import path from 'path';
import { expect } from 'chai';
import api from '../api';

export const makeUrl = (imgDir, filename) => path.normalize(`file://${path.join(imgDir, filename)}`);
export const getRadio = (app, device, i) => app.client.element(`//*[@id="${device}"]/div/label/ul/li[${i}]/div/label`);
export const getChbox = (app, device, i) => {
  const el = app.client.element(`//*[@id="${device}"]/div/label/ul/li[${i}]/label/input`);
  return el;
};

export const getPath = (app, el, pth) => app.client.elementIdElement(el.ELEMENT, pth);

export const assertRadioLabels = async (app, device, expectedLabels) => {
  const radios = `//*[@id="${device}"]/div/label/ul/li/div/label`;
  await app.client.waitForExist(radios);
  const elements = await app.client.elements(radios);
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

  const res = await api.post('/project', {
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

  expect(res).not.to.have.property('status');
};

export const clickButton = async (app, tag, text) => {
  await app.client.waitForText(tag, text);
  await sleep(1500);
  const el = await app.client.element(`${tag}=${text}`).element('../../button');
  await app.client.elementIdClick(el.value.ELEMENT);
};

export const clickLink = async (app, tag, text) => {
  await app.client.waitForText(tag, text);
  await app.client.element(`${tag}*=${text}`).element('..').click();
};

export const itNavigatesToProject = (app, appPath, xml) => {
  it('navigates to project page', async () => {
    await createProject(appPath, xml);
    try {
      await app.client.refresh();
    } catch (error) {
      const errorIsNavigatedError = error.message.includes('Inspected target navigated or closed');

      if (!errorIsNavigatedError) {
        throw error;
      }
    }
    await clickButton(app, 'span', 'testproj0');
    await app.client.waitForText('span', 'Begin assess');
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
    const template = "input[placeholder='Select template']";
    await app.client.waitForExist(template);
    await app.client.element(template).setValue('Classification');
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

export const getRadioState = async (app, device) => {
  const elements = await app.client.elements(`//*[@id="${device}"]/div/label/ul/li/div/label/input`);
  const disabled = await Promise.all(elements.value.map(async (el) => {
    const isSelected = await app.client.elementIdSelected(el.ELEMENT);
    return isSelected.value;
  }));

  return disabled;
};
export const getChecked = async (app, device) => {
  const elements = await app.client.elements(`//*[@id="${device}"]/div/label/ul/li/label/input`);
  const checked = await Promise.all(elements.value.map(async (el) => {
    const ch = await app.client.elementIdSelected(el.ELEMENT, 'checked');
    return ch.value;
  }));

  return checked;
};

const sortObj = (obj) => {
  // eslint-disable-next-line no-nested-ternary
  return obj === null || typeof obj !== 'object'
    ? obj
    : Array.isArray(obj)
      ? obj.map(sortObj)
      : Object.assign({},
        ...Object.entries(obj)
          .sort(([keyA], [keyB]) => keyA.localeCompare(keyB))
          .map(([k, v]) => ({ [k]: sortObj(v) })));
}

export const deterministicStrigify = (obj) => {
  return JSON.stringify(sortObj(obj));
};

export function expectSampleMarkupToBeEq(app, appPath, markup) {
  it('displays sample markup on project page', async () => {
    await clickLink(app, 'span', 'testproj0');
    await app.client.waitForExist('ul');

    const imgDir = path.join(appPath, 'src', 'e2e', 'test_data', 'proj0');

    const pathText = await getSamplePath(app, 'kek0.jpg').getText();
    const cl = await getSampleClass(app, 'kek0.jpg').getText();
    expect(pathText).to.be.eq(path.join(imgDir, 'kek0.jpg') + ':');
    expect(cl).to.be.eq(deterministicStrigify(markup));
  });
}
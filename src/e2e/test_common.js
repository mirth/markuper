import path from 'path';
import { expect } from 'chai';
import api from '../api';

export const makeUrl = (imgDir, filename) => path.normalize(`file://${path.join(imgDir, filename)}`);
export const getBtn = (app, i) => app.client.element('//*[@id="grid"]').element(`./div/div[2]/div/div/div/div[1]/ul/li[${i}]/div/button`);
export const getChbox = (app, device, i) => {
  const el = app.client.element(`//*[@id="${device}"]/div/ul/li[${i}]/label/input`);
  return el;
};

export const assertButtonLabels = async (app) => {
  let btnTxt = await getBtn(app, 1).element('.//span').getText();
  expect(btnTxt).to.be.eq('Cat');
  btnTxt = await getBtn(app, 2).element('.//span').getText();
  expect(btnTxt).to.be.eq('Dog');
  btnTxt = await getBtn(app, 3).element('.//span').getText();
  expect(btnTxt).to.be.eq('Chuk');
  btnTxt = await getBtn(app, 4).element('.//span').getText();
  expect(btnTxt).to.be.eq('Gek');
};

function sleep(ms) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

export const itNavigatesToProject = (app, appPath, xml) => {
  it('navigates to project page', async () => {
    await createProject(appPath, xml);
    await app.client.refresh();
    await app.client.waitUntilTextExists('span', 'testproj0');
    await sleep(1500);
    await app.client.element("button/*[@innertext='testproj0']").click();
  });
}

export const createProject = async(appPath, xml) => {
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
}
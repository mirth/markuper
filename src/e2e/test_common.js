import path from 'path';
import { expect } from 'chai';

export const makeUrl = (imgDir, filename) => path.normalize(`file://${path.join(imgDir, filename)}`);
export const getBtn = (app, i) => app.client.element('//*[@id="grid"]').element(`./div/div[2]/div/div/div/div[1]/ul/li[${i}]/div/button`);
export const getChbox = (app, device, i) => {
  const el = app.client.element(`//*[@id="${device}"]/div/ul/li[${i}]/label/input`);
  return el;
}

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

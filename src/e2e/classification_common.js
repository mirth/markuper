/* eslint-disable no-restricted-syntax */
import { expect } from 'chai';
import _ from 'lodash';
import {
  getPath, radioClick, assertRadioLabels, sleep,
  clickButton, getRadioState, getChecked, checkboxClick,
} from './test_common';

function expectFocusIsOn(app, devices, device) {
  it(`focused on device ${device}`, async () => {
    const cls = await Promise.all(devices.map(async (dev) => {
      const cl = await app.client.element(`//*[@id="${dev}"]/div`).getAttribute('class');
      return cl;
    }));

    const pairs = _.zip(devices, cls);

    for (const [iterName, iterCl] of pairs) {
      if (iterName === device) {
        expect(iterCl).to.have.string('selected');
      } else {
        expect(iterCl).to.not.have.string('selected');
      }
    }

    const devSubmitCl = await app.client.element('//*[@id="device_submit"]/div/button').getAttribute('class');

    if (device === 'devSubmit') {
      expect(devSubmitCl).to.have.string('filled');
    }
  });
}

async function assertCheckboxLabels(app, device) {
  const inputs = `//*[@id="${device}"]/div/label/ul/li/label/input`;
  await app.client.waitForExist(inputs);
  const elements = await app.client.elements(inputs);
  const labels = await Promise.all(elements.value.map(async (el) => {
    const pth = await getPath(app, el, '..');
    const txt = await app.client.elementIdText(pth.value.ELEMENT);
    return txt.value;
  }));
  expect(labels).to.be.deep.eq(['Black', 'White', 'Pink']);
}

export function itSubmitsSample(app) {
  it('submits the sample', async () => {
    await app.client.keys('Enter');
    await sleep(1500);
    await clickButton(app, 'span', 'testproj0');
    await sleep(1500);
  });
}

export function TestCheckboxRadioRadio(app, prefix) {
  const CN = (name) => `${prefix}/${name}`;

  const focusIsOn = (device) => {
    expectFocusIsOn(app, [CN('0'), CN('1'), CN('2')], device);
  };

  it('displays correct devices labels', async () => {
    await assertCheckboxLabels(app, CN('0'));
    await assertRadioLabels(app, CN('1'), ['Cat', 'Dog', 'Chuk', 'Gek']);
    await assertRadioLabels(app, CN('2'), ['Smoll', 'Big']);
  });

  focusIsOn(CN('0'));

  it('displays checkbox 1, 3 as checked', async () => {
    await checkboxClick(app, CN('0'), 1);
    await app.client.keys('2');
    await app.client.keys('1');
    await checkboxClick(app, CN('0'), 1);
    await app.client.keys('2');

    const checked = await getChecked(app, CN('0'));
    expect(checked).to.be.deep.eq([true, false, false]);
    await app.client.keys('Enter');
  });

  focusIsOn(CN('1'));

  it('displays radio 4 selected', async () => {
    await radioClick(app, CN('1'), 1);
    await radioClick(app, CN('1'), 1);
    await radioClick(app, CN('1'), 4);

    const selected = await getRadioState(app, CN('1'));
    expect(selected).to.be.deep.eq([false, false, false, true]);
    await app.client.keys('Enter');
  });

  focusIsOn(CN('2'));

  it('displays radio 2 selected', async () => {
    await radioClick(app, CN('2'), 2);
    const selected = await getRadioState(app, CN('2'));
    expect(selected).to.be.deep.eq([false, true]);
    await app.client.keys('Enter');
  });

  focusIsOn('device_submit');
}

export function TestRadioCheckbox(app, prefix) {
  const CN = (name) => `${prefix}/${name}`;

  const focusIsOn = (device) => {
    expectFocusIsOn(app, [CN('0'), CN('1')], device);
  };

  it('displays correct devices labels', async () => {
    await sleep(2000);
    await assertRadioLabels(app, CN('0'), ['Cat', 'Dog', 'Chuk', 'Gek']);
    await assertCheckboxLabels(app, CN('1'));
  });

  focusIsOn(CN('0'));

  it('tries to submit empty radio field', async () => {
    await app.client.keys('Enter');
    await sleep(1500);
  });

  focusIsOn(CN('0'));

  it('displays radio 3 selected', async () => {
    await app.client.keys('3');
    const selected = await getRadioState(app, CN('0'));
    expect(selected).to.be.deep.eq([false, false, true, false]);
  });

  focusIsOn(CN('0'));

  it('submits radio field', async () => {
    await app.client.keys('Enter');
  });

  focusIsOn(CN('1'));

  it('displays radio 1 selected', async () => {
    await radioClick(app, CN('0'), 1);
    const selected = await getRadioState(app, CN('0'));
    expect(selected).to.be.deep.eq([true, false, false, false]);
  });

  focusIsOn(CN('1'));

  it('displays checkbox 2 as checked', async () => {
    await app.client.keys('2');

    const checked = await getChecked(app, CN('1'));
    expect(checked).to.be.deep.eq([false, true, false]);
  });

  focusIsOn(CN('1'));

  it('displays checkboxes 1,2 as checked', async () => {
    await app.client.keys('1');

    const checked = await getChecked(app, CN('1'));
    expect(checked).to.be.deep.eq([true, true, false]);
  });

  focusIsOn(CN('1'));

  it('displays all checkboxes as checked', async () => {
    await app.client.keys('3');

    const checked = await getChecked(app, CN('1'));
    expect(checked).to.be.deep.eq([true, true, true]);
  });

  focusIsOn(CN('1'));

  it('displays checkbox 2 as unchecked', async () => {
    await app.client.keys('2');

    const checked = await getChecked(app, CN('1'));
    expect(checked).to.be.deep.eq([true, false, true]);
    await app.client.keys('Enter');
  });

  focusIsOn('device_submit');
}

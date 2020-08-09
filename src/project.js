/* eslint-disable no-nested-ternary */
/* eslint-disable no-restricted-syntax */

import _ from 'lodash';
import { push } from 'svelte-spa-router';
import { fetchProject, activeSample } from './store';
import api from './api';

export function goToProject(projectId) {
  return async () => {
    await fetchProject(projectId);
    push(`/project/${projectId}`);
  };
}

export function getFieldsInOrderFor(template) {
  const groupsOrder = template.fields_order;

  const byGroup = _.groupBy([
    ...(template.radios || []),
    ...(template.checkboxes || []),
    ...(template.bounding_boxes || []),
  ], 'group');
  const fields = groupsOrder.flatMap(((group) => byGroup[group]));

  const owner = template.group;
  for (const field of fields) {
    field.owner = owner;
  }

  return fields;
}

export function getProjectIDFromSampleID(sampleID) {
  const [projID] = sampleID.split('-');
  return projID;
}

// FIXME test_common.js
const sortObj = (obj) => (obj === null || typeof obj !== 'object'
  ? obj
  : Array.isArray(obj)
    ? obj.map(sortObj)
    : Object.assign({},
      ...Object.entries(obj)
        .sort(([keyA], [keyB]) => keyA.localeCompare(keyB))
        .map(([k, v]) => ({ [k]: sortObj(v) }))));


export const deterministicStrigify = (obj) => JSON.stringify(sortObj(obj));

function iterateXMLNodes(root, fn) {
  for (const node of root.childNodes) {
    fn(node);
    iterateXMLNodes(node, fn);
  }
}

function getRandomColor() {
  const letters = '0123456789ABCDEF';
  let color = '#';
  for (let i = 0; i < 6; i += 1) {
    color += letters[Math.floor(Math.random() * 16)];
  }
  return color;
}


export function enrichWithColor(templateXML) {
  const parser = new DOMParser();
  const xmlDoc = parser.parseFromString(templateXML, 'text/xml');
  const root = xmlDoc.getElementsByTagName('content')[0];

  iterateXMLNodes(root, (node) => {
    if (node.nodeName === '#text') {
      return;
    }

    const color = getRandomColor();
    node.setAttribute('color', color);
  });

  const newXml = `<content>${xmlDoc.documentElement.innerHTML}</content>`;

  return newXml;
}

/* eslint-disable no-restricted-syntax */

import _ from 'lodash';
import { push } from 'svelte-spa-router';
import { fetchProject } from './store';

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

// FIXME test_common.js
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
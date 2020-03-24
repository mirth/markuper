import _ from 'lodash';

export function getFieldsInOrderFor(template) {
  const groupsOrder = template.fields_order;
  const byGroup = _.groupBy([
    ...template.radios,
    ...template.checkboxes,
  ], 'group');
  const fields = groupsOrder.flatMap(((group) => byGroup[group]));

  return [fields, groupsOrder];
}

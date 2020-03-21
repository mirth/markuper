import _ from 'lodash';

export function makeLabelsWithKeys(fieldLabels) {
  const allKeys = _.range(1, 10).map(String);
  const labelsNumber = Math.min(fieldLabels.length, allKeys.length);
  const keys = allKeys.slice(0, labelsNumber);
  const labelsWithKeys = _.zip(fieldLabels, keys).concat(fieldLabels.slice(labelsNumber));

  return [keys, labelsWithKeys];
}

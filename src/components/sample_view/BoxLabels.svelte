<script>
/* eslint-disable no-restricted-syntax */
/* eslint-disable no-unused-vars */
import _ from 'lodash';
import { getFieldsInOrderFor } from '../../project';

export let markup;
export let template;

const fields = getFieldsInOrderFor(template.bounding_boxes[0]);


const byGroupByLabel = {};
for (const field of fields) {
  byGroupByLabel[field.group] = _.keyBy(field.labels, 'value');
}

const notBoxes = Object.entries(markup).filter(([k, _v]) => k !== 'box');
const labels = notBoxes.flatMap(([group, values]) => {
  if (_.isArray(values)) {
    return values.map((value) => [group, byGroupByLabel[group][value]]);
  }

  return [[group, byGroupByLabel[group][values]]];
});

</script>



<ul>
  {#each labels as [group, value]}
    <li>
      <span style={`background-color: ${value.display_color}`}>
        {value.display_name}
      </span>
    </li>
  {/each}
</ul>

<style>
li {
  list-style-type: none;
  display: inline;
}
ul {
  margin: 0;
  padding: 0;
}
</style>
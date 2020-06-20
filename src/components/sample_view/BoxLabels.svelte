<script>
import _ from 'lodash';
import { activeProject } from '../../store';
import { getFieldsInOrderFor } from '../../project';

export let markup;
export let template;

let fields = getFieldsInOrderFor(template.bounding_boxes[0]);


let byGroupByLabel = {}
for(const field of fields) {
  byGroupByLabel[field.group] = _.keyBy(field.labels, 'value');
}

const notBoxes = Object.entries(markup).filter(([k, v]) => k !== 'box');
const labels = notBoxes.flatMap(([group, values]) => {
  if(_.isArray(values)) {
    return values.map(value => [group, byGroupByLabel[group][value]]);
  }

  return [[group, byGroupByLabel[group][values]]];
})

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
<script>

import ControlList from './ControlList.svelte';
import { assessState, activeMarkup } from '../../store';
import { getFieldsInOrderFor } from '../../project';

export let submitMarkupAndFetchNext;
export let sample;;

const rootPrefix = 'root';
let fieldIter = 0;
$: $activeMarkup = (sample.markup && sample.markup.markup) || {};

// const [fields, groupsOrder] = getFieldsInOrderFor(sample.project.template);
// groupsOrder.push(submitGroup);

// console.log('KEKE1: ', fields)
// console.log('KEKE2: ', groupsOrder)


$assessState.controlPrefix = rootPrefix;

</script>


<ControlList
  {submitMarkupAndFetchNext}
  templateLike={sample.project.template}
  prefix={rootPrefix}
  />




<script>

import ControlList from './ControlList.svelte';
import { assessState, activeMarkup } from '../../store';
import { getFieldsInOrderFor } from '../../project';

export let submitMarkupAndFetchNext;
export let sample;

let fieldIter = 0;
$: $activeMarkup = (sample.markup && sample.markup.markup) || {};


$assessState.focusedGroup = sample.project.template.fields_order[0];
$assessState.curFocusedOwner = 'root';
$assessState.lastFocusedOwner = 'root';
</script>


<ControlList
  {submitMarkupAndFetchNext}
  owner={sample.project.template}
  />

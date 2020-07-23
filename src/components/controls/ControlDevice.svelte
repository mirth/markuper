<script>

import ControlList from './ControlList.svelte';
import { assessState, sampleMarkup } from '../../store';
import { getFieldsInOrderFor } from '../../project';
import { submitGroup } from '../../control';

export let submitMarkupAndFetchNext;
export let sample;

function ownerByChild(childGroup) {
  const rootFields = getFieldsInOrderFor(sample.project.template);
  {
    const idx = sample.project.template.fields_order.indexOf(childGroup);
    if (idx !== -1) {
      const field = rootFields[idx];
      return [sample.project.template, field];
    }
  }


  const owner = rootFields.filter((o) => o.type === 'bounding_box').find((o) => {
    const idx = o.fields_order.indexOf(childGroup);
    return idx !== -1;
  });

  const field = getFieldsInOrderFor(owner)[owner.fields_order.indexOf(childGroup)];

  return [owner, field];
}

function tryIncrementGroup(callerGroup) {
  if ($assessState.focusedGroup === submitGroup) {
    return;
  }

  if (callerGroup !== $assessState.focusedGroup) {
    return;
  }

  const [owner] = ownerByChild($assessState.focusedGroup);
  const ownerFieldsOrder = owner.fields_order;
  const idx = ownerFieldsOrder.indexOf($assessState.focusedGroup);

  const isLastControl = idx === (ownerFieldsOrder.length - 1);
  if (isLastControl) {
    const isNested = owner.group;
    if (isNested) {
      $assessState.focusedGroup = owner.group;
    } else {
      $assessState.focusedGroup = submitGroup;
    }
  } else {
    $assessState.focusedGroup = ownerFieldsOrder[idx + 1];
  }
}


/* eslint-disable prefer-const */
$sampleMarkup = (sample.markup && sample.markup.markup) || {};
$assessState.markup = {};

[$assessState.focusedGroup] = sample.project.template.fields_order;

</script>

<ControlList
  {submitMarkupAndFetchNext}
  onFieldCompleted={tryIncrementGroup}
  owner={sample.project.template}
  />

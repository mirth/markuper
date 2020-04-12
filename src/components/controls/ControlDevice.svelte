<script>

import ControlList from './ControlList.svelte';
import { assessState, sampleMarkup } from '../../store';
import { getFieldsInOrderFor } from '../../project';
import { submitGroup } from '../../control';

export let submitMarkupAndFetchNext;
export let sample;

let keyDown;

$sampleMarkup = (sample.markup && sample.markup.markup) || {};
$assessState.markup = {};

function handleKeydown(event) {
  keyDown = event.key;
}

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

function isFieldFilled(field) {
  if (field.type === 'radio' && !Object.hasOwnProperty.call($assessState.markup, field.group)) {
    return false;
  }

  return true;
}

function pushActiveMarkup() {
  const [owner, field] = ownerByChild($assessState.focusedGroup);

  const $assessStateMarkup = JSON.parse(JSON.stringify($assessState.markup));
  if (owner.group) {
    $sampleMarkup[owner.group] = [
      ...$sampleMarkup[owner.group],
      $assessStateMarkup,
    ];
  } else if (field.type !== 'bounding_box') {
    $sampleMarkup = $assessStateMarkup;
  }
}

function tryIncrementIter() {
  const [owner, field] = ownerByChild($assessState.focusedGroup);
  const ownerFieldsOrder = owner.fields_order;
  const idx = ownerFieldsOrder.indexOf($assessState.focusedGroup);

  if (idx === (ownerFieldsOrder.length - 1)) {
    if (isFieldFilled(field)) {
      pushActiveMarkup();
      if (!owner.group) {
        $assessState.focusedGroup = submitGroup;
      } else {
        $assessState.focusedGroup = owner.group;
        $assessState.markup = {};
      }
    }
  } else if (isFieldFilled(field)) {
    $assessState.focusedGroup = ownerFieldsOrder[idx + 1];
  }
}

function handleKeyup(event) {
  if (event.key !== keyDown) {
    keyDown = null;
    return;
  }

  if (event.key !== 'Enter') {
    return;
  }

  if ($assessState.focusedGroup === submitGroup) {
    return;
  }

  tryIncrementIter();
}

[$assessState.focusedGroup] = sample.project.template.fields_order;

</script>

<svelte:window on:keydown={handleKeydown} on:keyup={handleKeyup}/>

<ControlList
  {submitMarkupAndFetchNext}
  owner={sample.project.template}
  />

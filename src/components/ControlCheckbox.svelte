<script>
import Checkbox from 'svelte-atoms/Checkbox.svelte';
import Spacer from 'svelte-atoms/Spacer.svelte';
import { makeLabelsWithKeys } from '../control';
import { activeMarkup, assessState } from '../store';
import KeyboardButton from './KeyboardButton.svelte';

export let field;

const [keys, labelsWithKeys] = makeLabelsWithKeys(field.labels);

const checkedByLabel = {};
let keyDown;

function handleKeydown(event) {
  if (field.group !== $assessState.focusedGroup) {
    return;
  }

  keyDown = event.key;
}

function updateMarkupWith(labelValue) {
  checkedByLabel[labelValue] = !checkedByLabel[labelValue];
  $activeMarkup[field.group] = Object.entries(
    checkedByLabel,
  )
    .filter((kv) => kv[1])
    .map((kv) => kv[0]);
}

async function handleKeyup(event) {
  if (field.group !== $assessState.focusedGroup) {
    return;
  }

  if (event.key !== keyDown) {
    return;
  }

  if (!keys.includes(event.key)) {
    return;
  }

  const labelIndex = parseInt(event.key, 10) - 1;
  const label = field.labels[labelIndex];

  updateMarkupWith(label.value);

  keyDown = null;
}

function onChangeFor(label) {
  return (ev) => {
    if (ev.target.value === label.value) {
      updateMarkupWith(label.value);
    }
  };
}

</script>

<svelte:window on:keydown={handleKeydown} on:keyup={handleKeyup}/>
<div>
<ul>
{#each labelsWithKeys as [label, key]}
  <li>
    <Checkbox
      checked={checkedByLabel[label.value]}
      on:change={onChangeFor(label.value)}
      >
      {label.vizual}
    </Checkbox>
    <Spacer size={8} />
    <KeyboardButton {key} isKeyDown={key === keyDown} />
  </li>
{/each}
</ul>
</div>

<style>
ul {
  display: inline;
  padding: 0;
}

li {
  list-style-type: none;
  display: inline-block;
  margin-right: 20px;
}

div {
  padding-top: 25px;
}
</style>
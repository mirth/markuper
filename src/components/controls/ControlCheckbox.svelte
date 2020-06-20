<script>
import Checkbox from 'svelte-atoms/Checkbox.svelte';
import Spacer from 'svelte-atoms/Spacer.svelte';
import { makeLabelsWithKeys } from '../../control';
import { sampleMarkup, assessState, isFieldSelected } from '../../store';
import KeyboardButton from './KeyboardButton.svelte';

export let field;

const [keys, labelsWithKeys] = makeLabelsWithKeys(field.labels);

let keyDown;
let checked = new Set([]);
let isSelected = false;
if (!field.owner && $sampleMarkup[field.group]) {
  $assessState.markup[field.group] = $sampleMarkup[field.group];
}

$: checked = new Set($assessState.markup[field.group] || []);
$: isSelected = isFieldSelected(field, $assessState);

function handleKeydown(event) {
  if (!isSelected) {
    return;
  }

  keyDown = event.key;
}

function updateMarkupWith(labelValue) {
  if (checked.has(labelValue)) {
    checked.delete(labelValue);
  } else {
    checked.add(labelValue);
  }

  $assessState.markup[field.group] = Array.from(checked);
}

async function handleKeyup(event) {
  if (!isSelected) {
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

function onChangeFor(labelValue) {
  return () => {
    updateMarkupWith(labelValue);
  };
}

</script>

<svelte:window on:keydown={handleKeydown} on:keyup={handleKeyup}/>
<ul>
{#each labelsWithKeys as [label, key]}
  <li>
    <Checkbox
      checked={checked.has(label.value)}
      on:change={onChangeFor(label.value)}
      >
      <span style={`color: ${label.display_color}`}>{label.display_name}</span>
    </Checkbox>
    <Spacer size={8} />
    {#if isSelected}
      <KeyboardButton {key} isKeyDown={key === keyDown} />
    {/if}
  </li>
{/each}
</ul>

<style>
ul {
  display: inline;
  padding: 0;
}

li {
  list-style-type: none;
  display: inline-block;
  margin-right: 8px;
  padding: 0px;
}

</style>
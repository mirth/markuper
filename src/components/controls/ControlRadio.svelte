<script>
import Cell from 'svelte-atoms/Grids/Cell.svelte';
import Spacer from 'svelte-atoms/Spacer.svelte';
import Radio from 'svelte-atoms/Radio.svelte';
import { makeLabelsWithKeys } from '../../control';
import { sampleMarkup, assessState, isFieldSelected } from '../../store';
import KeyboardButton from './KeyboardButton.svelte';

export let field;

const [keys, labelsWithKeys] = makeLabelsWithKeys(field.labels);

let keyDown;
let isSelected = false;

if(!field.owner && $sampleMarkup[field.group]) {
  $assessState.markup[field.group] = $sampleMarkup[field.group];
}

$: radio = $assessState.markup[field.group];

$: isSelected = isFieldSelected(field, $assessState);



function handleKeydown(event) {
  if (!isSelected) {
    return;
  }

  keyDown = event.key;
}

function updateRadio(labelValue) {
  radio = labelValue;
  $assessState.markup[field.group] = radio;
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

  updateRadio(label.value);

  keyDown = null;
}

</script>

<svelte:window on:keydown={handleKeydown} on:keyup={handleKeyup} />

<ul>
{#each labelsWithKeys as [label, key]}
  <li>
    <Cell>
      <Radio
        bind:group={radio}
        value={label.value}
      >
        {label.vizual}
      </Radio>
      <Spacer size={8} />
      {#if isSelected}
        <KeyboardButton {key} isKeyDown={key === keyDown} />
      {/if}
    </Cell>
  </li>
{/each}
</ul>

<style>

ul {
  padding: 0;
  margin: 0;
}

li {
  list-style-type: none;
  display: inline-block;
  text-align: center;
}

</style>
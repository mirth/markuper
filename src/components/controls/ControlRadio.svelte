<script>
import Cell from 'svelte-atoms/Grids/Cell.svelte';
import Spacer from 'svelte-atoms/Spacer.svelte';
import Radio from 'svelte-atoms/Radio.svelte';
import { makeLabelsWithKeys } from '../../control';
import { activeMarkup, assessState, isFieldSelected } from '../../store';
import KeyboardButton from './KeyboardButton.svelte';

export let field;

const [keys, labelsWithKeys] = makeLabelsWithKeys(field.labels);

let keyDown;
let isSelected = false;
$: isSelected = isFieldSelected(field, $assessState);
let radio = $activeMarkup[field.group];
$: if (radio) {
  $activeMarkup[field.group] = radio;
}


function handleKeydown(event) {
  if (!isSelected) {
    return;
  }

  keyDown = event.key;
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
  radio = label.value;
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
}

li {
  list-style-type: none;
  display: inline-block;
  text-align: center;
}

</style>
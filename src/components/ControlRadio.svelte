<script>
import Cell from 'svelte-atoms/Grids/Cell.svelte';
import Spacer from 'svelte-atoms/Spacer.svelte';
import Button from 'svelte-atoms/Button.svelte';
import { makeLabelsWithKeys } from '../control';
import { activeMarkup, assessState } from '../store';
import KeyboardButton from './KeyboardButton.svelte';

export let field;

const [keys, labelsWithKeys] = makeLabelsWithKeys(field.labels);

let keyDown;

function handleKeydown(event) {
  if (field.group !== $assessState.focusedGroup) {
    return;
  }

  keyDown = event.key;
}

function handleButtonClick(label) {
  return async () => {
    $activeMarkup[field.group] = label.value;
  };
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
  await handleButtonClick(label)();
  keyDown = null;
}

</script>

<svelte:window on:keydown={handleKeydown} on:keyup={handleKeyup} />

<ul>
{#each labelsWithKeys as [label, key]}
  <li>
    <Cell>
      <Button
        on:click={handleButtonClick(label)}
        disabled={$activeMarkup[field.group] === label.value}
        style='display: inline; min-width: 60px;'
      >
        {label.vizual}
      </Button>
      <Spacer size={8} />
      <KeyboardButton {key} isKeyDown={key === keyDown} />
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
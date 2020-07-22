<script>
import Cell from 'svelte-atoms/Grids/Cell.svelte';
import Spacer from 'svelte-atoms/Spacer.svelte';
import Radio from 'svelte-atoms/Radio.svelte';
import { makeLabelsWithKeys } from '../../control';
import { sampleMarkup, assessState, isFieldSelected } from '../../store';
import KeyboardButton from './KeyboardButton.svelte';
import WithEnterForGroup from './WithEnterForGroup.svelte';


export let field;
export let onFieldCompleted;

function handleKeydown(event) {
  if (!isSelected) {
    return;
  }

  keyDown = event.key;
}

function updateRadio(labelValue) {
  $sampleMarkup[field.group] = labelValue
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

function onEnterPressed() {
  if(!$sampleMarkup[field.group]) {
    return;
  }

  onFieldCompleted(field.group);
}

const [keys, labelsWithKeys] = makeLabelsWithKeys(field.labels);

let keyDown;
let isSelected = false;

$: radio = $sampleMarkup[field.group];
$: isSelected = isFieldSelected(field, $assessState);

</script>

<svelte:window on:keydown={handleKeydown} on:keyup={handleKeyup} />

<WithEnterForGroup {field} {onEnterPressed}/>

<ul>
{#each labelsWithKeys as [label, key]}
  <li>
    <Cell>
      <Radio
        bind:group={radio}
        value={label.value}
        on:click={() => updateRadio(label.value)}
      >
        <span style={`color: ${label.display_color}`}>{label.display_name}</span>
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
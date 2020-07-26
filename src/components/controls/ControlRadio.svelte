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

const [, labelsWithKeys] = makeLabelsWithKeys(field.labels);

$: isSelected = isFieldSelected(field, $assessState);
$: radio = $sampleMarkup[field.group];

function updateRadio(labelValue) {
  $sampleMarkup[field.group] = labelValue;
}

async function handleKeyboardKeyPressed(key) {
  const labelIndex = parseInt(key, 10) - 1;
  const label = field.labels[labelIndex];

  updateRadio(label.value);
}

function onEnterPressed() {
  if (!$sampleMarkup[field.group]) {
    return;
  }

  onFieldCompleted(field.group);
}

</script>

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
        <KeyboardButton {field} {key} onKeyPressed={handleKeyboardKeyPressed} />
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
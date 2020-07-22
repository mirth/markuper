<script>
import _ from 'lodash';
import Checkbox from 'svelte-atoms/Checkbox.svelte';
import Spacer from 'svelte-atoms/Spacer.svelte';
import { makeLabelsWithKeys } from '../../control';
import { sampleMarkup, assessState, isFieldSelected } from '../../store';
import KeyboardButton from './KeyboardButton.svelte';
import WithEnterForGroup from './WithEnterForGroup.svelte';

export let onFieldCompleted;
export let field;

function handleKeydown(event) {
  if (!isSelected) {
    return;
  }

  keyDown = event.key;
}

function saveMarkup() {
  const values = checked.map((checked, i) => {
    return checked ? field.labels[i].value : null;
  });

  $sampleMarkup[field.group] = _.compact(values);
}

function updateCheckedWith(labelIndex) {
  checked[labelIndex] = !checked[labelIndex];
  saveMarkup();
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
  updateCheckedWith(labelIndex);

  keyDown = null;
}

function onEnterPressed() {
  onFieldCompleted(field.group);
}

function tryInitWithSample() {
  const checked = new Array(field.labels.length);

  const markuped = $sampleMarkup[field.group] || [];
  for(let i = 0; i < field.labels.length; i++) {
    if(markuped.indexOf(field.labels[i].value) !== -1) {
      checked[i] = true;
    }
  }

  return checked;
}

const [keys, labelsWithKeys] = makeLabelsWithKeys(field.labels);
let isSelected = false;
let keyDown;

const checked = tryInitWithSample();

$: isSelected = isFieldSelected(field, $assessState);

</script>

<WithEnterForGroup {field} {onEnterPressed}/>

<svelte:window on:keydown={handleKeydown} on:keyup={handleKeyup}/>


<ul>
{#each labelsWithKeys as [label, key], labelIndex}
  <li>
    <Checkbox bind:checked={checked[labelIndex]}>
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
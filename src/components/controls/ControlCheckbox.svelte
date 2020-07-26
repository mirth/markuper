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

const [, labelsWithKeys] = makeLabelsWithKeys(field.labels);

$: isSelected = isFieldSelected(field, $assessState);

function tryInitWithSample(smplMarkup) {
  const check = new Array(field.labels.length);

  const markuped = smplMarkup[field.group] || [];
  for (let i = 0; i < field.labels.length; i += 1) {
    if (markuped.indexOf(field.labels[i].value) !== -1) {
      check[i] = true;
    }
  }

  return check;
}

$: checked = tryInitWithSample($sampleMarkup);

function saveMarkup() {
  const values = checked.map((check, i) => (check ? field.labels[i].value : null));

  $sampleMarkup[field.group] = _.compact(values);
}

function updateCheckedWith(labelIndex) {
  checked[labelIndex] = !checked[labelIndex];
  saveMarkup();
}

function handleKeyboardKeyPressed(key) {
  const labelIndex = parseInt(key, 10) - 1;
  updateCheckedWith(labelIndex);
}

function onEnterPressed() {
  onFieldCompleted(field.group);
}

</script>

<WithEnterForGroup {field} {onEnterPressed}/>


<ul>
{#each labelsWithKeys as [label, key], labelIndex}
  <li>
    <Checkbox bind:checked={checked[labelIndex]}>
      <span style={`color: ${label.display_color}`}>{label.display_name}</span>
    </Checkbox>
    <Spacer size={8} />
    {#if isSelected}
      <KeyboardButton {field} {key} onKeyPressed={handleKeyboardKeyPressed} />
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
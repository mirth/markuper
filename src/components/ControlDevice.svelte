<script>
import _ from 'lodash';
import Cell from 'svelte-atoms/Grids/Cell.svelte';
import Spacer from 'svelte-atoms/Spacer.svelte';
import Button from 'svelte-atoms/Button.svelte';

export let makeHandleAssess;
export let sample;

const field = sample.template.radios[0];
const allKeys = _.range(1, 10).map(String);
const labelsNumber = Math.min(field.labels.length, allKeys.length);
const keys = allKeys.slice(0, labelsNumber);


const labelsWithKeys = _.zip(field.labels, keys).concat(field.labels.slice(labelsNumber));

let keyDown;

function handleKeydown(event) {
  keyDown = event.key;
}

async function handleKeyup(event) {
  if (event.key !== keyDown) {
    return;
  }

  if (!keys.includes(event.key)) {
    return;
  }

  const labelIndex = parseInt(event.key, 10) - 1;
  const label = field.labels[labelIndex];
  await makeHandleAssess(field, label)();
}

</script>

<svelte:window on:keydown={handleKeydown} on:keyup={handleKeyup}/>


<ul>
{#each labelsWithKeys as [label, key]}
  <li>
    <Cell>
      <Button
        on:click={makeHandleAssess(field, label)}
        disabled={sample.markup && sample.markup.markup.class === label.value}
        style='display: inline; min-width: 60px;'
      >
        {label.name}
      </Button>
      <Spacer size={8} />
      <kbd class:kbd-down='{key === keyDown}'>{key}</kbd>
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

kbd {
  background-color: #eee;

  border-radius: 4px;
  font-size: 1em;
  padding: 0.2em 0.5em;
  border-top: 5px solid rgba(255,255,255,0.5);
  border-left: 5px solid rgba(255,255,255,0.5);
  border-right: 5px solid rgba(0,0,0,0.2);
  border-bottom: 5px solid rgba(0,0,0,0.2);
  color: #555;
}

.kbd-down {
  background-color: rgb(197, 50, 50);
}

</style>
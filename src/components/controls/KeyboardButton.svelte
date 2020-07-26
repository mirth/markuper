<script>
import { assessState, isFieldSelected } from '../../store';
import { makeLabelsWithKeys } from '../../control';

export let field;
export let key;
export let onKeyPressed;


let keyDown;
const [keys] = field.labels ? makeLabelsWithKeys(field.labels) : ['Enter'];

$: isSelected = isFieldSelected(field, $assessState);

function handleKeydown(event) {
  if (!isSelected) {
    return;
  }

  if (key !== event.key) {
    return;
  }

  keyDown = event.key;
}

async function handleKeyup(event) {
  if (!isSelected) {
    return;
  }

  if (key !== event.key) {
    return;
  }

  if (event.key !== keyDown) {
    return;
  }

  if (!keys.includes(event.key)) {
    return;
  }

  onKeyPressed(event.key);

  keyDown = null;
}

</script>

<kbd class:kbd-down={key === keyDown}>{key}</kbd>

<svelte:window on:keydown={handleKeydown} on:keyup={handleKeyup}/>


<style>
kbd {
  display: inline-block;

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
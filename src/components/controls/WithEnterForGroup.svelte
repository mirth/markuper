<script>

import { isFieldSelected, assessState } from '../../store';

export let field;
export let onEnterPressed;

let keyDown;

function handleKeydown(event) {
  if(!isFieldSelected(field, $assessState)) {
    return;
  }

  keyDown = event.key;
}

function handleKeyup(event) {
  if(!isFieldSelected(field, $assessState)) {
    return;
  }

  if (event.key !== keyDown) {
    return;
  }


  if (event.key === 'Enter') {
    onEnterPressed();
  }
}

</script>

<svelte:window on:keydown={handleKeydown} on:keyup={handleKeyup}/>
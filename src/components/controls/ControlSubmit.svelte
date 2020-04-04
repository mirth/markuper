<script>
import Button from 'svelte-atoms/Button.svelte';
import KeyboardButton from './KeyboardButton.svelte';


export let submitMarkupAndFetchNext;
export let isSelected;

let keyDown = null;

function handleKeydown(event) {
  if (!isSelected) {
    return;
  }

  keyDown = event.key;
}

function handleKeyup(event) {
  if (!isSelected) {
    return;
  }

  if (event.key !== keyDown) {
    return;
  }

  if (event.key === 'Enter') {
    submitMarkupAndFetchNext();
  }
}

</script>

<svelte:window on:keydown={handleKeydown} on:keyup={handleKeyup}/>

<div>
  <Button
    on:click={submitMarkupAndFetchNext}
    type={isSelected ? 'filled' : 'flat'}
    id='device_submit'
    style='display: inline;'
    >
      Submit
  </Button>
  {#if isSelected}
    <KeyboardButton key={'Enter'} isKeyDown={keyDown === 'Enter'} />
  {/if}
</div>
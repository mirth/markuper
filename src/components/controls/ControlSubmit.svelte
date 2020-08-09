<script>
import Button from 'svelte-atoms/Button.svelte';
import KeyboardButton from './KeyboardButton.svelte';
import { submitGroup } from '../../control';

export let submitMarkupAndFetchNext;
export let isSelected;

</script>


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
    <KeyboardButton field={{ group: submitGroup }} key={'Enter'} bind:onKeyPressed={submitMarkupAndFetchNext} />
  {/if}
</div>
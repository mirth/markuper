<script>
import Input from 'svelte-atoms/Input.svelte';
import { dirty } from '../forms';

export let dataSource;
export let disabled;

let isNewSourceValid = false;
$: isNewSourceValid = dataSource.isValid();

const dirties = new Set([]);
const [makeDirty, isDirty] = dirty(dirties, 'dataSource');

</script>

<Input
  bind:value={dataSource.source_uri}
  placeholder='/some/path or /some/glob/*.jpg'
  size='small'
  {disabled}
  on:input={makeDirty}
  />
{#if !isNewSourceValid && isDirty()}
  <small><span>This doesn't look like valid path or glob</span></small>
{/if}

<style>
span {
  color: red;
}
</style>
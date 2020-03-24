<script>
import Input from 'svelte-atoms/Input.svelte';

export let dataSource;
export let disabled;

let isNewSourceValid = false;
$: isNewSourceValid = dataSource.isValid();

</script>

<Input
  bind:value={dataSource.source_uri}
  placeholder='/some/path or /some/glob/*.jpg'
  size='small'
  {disabled}
  />
{#if !isNewSourceValid}
  <small><span>This doesn't look like valid path or glob</span></small>
{/if}

<style>
span {
  color: red;
}
</style>
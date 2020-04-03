<script>
import { onMount } from 'svelte';
import { push } from 'svelte-spa-router';
import api from '../api';
import IconButton from 'svelte-atoms/IconButton.svelte';
import Spacer from 'svelte-atoms/Spacer.svelte';

let version;
onMount(async () => {
  version = await api.get('/version');
});
</script>

<div>
  <span><small>{version}</small></span>
  <Spacer size={16} />
  <IconButton
    icon='monitor'
    on:click={() => push(`/`)}
    type='filled' status='primary'
    style='display: inline;'
  />
</div>


<style>
div {
  text-align: center;
}
</style>
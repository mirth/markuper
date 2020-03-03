<script>
import Input from 'svelte-atoms/Input.svelte';

export let dataSource;

</script>

<Input bind:value={dataSource.source_uri} placeholder="/some/path" minlength="1" size="small"/>
<script>
import { onMount } from 'svelte';
import { fetchProjectList } from '../store';
import SideBar from './SideBar.svelte';

onMount(fetchProjectList);

</script>

<style>
#grid {
  display: grid;
  height: 100px;
  grid-template-columns: repeat(6, 1fr);
  grid-template-rows: 100px;
}
</style>

<div id='grid'>
  <SideBar />
  <slot></slot>
</div>
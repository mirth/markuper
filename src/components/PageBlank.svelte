<script>
import { onMount } from 'svelte';
import Container from "svelte-atoms/Grids/Container.svelte";
import Row from "svelte-atoms/Grids/Row.svelte";
import Cell from "svelte-atoms/Grids/Cell.svelte";
import { fetchProjectList } from '../store';
import SideBar from './SideBar.svelte';

onMount(fetchProjectList);

</script>


<Container>
<div id='grid'>
  <Row>
    <Cell xs={2}>
      <SideBar />
    </Cell>
    <Cell xs={10}>
      <Container>
        <slot></slot>
      </Container>
    </Cell>
  </Row>
</div>
</Container>

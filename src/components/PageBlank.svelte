<script>
import { onMount } from 'svelte';
import Container from 'svelte-atoms/Grids/Container.svelte';
import Row from 'svelte-atoms/Grids/Row.svelte';
import Cell from 'svelte-atoms/Grids/Cell.svelte';
import { fetchProjectList } from '../store';
import SideBar from './SideBar.svelte';
import api from '../api';

let version
onMount(async () => {
  fetchProjectList();
  version = await api.get('/version');
});

</script>


<Container>
<div id='grid'>
  <Row>
    <Cell xs={2}>
      <span>{version}</span>
      <SideBar />
    </Cell>
    <Cell xs={10}>
      <Container>
        <Row>
          <Cell>
            <slot></slot>
          </Cell>
        </Row>
      </Container>
    </Cell>
  </Row>
</div>
</Container>

<script>
import { onMount } from 'svelte';
import Container from 'svelte-atoms/Grids/Container.svelte';
import Row from 'svelte-atoms/Grids/Row.svelte';
import Cell from 'svelte-atoms/Grids/Cell.svelte';
import api from '../api';

let version;
onMount(async () => {
  version = await api.get('/version');
});

</script>


<Container>
<div id='grid'>
  <span>{version}</span>
  <Row>
    <Cell xs={12}>
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

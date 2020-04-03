<script>
import Container from 'svelte-atoms/Grids/Container.svelte';
import Row from 'svelte-atoms/Grids/Row.svelte';
import Cell from 'svelte-atoms/Grids/Cell.svelte';
import SideBar from './SideBar.svelte';



</script>


<Container>
<div id='grid'>
  <Row>
    <Cell xs={1}>
      <SideBar />
    </Cell>
    <Cell xs={11}>
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

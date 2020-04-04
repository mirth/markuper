<script>
import { push } from 'svelte-spa-router';
import Row from 'svelte-atoms/Grids/Row.svelte';
import Cell from 'svelte-atoms/Grids/Cell.svelte';
import Block from 'svelte-atoms/Block.svelte';
import Spacer from 'svelte-atoms/Spacer.svelte';
import Button from 'svelte-atoms/Button.svelte';
import Container from 'svelte-atoms/Grids/Container.svelte';
import ProjectList from './ProjectList.svelte';

</script>

<Container>
<Row>
<Cell xs={4}><Block></Block></Cell>
<Cell xs={4}><Block>
  <Spacer size={128} />
  <div>
    <Button
      on:click={() => push('/create_project')}
      iconRight='chevron-right'
      style='display: inline;'
    >
      Create new project
    </Button>
  </div>
</Block></Cell>
<Cell xs={4}><Block></Block></Cell>
</Row>
<ProjectList />
</Container>


<style>
div {
  text-align: center;
}
</style>
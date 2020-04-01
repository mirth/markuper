<script>
import _ from 'lodash';
import { onMount } from 'svelte';
import Button from 'svelte-atoms/Button.svelte';
import { projects, fetchProject } from '../store';
import Row from 'svelte-atoms/Grids/Row.svelte';
import Cell from 'svelte-atoms/Grids/Cell.svelte';
import { fetchProjectList } from '../store';
import Container from 'svelte-atoms/Grids/Container.svelte';
import { goToProject  } from '../project';


let grouped = [];
$: grouped = _.chunk($projects, 4);

onMount(fetchProjectList);

</script>

<Container>
<Row>
<Cell xsOffset={1} xs={10}>
  {#each grouped as chunk}
    <Row>
      {#each chunk as project}
      <Cell>
        <div>
          <Button
            type='empty'
            on:click={goToProject(project.project_id)}
            style='padding: 0;'
          >
            {project.description.name}
          </Button>
        </div>
      </Cell>
      {/each}
    </Row>
  {/each}
</Cell>
</Row>
</Container>

<style>
div {
  margin-top: 30px;
}

</style>
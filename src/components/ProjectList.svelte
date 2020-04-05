<script>
import _ from 'lodash';
import { onMount } from 'svelte';
import Block from "svelte-atoms/Block.svelte";
import Button from 'svelte-atoms/Button.svelte';
import Row from 'svelte-atoms/Grids/Row.svelte';
import Cell from 'svelte-atoms/Grids/Cell.svelte';
import Container from 'svelte-atoms/Grids/Container.svelte';
import { fetchProjectList, projects } from '../store';

import { goToProject } from '../project';


let grouped = [];
$: grouped = _.chunk($projects, 4);

onMount(fetchProjectList);

</script>

<Container>
<Row>
<Cell xsOffset={3} xs={9}>
  {#each grouped as chunk}
    <Row>
      {#each chunk as project}
        <Cell style='padding: 0; margin: 0;' xs=2>
          <div style='margin-bottom: 12px;'>
            <Button
              type='empty'
              on:click={goToProject(project.project_id)}
              style='padding: 0; margin: 0; display: block;'
            >
              <Block type='block3' style='width: 120px;'>
                {project.description.name}
              </Block>
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
  text-align: center;
}

</style>
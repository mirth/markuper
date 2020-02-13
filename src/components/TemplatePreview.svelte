<script>
import FieldRadio from './FieldRadio.svelte'

export let template;
</script>

<h1>{template.template.task}</h1>

<ul>
  {#each template.template.radios as field}
    <li>
      <FieldRadio {field} />
    </li>
  {/each}
</ul>
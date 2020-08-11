<script>
import { deterministicStrigify, getProjectIDFromSampleID } from '../../project';

export let assessedList;

function formatMarkup(markup) {
  if (markup === null) {
    return '[SKIPPED]';
  }

  return deterministicStrigify(markup);
}
</script>

<ul>
  {#each assessedList as forSample}
    <li>
      <p>
        <a href={`#/project/${getProjectIDFromSampleID(forSample.sample_id)}/samples/${forSample.sample_id}`}>
          <small>{forSample.sample_uri}:</small>
        </a>
        <small>{formatMarkup(forSample.sample_markup.markup)}</small>
      </p>
    </li>
  {/each}
</ul>

<style>
p {
  margin: 0;
}
small {
  font-size: 0.75em;
}
</style>
<script>
import { fetchProjectStats } from './../../store'

export let projectId;

$: projStats = fetchProjectStats(projectId)
</script>

{#await projStats then projStats}
  <p>
    ({projStats.assessed_number_of_samples} / {projStats.total_number_of_samples})
  </p>
{/await}

<style>
p {
  display: inline;
}
</style>
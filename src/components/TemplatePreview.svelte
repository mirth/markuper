<script>
import FieldRadio from './FieldRadio.svelte';

export let template;

const firstRadio = template.template.radios[0];
</script>

<FieldRadio field={firstRadio} />
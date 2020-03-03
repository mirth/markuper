<script>
import Typography from 'svelte-atoms/Typography.svelte';
import FieldRadio from './FieldRadio.svelte';

export let template;

const firstRadio = template.template.radios[0];
</script>

<Typography type="subheader" block>Add label:</Typography>
<FieldRadio field={firstRadio} />
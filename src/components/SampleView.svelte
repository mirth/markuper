<script>
import ImageView from './sample_view/ImageView.svelte';
import AudioView from './sample_view/AudioView.svelte';

export let sample;

</script>

{#if sample.sample.media_type === 'IMAGE_FILE_TYPE'}
  <ImageView {sample}/>
{/if}

{#if sample.sample.media_type === 'AUDIO_FILE_TYPE'}
  <AudioView {sample} />
{/if}

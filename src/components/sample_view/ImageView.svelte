<script>
import { sampleView, sampleMarkup, assessState } from '../../store';
import BoxLabels from './BoxLabels.svelte';


export let sample;

const field = sample.project.template.bounding_boxes[0];
let boxes = [];
let activeBox;
let resizeObserver;

function scaleBox(raw) {
  const scale = $assessState.imageElement.naturalWidth / $assessState.imageElement.clientWidth;
  return {
    x: raw.x / scale,
    y: raw.y / scale,
    width: raw.width / scale,
    height: raw.height / scale,
  };
}

function scaleBoxes(markupForGroup, assessStateBox) {
  if (!$assessState.imageElement) {
    return;
  }

  activeBox = assessStateBox ? scaleBox(assessStateBox) : null;
  boxes = (markupForGroup || []).map((mark) => ({
    ...mark,
    box: scaleBox(mark.box),
  }));
}

$: if (field) {
  scaleBoxes($sampleMarkup[field.group], $assessState.box);
}

$: if (field && $assessState.box) {
  scaleBoxes($sampleMarkup[field.group], $assessState.box);
}

$: if (field && !resizeObserver && $assessState.imageElement) {
  resizeObserver = new ResizeObserver(() => {
    scaleBoxes($sampleMarkup[field.group], $assessState.box);
  });
  resizeObserver.observe($assessState.imageElement);
}

</script>

<div id='image-container'>
  <img
    src='file://{sample.sample.media_uri}'
    alt='KEK'
    draggable=false
    bind:this={$assessState.imageElement}
    />

  {#each boxes as box, i}
    <div style={`
      width: ${box.box.width}px;
      height: ${box.box.height}px;
      left: ${box.box.x}px;
      top: ${box.box.y}px;
    `} class='box' class:box-selected={$sampleView.selectedBox === i}>
      <BoxLabels markup={box} template={sample.project.template} />
    </div>
  {/each}

  {#if activeBox}
    <div style={`
      width: ${activeBox.width}px;
      height: ${activeBox.height}px;
      left: ${activeBox.x}px;
      top: ${activeBox.y}px;
    `} class='box'/>
  {/if}
</div>

<style>

img {
  z-index:1;

  width: 100%;
  display:block;
  box-sizing: border-box;
}

.box {
  position: absolute;
  border: 2px solid green;

  box-sizing: border-box;
}

.box-selected {
  border: 4px solid green;
}

#image-container {
  position:relative;
  border: 2px solid black;
}

</style>
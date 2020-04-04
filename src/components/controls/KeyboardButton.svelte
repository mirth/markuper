<script>
export let key;
export let isKeyDown;

</script>

<kbd class:kbd-down={isKeyDown}>{key}</kbd>

<style>
kbd {
  display: inline-block;

  background-color: #eee;
  border-radius: 4px;
  font-size: 1em;
  padding: 0.2em 0.5em;
  border-top: 5px solid rgba(255,255,255,0.5);
  border-left: 5px solid rgba(255,255,255,0.5);
  border-right: 5px solid rgba(0,0,0,0.2);
  border-bottom: 5px solid rgba(0,0,0,0.2);
  color: #555;
}

.kbd-down {
  background-color: rgb(197, 50, 50);
}
</style>
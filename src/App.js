import svelte from 'svelte/compiler';
import App from './App.svelte';

const app = new App({
  target: document.body,
  props: {
    version: svelte.VERSION,
  },
});

export default app;

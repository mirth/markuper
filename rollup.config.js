import resolve from '@rollup/plugin-node-resolve';
import commonjs from 'rollup-plugin-commonjs'
import svelte from 'rollup-plugin-svelte';
import json from '@rollup/plugin-json';

export default {
  input: 'src/App.js',
  output: {
    file: 'public/bundle.js',
    format: 'iife',
    sourcemap: false,
  },
  plugins: [
    resolve({
      browser: true,
    }),
    commonjs(),
    svelte({
      css: function (css) {
        css.write('public/bundle.css');
      },
    }),
    json(),
  ]
}

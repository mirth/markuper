/* eslint-disable import/prefer-default-export */
import _ from 'lodash';
import { writable } from 'svelte/store';
import api from './api';

export const newProject = writable({
  name: '',
});

export const projects = writable([]);

export async function fetchProjectList() {
  const res = await api.get('/projects');
  const p = _.orderBy(res.projects, 'created_at', 'desc');

  projects.set(p);
}

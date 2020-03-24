/* eslint-disable import/prefer-default-export */
import _ from 'lodash';
import { writable } from 'svelte/store';
import api from './api';

export const projects = writable([]);
export const activeProject = writable({
  assessed: [],
  template: {
    fields_order: [],
    radios: [],
    checkboxes: [],
  },
  description: {
    name: '',
  },
  data_sources: [],
});
export const activeMarkup = writable({});
export const assessState = writable({
  focusedGroup: null,
});

export async function fetchProjectList() {
  const res = await api.get('/projects');
  const p = _.orderBy(res.projects, 'created_at', 'desc');

  projects.set(p);
}

export async function fetchProject(projectId) {
  const proj = await api.get(`/project/${projectId}`);
  const assessed = await api.get(`/project/${projectId}/assessed`);
  proj.assessed = assessed;

  activeProject.set(proj);
}

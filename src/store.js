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
  stats: {},
});
export const activeSample = writable({
  sample: null,
  project: {
    template: {
      fields_order: [],
    },
    description: {
      name: '',
    },
  },
});
export const sampleMarkup = writable({});
export const sampleView = writable({
  selectedBox: null,
});
export const assessState = writable({
  focusedGroup: null,
  lastTouchedGroup: null,
  box: null,
});

export const dataView = writable({
  imageElement: null,
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
  const stats = await api.get(`/project/${projectId}/stats`);
  proj.stats = stats;

  activeProject.set(proj);
}

export async function fetchProjectStats(projectId) {
  const stats = await api.get(`/project/${projectId}/stats`);

  return stats;
}

export function isFieldSelected(field, state) {
  return field.group === state.focusedGroup;
}

export async function fetchNextSample(projectID) {
  const sample = await api.get(`/project/${projectID}/next`);

  return sample;
}

/* eslint-disable import/prefer-default-export */
import { writable } from 'svelte/store';
import api from './api';

export const newProject = writable({
  name: '',
});

export const projects = writable([]);

export async function fetchProjectList() {
  const res = await api.get('/projects');
  projects.set(res.projects);
}

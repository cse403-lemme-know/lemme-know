import { writable } from 'svelte/store';

export const startDate = writable(undefined);
export const endDate = writable(undefined);
export const groupName = writable('');
export const groupId = writable(undefined);
export const userId = writable(undefined);

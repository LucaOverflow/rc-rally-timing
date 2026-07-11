import PocketBase from 'pocketbase';

const BASE_URL = import.meta.env.DEV ? 'http://127.0.0.1:8090/' : '/';
export const pb = new PocketBase(BASE_URL);
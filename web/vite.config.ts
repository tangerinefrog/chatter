import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';
import { loadEnv } from 'vite';
import path from 'path';

export default defineConfig(({mode}) => {
	const env = loadEnv(mode, path.resolve('../'), '');
	const url = new URL(env.WEB_ADDR);
	const port = Number(url?.port);

	return {
		plugins: [sveltekit()],
		server: {
			port: port ? port : 5473
		}
	}
});

import { sveltekit } from '@sveltejs/kit/vite';

/** @type {import('vite').UserConfig} */
const config = {
	plugins: [sveltekit()],
	test: {
		include: ['src/**/*.{test,spec}.{js,ts}']
	},
	preview: {
		port: 5173
	},
	build: {
		// because we build and preview to debug.
		minify: false
	}
};

export default config;

import { sveltekit } from '@sveltejs/kit/vite';

/** @type {import('vite').UserConfig} */
const config = {
	plugins: [
		sveltekit({
			onwarn: (warning, handler) => {
				if (
					warning.message &&
					warning.message.includes(
						'visible, non-interactive elements with an on:click event must be accompanied by an on:keydown, on:keyup, or on:keypress event'
					)
				) {
					return;
				}
				handler(warning);
			}
		})
	],
	test: {
		include: ['src/**/*.{test,spec}.{js,ts}']
	},
	preview: {
		port: 5173
	}
};

export default config;

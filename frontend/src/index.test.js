import { describe, it, expect } from 'vitest';
import { render } from '@testing-library/svelte';
import { mockPageFunction, navigate, checkInput } from './__mocks__/MockPage'

describe('Mocking page.svelte', () => {
	it('adds 1 + 2 to equal 3', () => {
		expect(1 + 2).toBe(3);
	});

	it('renders pages', async () => {
		const result = mockPageFunction();
		expect(result).toBe(true);
	});

	it('navigates to pages', async () => {
		const message = navigate('/dashboard');
		expect(message).toBe('Navigated to /dashboard');
	});

	it('checking for valid input', async () => {
		expect(checkInput('test')).toBe(true);
	});

	it('checking for empty input', async () => {
		expect(checkInput(' ')).toBe(false);
	});
});

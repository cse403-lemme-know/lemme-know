import { expect, test } from '@playwright/test';

test('Page has header and message', async ({ page }) => {
	await page.goto('http://localhost:8080/');
	expect(await page.textContent('h1')).toBe('LemmeKnow');
	expect(await page.textContent('h3')).toContain(
		'SCHEDULE HANGOUTS, PLAN ROAD-TRIPS, SHARE CALENDARS, EVERYTHING, EVERYWHERE, ALL AT ONCE.'
	);
});

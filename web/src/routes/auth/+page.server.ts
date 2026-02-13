import { fail, redirect } from '@sveltejs/kit';
import type { Actions } from './$types';

const BACKEND_URL = 'http://localhost:8080';

export const actions: Actions = {
	login: async ({ request, fetch }) => {
		const data = await request.formData();
		const username = data.get('username');
		const password = data.get('password');

		const res = await fetch(`${BACKEND_URL}/login`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify({ username, password }),
			credentials: 'include'
		});

		if (!res.ok) {
			const body = await res.json().catch(() => ({}));
			return fail(res.status, { error: body.message ?? 'Invalid credentials' });
		}

		throw redirect(303, '/');
	},

	signup: async ({ request, fetch }) => {
		const data = await request.formData();
		const username = data.get('username');
		const password = data.get('password');

		const res = await fetch(`${BACKEND_URL}/signup`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify({ username, password }),
			credentials: 'include'
		});

		if (!res.ok) {
			const body = await res.json().catch(() => ({}));
			return fail(res.status, { error: body.message ?? 'Signup failed' });
		}

		throw redirect(303, '/');
	}
};
import { redirect } from '@sveltejs/kit';
import type { Handle } from '@sveltejs/kit';

export const handle: Handle = async ({ event, resolve }) => {
	const jwt = event.cookies.get('auth_token');
	const path = event.url.pathname;

	const isAuthPage = path.startsWith('/auth');

	if (!isAuthPage && !jwt) {
		throw redirect(303, '/auth');
	}

	if (isAuthPage && jwt) {
		throw redirect(303, '/');
	}

	return resolve(event);
};
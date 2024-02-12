async function getUser() {
	try {
		const response = await fetch(`//${location.host}/api/user/`);
		const user = await response.json();
		return user;
	} catch (e) {
		return null;
	}
}

async function createGroup(name) {
	try {
		const response = await fetch(`//${location.host}/api/group/`, {
			method: 'PATCH',
			body: JSON.stringify({ name })
		});
		const result = await response.json();
		return result.groupId;
	} catch (e) {
		return null;
	}
}

async function createAvailability(groupId, availability) {
	try {
		const response = await fetch(`//${location.host}/api/group/${groupId}/availability/`, {
			method: 'PATCH',
			headers: {
				'Content-Type': 'application/json',
			},
			body: JSON.stringify(availability),
			credentials: 'include',
		});
		if (!response.ok) {
			throw new Error('response not ok');
		}
		return await response.json();
	} catch (e) {
		console.error('Error creating availability:', e);
		return null;
	}
}

getUser().then(console.log);

export { getUser, createGroup, createAvailability };

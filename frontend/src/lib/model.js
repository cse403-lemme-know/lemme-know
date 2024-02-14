// @ts-nocheck
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
				'Content-Type': 'application/json'
			},
			body: JSON.stringify(availability)
		});
		if (response.status === 200) {
			console.log('success for creating availability');
		}
	} catch (e) {
		console.error('Error creating availability:', e);
		return null;
	}
}

async function createTask(groupId, title) {
	try {
		return await fetch(`//${location.host}/api/group/${groupId}/task/`, {
			method: 'PATCH',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({ title })
		});
	} catch (e) {
		console.error('Error creating task:', e);
		return null;
	}
}

async function createPoll(groupId, title, options) {
	try {
		const response = await fetch(`//${location.host}/api/group/${groupId}/poll/`, {
			method: 'PATCH',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({ title, options })
		});
		const result = await response.json();
	} catch (e) {
		return null;
	}
}

async function updateVotes(groupID, votes) {
	try {
		const response = await fetch(`//${location.host}/api/group/${groupID}/poll/`, {
			method: 'PATCH',
			body: JSON.stringify({ votes })
		});
		const result = await response.json();
	} catch (e) {
		return null;
	}
}

async function deletePoll(groupID) {
	try {
		const response = await fetch(`//${location.host}/api/group/${groupID}/poll/`, {
			method: 'DELETE'
		});
		const result = await response.json();
	} catch (e) {
		return null;
	}
}

async function sendMessage(groupID, content) {
	try {
		const response = await fetch(`//${location.host}/api/group/${groupID}/chat/`, {
			method: 'PATCH',
			body: JSON.stringify({ content })
		});
		const result = await response.json();
	} catch (e) {
		return null;
	}
}

async function fetchMessages(groupID, start, end) {
	try {
		const response = await fetch(`//${location.host}/api/group/${groupID}/chat/?` + new URLSearchParams({ start, end}), {
			method: 'GET',
		});
		const result = await response.json();
		if (result.continue == true) {
			result.messages[result.messages.length - 1].timestamp + 1;
		}
	} catch (e) {
		return null;
	}
}

getUser().then(console.log);

export { getUser, createGroup, 
	createAvailability, createTask, 
	createPoll, updateVotes, deletePoll, 
	sendMessage, fetchMessages};

<script>
	// @ts-nocheck

	import { onMount } from 'svelte';
	import dayjs from 'dayjs';
	import { writable, get } from 'svelte/store';
	import { startDate, endDate, groupId } from '$lib/stores';
	import { createAvailability, createTask, fetchMessages, sendMessage } from '$lib/model';
	import PollCreationModal from './PollCreationModal.svelte';

	let start, end;
	let availableTimes = [];
	let availability = writable({});
	let successMsg = writable('');

	let tasks = writable([]);
	let taskMsg = writable('');
	let taskInput = '';
	let assignedInput = '';

	let isPoll = writable(false);

	onMount(async () => {
		start = get(startDate);
		end = get(endDate);

		start = dayjs(start);
		end = dayjs(end);

		function initializeAvailability(start, end) {
			let days = {};

			for (
				let current = start;
				current.isBefore(end.add(1, 'day'));
				current = current.add(1, 'day')
			) {
				const dateString = current.format('YYYY-MM-DD');
				days[dateString] = new Array(24).fill(false);
			}
			availability.set(days);
		}

		if (start.isValid() && end.isValid()) {
			initializeAvailability(start, end);
		} else {
			console.error('Invalid start or end date');
		}
	});

	function toggleSlot(day, hour) {
		availability.update((a) => {
			a[day][hour] = !a[day][hour];
			return a;
		});
	}

	async function addTask(taskDescription, assigneeName) {
		const currentGroup = get(groupId);
		console.log('adding task for group: ', currentGroup);
		if (!currentGroup) {
			taskMsg.set('No Group ID is set');
			return;
		}

		if (!taskDescription.trim()) {
			taskMsg.set('Task description is required.');
			return;
		}
		taskMsg.set('');

		try {
			const response = await createTask(currentGroup, taskDescription, assigneeName);
			if (response.ok) {
				tasks.update((currentTasks) => {
					const newTask = {
						id: currentTasks.length + 1,
						description: taskDescription,
						assignedTo: assigneeName,
						completed: false
					};
					return [...currentTasks, newTask];
				});
				taskInput = '';
				assignedInput = '';
				taskMsg.set(`Task added: ${taskDescription}`);
			} else {
				taskMsg.set(`Failed to add task: server error`);
			}
		} catch (e) {
			taskMsg.set('Failed to add task');
			console.error('task error ', e);
		}
	}

	/** @param {number} taskId */
	function toggleCompletion(taskId) {
		tasks.update((currentTasks) => {
			const index = currentTasks.findIndex((t) => t.id === taskId);
			if (index !== -1) {
				currentTasks[index].completed = !currentTasks[index].completed;
			}
			return currentTasks;
		});
	}
	// for chat box
	let messages = [];
	let newMessage = '';
	/**
	 * Send a message to the chat.
	 */
	function handleSendMessage() {
		if (newMessage.trim() !== '') {
			messages = [...messages, { text: newMessage, sender: 'user' }];
			content = newMessage.trim();
			sendMessage($groupId, content);
			newMessage = '';
			// need to handle Fetch messages
		}
	}
	/**
	 * send message if hit enter
	 * @param event
	 */
	function handleKeyPress(event) {
		if (event.key === 'Enter') {
			sendMessage();
		}
	}

	function openPoll() {
		isPoll.set(true);
	}

	async function saveAllAvailabilities() {
		const currentGroupId = $groupId;
		console.log('current group ', currentGroupId);
		if (!currentGroupId) {
			console.error('No group ID is set.');
			return;
		}

		const allAvailabilityData = [];

		for (const [date, slots] of Object.entries($availability)) {
			slots.forEach((slot, hour) => {
				if (slot) {
					const timeId = `${date}_${hour < 10 ? `0${hour}` : hour}:00`;
					if (!availableTimes.includes(timeId)) {
						allAvailabilityData.push({
							date: date,
							start: `${hour}:00`,
							end: `${hour + 1}:00`
						});
						availableTimes.push(timeId);
					}
				}
			});
		}

		try {
			for (const availabilityData of allAvailabilityData) {
				createAvailability(currentGroupId, availabilityData);
			}

			const times = allAvailabilityData
				.map((data) => `${data.date} from ${data.start} to ${data.end}`)
				.join(', ');
			successMsg.set('All availabilities saved successfully ' + times);
			console.log('Saved times:', JSON.stringify(availableTimes));
		} catch (error) {
			successMsg.set('Failed to save availability');
			console.error('Failed to save availability with error', error);
			availableTimes = {};
		}
	}
</script>

<header />

<main>
	<div class="content-wrap">
		<div class="menu-bar">
			<button class="menu-button">
				<img src="../menubar.png" alt="menu bar" class="hamburger-icon" />
				<span class="logo">LemmeKnow</span>
			</button>
			<button class="menu-button">
				<img src="../users.png" alt="menu bar" class="user-icon" />
				<span class="members-title">Members</span>
			</button>
			<button class="menu-button" on:click={openPoll}>
				<img src="../users.png" alt="menu bar" class="user-icon" />
				<span class="members-title">Create Poll</span>
			</button>
		</div>

		<div class="chatbox">
			<div class="messages">
				{#each messages as message (message.text)}
					<div class:message class:message.sender={message.sender}>
						{#if message.sender === 'user'}
							<strong class="user-message">You:</strong> {message.text}
						{:else if message.sender === 'system'}
							<em class="system-message">{message.text}</em>
						{/if}
					</div>
				{/each}
			</div>
			<div class="poll">
				{#if $isPoll}
					<PollCreationModal />
				{/if}
			</div>
			<div class="input-bar">
				<input
					class="input"
					bind:value={newMessage}
					placeholder="Type your message..."
					on:keydown={handleKeyPress}
				/>
				<button on:click={handleSendMessage} on:keyup={handleSendMessage}>Send Message</button>
			</div>
		</div>

		<div class="calendar-container">
			<span class="calendar-title">AVAILABILITY CALENDAR</span>
			{#each Object.keys($availability) as day}
				<div class="day">
					<h3>{day}</h3>
					<div class="slots">
						{#each $availability[day] as available, hour}
							<div
								class="slot {available ? 'available' : ''}"
								on:click|preventDefault={() => toggleSlot(day, hour)}
								on:keypress={() => toggleSlot(day, hour)}
							>
								{hour}:00
							</div>
						{/each}
					</div>
				</div>
			{/each}
			<button on:click={saveAllAvailabilities}>Save Availability</button>
			{#if $successMsg}
				<p>{$successMsg}</p>
			{/if}
			<form on:submit|preventDefault={() => addTask(taskInput, assignedInput)}>
				<input
					type="text"
					bind:value={taskInput}
					placeholder="Enter task description (50 characters max)"
					maxlength="50"
				/>
				<input
					type="text"
					bind:value={assignedInput}
					placeholder="Enter assignee name (50 characters max)"
					maxlength="50"
				/>
				<button type="submit" disabled={!taskInput.trim()}>Add Task</button>

				{#if $taskMsg}
					<p>{$taskMsg}</p>
				{/if}
			</form>
			{#each $tasks as task (task.id)}
				<div class="task-item">
					<input
						type="checkbox"
						bind:checked={task.completed}
						on:click={() => toggleCompletion(task.id)}
						on:keypress={() => toggleCompletion(task.id)}
					/>
					<span class={task.completed ? 'completed-task' : ''}>{task.description}</span>
					{#if task.assignedTo}
						<span>Assigned to: {task.assignedTo}</span>
					{/if}
				</div>
			{/each}
		</div>
	</div>
</main>

<style>
	.calendar-title {
		display: flex;
		align-items: flex-start;
		justify-content: center;
		flex-direction: column;
		text-align: center;
		font-size: 3rem;
		margin-top: 0.25rem;
		font-family: 'Baloo Bhai 2';
		margin-left: 1rem;
		font-weight: bolder;
		color: black;
	}
	.menu-bar {
		position: relative;
		display: flex;
		flex-direction: column;
		top: 0;
		left: 0;
		align-items: flex-start;
	}

	.user-icon {
		width: 3rem;
		display: block;
		margin-left: 1.5rem;
	}

	.menu-button {
		background: none;
		border: none;
		cursor: pointer;
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		line-height: 1;
		text-align: center;
	}

	.menu-button:focus {
		outline: none;
	}

	.logo {
		font-size: 1.5rem;
		margin-top: 0.25rem;
		font-family: 'Baloo Bhai 2';
		font-weight: bolder;
		color: #73a0e7;
	}

	img.hamburger-icon {
		width: 3rem;
		margin: 0 auto;
		display: block;
	}

	.members-title {
		display: flex;
		align-items: flex-start;
		justify-content: center;
		text-align: center;
		font-size: 1.5rem;
		margin-top: 0.25rem;
		margin-left: 1rem;
		font-family: 'Baloo Bhai 2';
		font-weight: bolder;
		color: black;
	}

	.menu-button:hover .hamburger-icon {
		transform: scale(1.2);
	}

	.menu-button:hover .user-icon {
		transform: scale(1.2);
	}
	.menu-button:hover .poll-icon {
		transform: scale(1.2);
	}
	.content-wrap {
		display: flex;
		flex-direction: row;
		gap: 2rem;
	}

	.calendar-container {
		display: flex;
		flex-direction: column;
		flex-wrap: wrap;
		margin-left: 2rem;
		margin-top: 3rem;
	}

	.day {
		border: 1px solid #ccc;
		padding: 10px;
	}

	.slots {
		display: grid;
		grid-template-columns: repeat(4, 1fr);
		gap: 5px;
	}

	.slot {
		padding: 5px;
		background-color: #f0f0f0;
		text-align: center;
		cursor: pointer;
	}

	.slot.available {
		background-color: #a9e1a9;
	}

	input {
		padding: 0.5rem;
		margin-bottom: 0.5rem;
		margin-top: 0.5rem;
		width: 80%;
		max-width: 15rem;
		text-align: center;
		font-size: 1rem;
		background-color: #c9e7e7;
		border-radius: 15px;
		border: 2px solid transparent;
	}

	button {
		padding: 0.5rem 1rem;
		background-color: #2774d0;
		color: white;
		border: none;
		align-items: center;
		font-size: 1rem;
		border-radius: 1rem;
		cursor: pointer;
	}

	button:hover {
		background-color: gray;
		color: white;
	}

	.task-item {
		display: flex;
		margin-bottom: 1rem;
		margin-top: 0.5rem;
		padding: 0.5rem;
		background-color: #f9f9f9;
		border-radius: 0.2rem;
		font-family: 'Baloo Bhai 2';
		align-items: center;
		font-size: 1.25rem;
	}

	.task-item input[type='checkbox'] {
		accent-color: #879db7;
		transform: scale(1.5);
		cursor: pointer;
	}

	.task-item .completed-task {
		text-decoration: line-through;
		color: #879db7;
	}

	.task-item span {
		margin-left: 1rem;
		color: #333;
		font-weight: bold;
	}

	button[type='submit']:disabled {
		background-color: #ccc;
		cursor: not-allowed;
	}
	/* chatbox style */
	.chatbox {
		display: flex;
		flex-direction: column;
		border: 2px solid #ccc;
		padding: 4rem 6rem 2rem 6rem;
		max-width: calc(90% - 10px);
		height: 700px;
		border-radius: 8px;
		box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
		overflow-y: auto; /* Add scrollbar when content exceeds the height */
		margin-right: 2rem;
	}

	.messages {
		flex-grow: 1;
		display: flex;
		flex-direction: column-reverse; /* Reverse the order of messages */
	}

	.message {
		margin: 8px 0;
		padding: 8px;
		background-color: #f0f0f0;
		border-radius: 4px;
	}

	.user-message {
		background-color: #e6f7ff;
		text-align: right;
	}

	.system-message {
		color: #888;
		font-style: italic;
	}

	.input-bar {
		display: flex;
		align-items: center;
		justify-content: space-between;
		margin-top: 10px;
	}
	input {
		padding: 0.5rem;
		margin-bottom: 0.5rem;
		width: 80%;
		max-width: 300px;
		text-align: center;
		font-size: 1rem;
		background-color: #eedaf4;
		border-radius: 15px;
		border: 2px solid transparent;
	}

	button {
		flex-shrink: 0;
	}
</style>

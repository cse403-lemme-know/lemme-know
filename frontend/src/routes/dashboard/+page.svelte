<script>
// @ts-nocheck

	import { onMount } from 'svelte';
	import dayjs from 'dayjs';
	import { writable, get } from 'svelte/store';
	import { startDate, endDate } from '$lib/stores';
  import PollCreationModal from './PollCreationModal.svelte';
	import {sendMessage, fetchMessages} from '$lib/model';

	let start;
	let end;

	let availability = writable({});
	let tasks = writable([]);
	let taskInput = '';
	let assignedInput = '';

	let isPollCreationModalOpen = writable(false);

	onMount(() => {
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
		initializeAvailability(dayjs(start), dayjs(end));

	});

	function toggleSlot(day, hour) {
		availability.update((a) => {
			a[day][hour] = !a[day][hour];
			return a;
		});
	}

	function addTask(taskDescription, assigneeName) {
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
	let messages = writable([]);
	let newMessage = '';
	/**
	 * Send a message to the chat.
	 */
	function sendMessages() {
		if (newMessage.trim() !== '') {
			messages = [...messages, { text: newMessage, sender: 'user' }];
			newMessage = '';
			// add logic here to handle the response from a server or another user.

		}
	}
	/**
	 * send message if hit enter
	 * @param event
	 */
	function handleKeyPress(event) {
		if (event.key === 'Enter') {
			sendMessages();
		}
	}
	// for poll
  function openPollCreationModal() {
		isPollCreationModalOpen.set(true);
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
			<button class="menu-button" on:click={openPollCreationModal}>
				<img src="../poll.png" alt="menu bar" class="poll-icon" />
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
			<!-- poll on and off -->
			{#if $isPollCreationModalOpen}
				<PollCreationModal />
			{/if}

			<div class="input-bar">
				<input
					class="input"
					bind:value={newMessage}
					placeholder="Type your message..."
					on:keydown={handleKeyPress}
				/>
				<button on:click={sendMessages} on:keyup={sendMessages}>Send Message</button>
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
								on:click={() => toggleSlot(day, hour)}
								on:keypress={() => toggleSlot(day, hour)}
							>
								{hour}:00
							</div>
						{/each}
					</div>
				</div>
			{/each}
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
				<button type="submit" disabled={!taskInput.trim() || !assignedInput.trim()}>Add Task</button
				>
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

	.poll-icon {
		align-items: center;
		justify-content: center;
		width: 3rem;
		display: block;

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
		display: flex;
		padding: 0;
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

	.content-wrap {
		display: flex;
		flex-direction: row;
	}

	.calendar-container {
		display: flex;
		flex-direction: column;
		flex-wrap: wrap;
		margin-left: 4rem;
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
		padding: 10px;
		width: 700px;
		height: 700px;
		margin: auto;
		border-radius: 8px;
		box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
		overflow-y: auto; /* Add scrollbar when content exceeds the height */
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
</style>

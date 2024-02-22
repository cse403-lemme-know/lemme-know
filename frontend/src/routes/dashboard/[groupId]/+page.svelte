<script>
	// @ts-nocheck

	import { onMount } from 'svelte';
	import dayjs from 'dayjs';
	import { get, writable } from 'svelte/store';
	import {
		createAvailability,
		createTask,
		deleteAvailability,
		deleteTask,
		getGroup,
		groups,
		refreshGroup
	} from '$lib/model';
	import { goto } from '$app/navigation';
	import Chat from './Chat.svelte';
	import { page } from '$app/stores';

	$: groupId = $page.params.groupId;

	let start, end;
	let availableTimes = [];
	let availability = writable({});
	let successMsg = writable('');
	$: group = $groups[groupId];
	let groupData = {};
	$: group = $groups[groupId];

	let tasks = writable([]);
	let taskMsg = writable('');
	let taskInput = '';
	let assignedInput = '';
	let isPoll = false;

	onMount(async () => {
		// TODO: Refactor to avoid needing this.
		let g = group;
		if (!g) {
			await refreshGroup(groupId);
			g = get(groups)[groupId];
			if (!g) {
				goto('/');
				return;
			}
		}
		const calendarMode = g.calendarMode.split(' to ');
		const dateFormat = 'YYYY-MM-DD';
		console.log(calendarMode);

		start = dayjs(calendarMode[0], dateFormat);
		end = dayjs(calendarMode[1], dateFormat);

		function initializeAvailability(start, end) {
			let days = {};
			let loopEndDate = end.add(1, 'day');
			for (let current = start; current.isBefore(loopEndDate); current = current.add(1, 'day')) {
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

	async function addTask(title) {
		console.log('adding task for group: ', groupId);
		if (!groupId) {
			taskMsg.set('No Group ID is set');
			return;
		}

		// if (!taskDescription.trim()) {
		// 	taskMsg.set('Task description is required.');
		// 	return;
		// }
		taskMsg.set('');

		try {
			const response = await createTask(groupId, title);
			if (response.ok) {
				await updateGroupData(groupId);
				tasks.set(
					groupData.tasks.map((task) => ({
						id: task.taskId,
						description: task.title,
						assignedTo: task.assignee,
						completed: task.complete
					}))
				);
				taskInput = '';
				assignedInput = '';
				taskMsg.set(`Task added: ${title}`);
			} else {
				taskMsg.set(`Failed to add task: server error`);
			}
		} catch (e) {
			taskMsg.set('Failed to add task');
			console.error('task error ', e);
		}
		await updateGroupData(groupId);
	}

	function toggleCompletion(taskId) {
		tasks.update((currentTasks) => {
			const index = currentTasks.findIndex((t) => t.id === taskId);
			if (index !== -1) {
				currentTasks[index].completed = !currentTasks[index].completed;
			}
			return currentTasks;
		});
	}

	function openPoll() {
		isPoll = true;
	}

	async function saveAllAvailabilities() {
		const currentGroupId = groupId;
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
				createAvailability(groupId, availabilityData);
			}

			const times = allAvailabilityData
				.map((data) => `${data.date} from ${data.start} to ${data.end}`)
				.join(', ');
			successMsg.set('All availabilities saved successfully ' + times);
			console.log('GroupID', groupId);
			console.log('Saved times:', JSON.stringify(availableTimes));
		} catch (error) {
			successMsg.set('Failed to save availability');
			console.error('Failed to save availability with error', error);
			availableTimes = {};
		}
	}

	async function removeAvailability(selectedDay, selectedHour) {
		const formattedHour = `${selectedHour < 10 ? `0${selectedHour}` : selectedHour}:00`;
		const currentData = await getGroup(groupId);
		const matchingAvailability = currentData.availabilities.find(
			(avail) => avail.date === selectedDay && avail.start === formattedHour
		);

		console.log(groupId);
		if (matchingAvailability) {
			await deleteAvailability(groupId, matchingAvailability.availabilityId);
			console.log(
				'making an attempt to delete availability with id: ',
				matchingAvailability.availabilityId
			);
			await updateGroupData(groupId);
			console.log(`Deleted availability with ID: ${matchingAvailability.availabilityId}`);
		} else {
			console.error('No matching availability found to delete');
		}
	}

	async function updateGroupData(groupId) {
		try {
			groupData = await getGroup(groupId);
			console.log('group after update: ', groupData);
		} catch (e) {
			console.error(e);
		}
	}

	async function deleteTaskWrapper(taskId) {
		try {
			await deleteTask(groupId, taskId);
			tasks.update((currentTasks) => {
				return currentTasks.filter((task) => task.id !== taskId);
			});
		} catch (error) {
			console.error(error);
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
				<img src="../poll.png" alt="menu bar" class="user-icon" />
				<span class="members-title">Create Poll</span>
			</button>
			<button
				on:click={() =>
					navigator.clipboard.writeText(`${window.location.origin}/dashboard/${get(groupId)}`)}
				class="invite-button">Invite Link!</button
			>
		</div>

		<Chat {groupId} {group} bind:isPoll />

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
								{#if available}
									<button on:click|preventDefault={() => removeAvailability(day, hour)}
										>Delete</button
									>
								{/if}
							</div>
						{/each}
					</div>
				</div>
			{/each}
			<button on:click={saveAllAvailabilities}>Save Availability</button>
			{#if $successMsg}
				<p>{$successMsg}</p>
			{/if}
			<form on:submit|preventDefault={() => addTask(taskInput)}>
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
						<span class={task.completed ? 'completed-task' : ''}
							>Assigned to: {task.assignedTo}</span
						>
					{/if}
					<button class="delete-task" on:click={() => deleteTaskWrapper(task.id)}>delete</button>
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

	:global(button) {
		padding: 0.5rem 1rem;
		background-color: #2774d0;
		color: white;
		border: none;
		align-items: center;
		font-size: 1rem;
		border-radius: 1rem;
		cursor: pointer;
	}

	:global(button:hover) {
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
		margin-left: -7.5rem;
		margin-right: -7.5rem;
	}

	.task-item .completed-task {
		text-decoration: line-through;
		color: #879db7;
	}

	.task-item span {
		margin-right: 1rem;
		color: #333;
		text-align: center;
		font-weight: bold;
	}

	:global(button[type='submit']:disabled) {
		background-color: #ccc;
		cursor: not-allowed;
	}

	:global(input) {
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

	:global(button) {
		flex-shrink: 0;
	}

	.invite-button {
		display: block;
		margin: 1rem auto;
		background-color: #76a6e7;
		font-weight: bolder;
		font-family: 'Baloo Bhai 2';
		font-size: large;
		color: black;
	}

	.invite-button:hover {
		background-color: #afaeae;
		color: white;
	}

	.delete-task {
		background-color: #879db7;
		color: black;
		border: none;
		cursor: pointer;
		margin-left: 1.5rem;
		padding: 0.5rem 1rem;
		display: inline-block;
		text-align: center;
		font-size: 1rem;
		border-radius: 0.3rem;
	}

	.delete-task:hover {
		background-color: gray;
		color: white;
	}
</style>

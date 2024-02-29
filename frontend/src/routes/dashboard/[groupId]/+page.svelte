<script>
	// @ts-nocheck

	import { onMount, onDestroy } from 'svelte';
	import dayjs from 'dayjs';
	import { get, writable } from 'svelte/store';
	import {
		createAvailability,
		createTask,
		deleteAvailability,
		deleteTask,
		getGroup,
		groups,
		refreshGroup,
		updateTask,
		userId
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

	let tasks = writable([]);
	let taskInput = '';
	let isPoll = false;
	let commonAvailability = writable({ isLoading: true, slots: [] });
	let initialLoad = true;
	let loadingTimeout;
	let lastCalculated = [];

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

		start = dayjs(calendarMode[0], dateFormat);
		end = dayjs(calendarMode[1], dateFormat);

		function initializeAvailability(start, end) {
			let days = {};
			let loopEndDate = end.add(1, 'day');
			for (let current = start; current.isBefore(loopEndDate); current = current.add(1, 'day')) {
				const dateString = current.format('YYYY-MM-DD');
				days[dateString] = new Array(16).fill(false);
			}
			availability.set(days);
		}

		if (start.isValid() && end.isValid()) {
			initializeAvailability(start, end);
			await loadExistingAvailabilities();
			await loadTasks(groupId);
		} else {
			console.error('Invalid start or end date');
		}

		const pollingInterval = setInterval(async () => {
			await getGroup(groupId);
			await calculateCommonAvailability();
		}, 1000);

		onDestroy(() => {
			clearInterval(pollingInterval);
		});
	});

	async function loadExistingAvailabilities() {
		const currentUser = get(userId);
		const groupData = await getGroup(groupId);
		if (groupData && groupData.availabilities) {
			const userAvailabilities = groupData.availabilities.filter(
				(avail) => avail.userId === currentUser
			);
			availability.update((a) => {
				userAvailabilities.forEach(({ date, start }) => {
					const hour = parseInt(start.split(':')[0], 10) - 7;
					if (a[date]) {
						a[date][hour] = true;
					}
				});
				return a;
			});
		}
	}

	function slotsChanged(newSlots, oldSlots) {
		if (newSlots.length !== oldSlots.length) return true;
		for (let i = 0; i < newSlots.length; i++) {
			if (newSlots[i] !== oldSlots[i]) {
				return true;
			}
		}
		return false;
	}

	async function calculateCommonAvailability() {
		if (initialLoad) {
			commonAvailability.set({ isLoading: true, slots: [] });
		}
		const groupData = await getGroup(groupId);
		let availabilityRanges = {};

		if (groupData && groupData.availabilities) {
			groupData.availabilities.forEach(({ userId, date, start, end }) => {
				const startTime = dayjs(`${date} ${start}`);
				const endTime = dayjs(`${date} ${end}`);
				if (!availabilityRanges[date]) {
					availabilityRanges[date] = [];
				}
				availabilityRanges[date].push({
					userId: userId,
					start: startTime,
					end: endTime
				});
			});

			let commonSlots = {};
			Object.keys(availabilityRanges).forEach((date) => {
				let ranges = availabilityRanges[date];
				ranges.sort((a, b) => a.start - b.start);

				let overlapping = [];
				ranges.forEach((currentRange, index) => {
					for (let i = index + 1; i < ranges.length; i++) {
						let nextRange = ranges[i];
						if (currentRange.end > nextRange.start) {
							let startOverlap = nextRange.start.isAfter(currentRange.start)
								? nextRange.start
								: currentRange.start;
							let endOverlap = nextRange.end.isBefore(currentRange.end)
								? nextRange.end
								: currentRange.end;
							if (
								!overlapping.some(
									(ov) => ov.start.isSame(startOverlap) && ov.end.isSame(endOverlap)
								)
							) {
								overlapping.push({ start: startOverlap, end: endOverlap });
							}
						}
					}
				});

				if (overlapping.length) {
					commonSlots[date] = overlapping;
				}
			});

			let formattedCommonSlots = [];
			Object.keys(commonSlots).forEach((date) => {
				commonSlots[date].forEach((slot) => {
					formattedCommonSlots.push(
						`${date} from ${slot.start.format('HH:mm')} to ${slot.end.format('HH:mm')}`
					);
				});
			});

			if (slotsChanged(formattedCommonSlots, lastCalculated)) {
				commonAvailability.set({ isLoading: false, slots: formattedCommonSlots });
				lastCalculated = formattedCommonSlots;
			} else {
				if (!initialLoad) {
					clearTimeout(loadingTimeout);
					commonAvailability.update((values) => ({ ...values, isLoading: false }));
				}
			}
		} else {
			if (!initialLoad) {
				clearTimeout(loadingTimeout);
				commonAvailability.update((values) => ({ ...values, isLoading: false }));
			}
		}

		initialLoad = false;
	}

	async function loadTasks(groupId) {
		try {
			const groupData = await getGroup(groupId);
			if (groupData && groupData.tasks) {
				tasks.set(
					groupData.tasks.map((task) => ({
						taskId: task.taskId,
						title: task.title,
						assignee: task.assignee,
						completed: task.completed
					}))
				);
			}
		} catch (error) {
			console.error('Failed to load tasks', error);
			tasks.set([]);
		}
	}

	function toggleSlot(day, hour) {
		availability.update((a) => {
			a[day][hour] = !a[day][hour];
			return a;
		});
	}

	async function addTask(title) {
		if (!groupId) {
			return;
		}

		try {
			const response = await createTask(groupId, title);
			if (response.ok) {
				await updateGroupData(groupId);
				tasks.set(
					groupData.tasks.map((task) => ({
						taskId: task.taskId,
						title: task.title,
						assignee: task.assignee,
						completed: task.completed
					}))
				);
				taskInput = '';
			} else {
				console.error('server error');
			}
		} catch (e) {
			console.error('task error ', e);
		}
		await updateGroupData(groupId);
	}

	async function toggleCompletion(taskId) {
		const task = $tasks.find((t) => t.taskId === taskId);
		if (task) {
			try {
				const newCompletedStatus = !task.completed;
				console.log('status', newCompletedStatus);
				const success = await updateTask(groupId, taskId, { completed: newCompletedStatus });

				if (success) {
					tasks.update((currentTasks) => {
						return currentTasks.map((t) =>
							t.taskId === taskId ? { ...t, completed: newCompletedStatus } : t
						);
					});
				} else {
					console.error('Failed to update task completion on server.');
				}
				await updateGroupData(groupId);
			} catch (error) {
				console.error('Error updating task completion:', error);
			}
		}
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
				hour = hour + 7;
				if (slot) {
					const timeId = `${date}_${hour < 10 ? `0${hour}` : hour}:00`;
					if (!availableTimes.includes(timeId)) {
						allAvailabilityData.push({
							date: date,
							start: `${hour < 10 ? `0${hour}` : hour}:00`,
							end: `${hour + 1 < 10 ? `0${hour + 1}` : hour + 1}:00`
						});
						availableTimes.push(timeId);
					}
				}
			});
		}
		console.log('Availabiltiy', availability);

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
				return currentTasks.filter((task) => task.taskId !== taskId);
			});
		} catch (error) {
			console.error(error);
		}
	}

	async function assignTaskToUser(taskId) {
		console.log('taskid:', taskId);
		console.log('user id:', userId);
		const taskData = {
			assignee: $userId
		};
		const success = await updateTask(groupId, taskId, taskData);
		if (success) {
			tasks.update((currentTasks) => {
				return currentTasks.map((t) => {
					if (t.taskId === taskId) {
						return { ...t, assignee: $userId };
					}
					return t;
				});
			});
			console.log('Task assigned');
		} else {
			console.error('Failed to assign task.');
		}
		console.log('before: ', groupData);
		await updateGroupData(groupId);
		console.log('after: ', groupData);
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
				on:click={() => {
					navigator.clipboard.writeText(`${window.location.origin}/dashboard/${groupId}`);
					document.querySelector('.invite-button').innerText = 'Copied to Clipboard!';
					setTimeout(() => {
						document.querySelector('.invite-button').innerText = 'Copy Invite Link!';
					}, 1500);
				}}
				class="invite-button"
				>Copy Invite Link!
			</button>
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
								{hour + 7}:00
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
			<button class="save-avail" on:click={saveAllAvailabilities}>SAVE AVAILABILITY</button>
			{#if $successMsg}
				<p class="success-msg">{$successMsg}</p>
			{/if}
			<form on:submit|preventDefault={() => addTask(taskInput)}>
				<input
					type="text"
					bind:value={taskInput}
					placeholder="Enter task description (50 characters max)"
					maxlength="50"
				/>
				<button type="submit" disabled={!taskInput.trim()}>Add Task</button>
			</form>
			{#each $tasks as task (task.taskId)}
				<div class="task-item">
					<input
						type="checkbox"
						bind:checked={task.completed}
						on:click={() => toggleCompletion(task.taskId)}
						on:keypress={() => toggleCompletion(task.taskId)}
					/>
					<span class={task.completed ? 'completed-task' : ''}>{task.title}</span>
					{#if task.assignee}
						<span class={task.completed ? 'completed-task' : ''}>Assigned to: {task.assignee}</span>
					{/if}
					<button class="delete-task" on:click={() => deleteTaskWrapper(task.taskId)}>delete</button
					>
					<button class="self-assign" on:click={() => assignTaskToUser(task.taskId)}
						>Self Assign</button
					>
				</div>
			{/each}
			{#if $commonAvailability.isLoading}
				<div class="common-availability-message">CALCULATING COMMON AVAILABILITIES...</div>
			{:else if $commonAvailability.slots.length > 0}
				<div class="common-availability-message">
					<strong>Common Availabilities:</strong>
					<ul>
						{#each $commonAvailability.slots as slot}
							<li>{slot}</li>
						{/each}
					</ul>
				</div>
			{:else}
				<div class="common-availability-message">NO COMMON AVAILABILITIES FOUND.</div>
			{/if}
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
		margin-bottom: 0.15rem;
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
		max-width: 25rem;
		text-align: center;
		font-size: 1rem;
		background-color: #c9e7e7;
		border-radius: 10px;
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
		padding: 0.5rem;
		background-color: #f9f9f9;
		border-radius: 0.2rem;
		font-family: 'Baloo Bhai 2';
		align-items: center;
		font-size: 1.25rem;
		margin-top: 0;
		margin-bottom: 0.5rem;
	}

	.task-item input[type='checkbox'] {
		accent-color: #879db7;
		transform: scale(1.5);
		cursor: pointer;
		margin-top: -1rem;
		margin-left: -4rem;
		margin-right: -4rem;
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

	.self-assign {
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

	.self-assign:hover {
		background-color: gray;
		color: white;
	}

	.common-availability-message {
		text-align: left;
		font-family: 'Baloo Da 2';
		font-style: oblique;
		margin-top: 1rem;
		font-size: 1.25rem;
	}

	.save-avail {
		background-color: #87c7fa;
		color: black;
		border: none;
		font-weight: bold;
		cursor: pointer;
		margin-top: 0.25rem;
		padding: 0.5rem;
		display: inline-block;
		text-align: center;
		font-size: 1rem;
		border-radius: 0.3rem;
		margin-bottom: 1rem;
	}

	.save-avail:hover {
		background-color: gray;
		color: white;
	}

	.success-msg {
		margin-top: 0.5rem;
		margin-bottom: 0.5rem;
	}
</style>

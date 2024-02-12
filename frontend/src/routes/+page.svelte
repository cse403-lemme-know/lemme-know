
<script>
	import { goto } from '$app/navigation';
	import { Datepicker } from 'svelte-calendar';
	import dayjs from 'dayjs';
	import { writable } from 'svelte/store';
	import "$lib/model.js"
	import { startDate, endDate, groupName } from '$lib/stores';
	let name = 'LemmeKnow';
	let errorMsg = writable('');

	async function handleButtonClick() {
		if ($startDate && $endDate) {
			if (dayjs($startDate).isAfter($endDate)) {
				$errorMsg = 'Start date must be before the end date';
			} else {
				$errorMsg = '';
				if (!$groupName || $groupName.trim().length === 0) {
					$errorMsg = 'Please enter a group name';
				} else {
					startDate.set($startDate);
					endDate.set($endDate);
					groupName.set($groupName);
					try {
						const response = await fetch('https://localhost:8080/api/group/', {
							method: 'PATCH',
							headers: {
								'Content-Type': 'application/json',
							},
							credentials: 'include',
							body: JSON.stringify({
								name: $groupName,
							}),
						});

						if (response.ok) {
							const data = await response.json();
							console.log('Group created with ID:', data.groupId);
							goto('/dashboard');
						} else {
							$errorMsg = 'Failed to create group';
						}
					} catch (error) {
						console.error('Error creating group:', error);
						$errorMsg = 'Error connecting to the server';
					}
				}
			}
		}
	}
</script>

<header>
	<h1>{name}</h1>
	<nav>
		<a href="/features/">FEATURES</a>
		<a href="/privacy">PRIVACY</a>
		<a href="/devteam">DEV TEAM</a>
	</nav>
</header>

<main>
	<h3>SCHEDULE HANGOUTS, PLAN ROAD-TRIPS, SHARE CALENDARS, EVERYTHING, EVERYWHERE, ALL AT ONCE.</h3>
	<div class="images">
		<div class="container">
			<img src="/highfive.png" alt="3 people trying to figure out their collective availability" />
			<img src="/cal.png" alt="family of four sharing their calendars" />
		</div>
		<div class="date-picker-container">
			<label for="startDate">Start Date:</label>
			<Datepicker bind:selected={$startDate} />
			<label for="endDate">End Date:</label>
			<Datepicker bind:selected={$endDate} />
		</div>
		<div class="input-container">
			<input type="text" bind:value={$groupName} placeholder="Enter Group name.." />
			<button on:click={handleButtonClick}>Let me know!</button>
			<span class="errorMsg">{$errorMsg}</span>
		</div>
		<div class="container">
			<img src="/road.png" alt="group of friends driving in the car" />
			<img src="/calendar.png" alt="three colleagues figuring out their work schedule" />
		</div>
	</div>
</main>

<style>
	main {
		max-width: 1200px;
		margin: 0 auto;
		padding: 2rem;
		display: flex;
		flex-direction: column;
		align-items: center;
		text-align: center;
		font-family: 'Baloo Da 2', sans-serif;
	}

	header {
		width: 100%;
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 1rem 0;
	}

	.errorMsg {
		color: red;
		margin-top: 0.5rem;
	}

	nav a {
		text-decoration: none;
		color: #000;
		margin: 0 1rem;
		font-weight: bold;
		font-family: 'Baloo Da 2';
	}

	h1 {
		margin: 0.5rem 0;
		font-family: 'Baloo Bhai 2';
		padding: 0 2rem;
	}

	h3 {
		margin: 0.5rem 0;
		font-family: 'Baloo Da 2';
		font-weight: bolder;
		font-size: 2rem;
		color: #5d75cd;
	}

	img {
		width: 100%;
		max-width: 500px;
		margin: 0.5rem 0.5rem;
	}

	.images {
		display: flex;
		justify-content: space-around;
		align-items: center;
		flex-wrap: wrap;
		width: 100%;
		gap: 1rem;
	}

	.container {
		flex: 1;
		max-width: 50%;
		margin: 1rem;
	}

	.input-container {
		display: flex;
		flex: 1;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		margin: 1rem;
		border-color: white;
		color: white;
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

	input:focus {
		outline: none;
		border-color: black;
	}

	button {
		padding: 0.5rem 1rem;
		background-color: #879db7;
		color: white;
		border: none;
		font-size: 1rem;
		border-radius: 1rem;
		cursor: pointer;
	}

	button:hover {
		background-color: gray;
		color: white;
	}

	.date-picker-container {
		display: flex;
		flex-direction: column;
		align-items: center;
		z-index: 10;
		position: relative;
	}

	label {
		display: block;
		font-weight: bold;
		font-size: 1rem;
		color: black;
	}

	@media (min-width: 640px) {
		.images,
		.container,
		.input-container {
			max-width: 100%;
		}

		nav a {
			margin: 0.5rem;
		}
	}
</style>

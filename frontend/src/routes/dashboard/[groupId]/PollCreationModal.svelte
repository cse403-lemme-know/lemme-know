<script>
	//@ts-nocheck
	import { createPoll, updateVotes, deletePoll } from '$lib/model';

	export let groupId;
	export let group;
	let title = '';
	let options = [];
	let votes = [];

	function addOption() {
		options = [...options, ''];
	}

	function removeOption(index) {
		options = options.filter((_, i) => i !== index);
	}

	function handleCreatePoll() {
		createPoll(groupId, title, options);
	}
	// function addVote(vote) {
	// 	votes = [...votes, vote];
	// }
	function addVote(vote) {
		if (votes.includes(vote)) {
			votes = votes.filter((v) => v !== vote);
		} else {
			votes = [...votes, vote];
		}
		updateVotes(groupId, votes);
	}
	function handleUpdateVotes(vote) {
		addVote(vote);
		updateVotes(groupId, votes);
	}
	function handleDeletePoll() {
		deletePoll(groupId);
		votes = [];
	}

	$: totalVotes = group.poll
		? group.poll.options.reduce((acc, option) => acc + option.votes.length, 0)
		: 0;

	// function getTotalVotes() {
	// 	if (!group.poll || !group.poll.options) {
	// 		return 0;
	// 	}
	// 	return group.poll.options.reduce((total, opt) => total + opt.votes.length, 0);
	// }

	function getPercentage(numVotes) {
		return totalVotes === 0 ? 0 : ((numVotes || 0) / totalVotes) * 100;
	}
</script>

{#if !group.poll}
	<div class="modal">
		<div class="modal-content">
			<h2>Create Poll</h2>
			<label for="pollName">Poll Name:</label>
			<input type="text" id="pollName" bind:value={title} />

			<h3>Options:</h3>
			{#each options as _, index}
				<div class="option">
					<input type="text" bind:value={options[index]} />
					<button on:click={() => removeOption(index)}>Remove</button>
				</div>
			{/each}
			<button on:click={addOption}>Add Option</button>

			<button on:click={handleCreatePoll}>Create Poll</button>
		</div>
	</div>
{:else}
	<div class="poll">
		<h2>{group.poll.title}</h2>
		{#each group.poll.options as option}
			<div class="option">
				<span>{option.name}</span>
				<!-- <button on:click={() => handleUpdateVotes(option.option)}>Vote</button> -->
				<button on:click={() => handleUpdateVotes(option.name)}
					>{votes.includes(option.name) ? 'Unvote' : 'Vote'}</button
				>

				<span>({option.votes.length} votes)</span>
				<span>({getPercentage(option.votes.length).toFixed(2)}%)</span>
			</div>
		{/each}
		<p>Total Votes: {totalVotes}</p>
		<button on:click={() => handleDeletePoll()}>Dismiss Poll</button>
	</div>
{/if}

<style>
	.modal {
			background-color: #f0f0f0;
			padding: 20px;
			border-radius: 5px;
			box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
			margin-bottom: 20px;
	}

	.poll {
			background-color: #f0f0f0;
			padding: 10px;
			border-radius: 5px;
			box-shadow: 0 0 5px rgba(0, 0, 0, 0.1);
			margin-bottom: 20px;
	}

	.option {
			display: flex;
			align-items: center;
			justify-content: space-between;
			padding: 5px;
			border-radius: 5px;
			background-color: #ffffff;
			box-shadow: 0 0 3px rgba(0, 0, 0, 0.1);
	}

	.option input[type="text"] {
			flex-grow: 1;
			margin-right: 5px;
			font-size: 14px;
	}

	.option button {
			margin-left: 5px;
			padding: 5px;
			border: none;
			border-radius: 3px;
			background-color: #4CAF50;
			color: white;
			cursor: pointer;
			font-size: 12px;
	}

	.poll h2 {
			margin-bottom: 5px;
			font-size: 18px;
	}

	.poll p {
			margin-top: 10px;
			font-weight: bold;
			font-size: 14px;
	}

	.poll button {
		padding: 5px 10px;
		border: none;
		border-radius: 5px;
		background-color: #f44336;
		color: white;
		cursor: pointer;
		margin-top: 10px;
		font-size: 14px;
	}

</style>

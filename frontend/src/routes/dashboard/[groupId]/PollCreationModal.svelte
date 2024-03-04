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
		console.log(groupId);
		createPoll(groupId, title, options);
		console.log(group.Poll);
	}

	function handleUpdateVotes() {
		updateVotes(groupId, votes);
	}
	function handleDeletePoll() {
		deletePoll(groupId);
	}

	function getTotalVotes() {
		return Object.values(votes).reduce((acc, curr) => acc + curr, 0);
	}

	function getPercentage(optionIndex) {
		const totalVotes = getTotalVotes();
		return totalVotes === 0 ? 0 : ((votes[optionIndex] || 0) / totalVotes) * 100;
	}
</script>

{#if !group.Poll}
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
		<h2>{group.Poll.title}</h2>
		{#each group.Poll.Options as option, index}
			<div class="option">
				<span>{option}</span>
				<button on:click={() => handleUpdateVotes(index)}>Vote</button>
				<span>({votes[index] || 0} votes)</span>
				<span>({getPercentage(index).toFixed(2)}% of total votes)</span>
			</div>
		{/each}
		<p>Total Votes: {getTotalVotes()}</p>
	</div>
{/if}

<style>
	.modal {
		background-color: #f0f0f0;
	}

	.poll {
		background-color: #f0f0f0;
		padding: 20px;
		border-radius: 5px;
	}

	.option {
		margin-bottom: 10px;
	}
</style>

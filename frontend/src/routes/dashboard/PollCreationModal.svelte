<script>
	//@ts-nocheck
	import { createPoll, updateVotes, deletePoll } from '$lib/model';
	import { groupId } from '$lib/stores';


	let title = '';
	let options = [];
	let pollData = null;
	let votes = [];

	function addOption() {
		options = [...options, ''];
	}

	function removeOption(index) {
		options = options.filter((_, i) => i !== index);
	}

  function handleCreatePoll() {
		pollData = createPoll($groupId, title, options);
	}

	function handleUpdateVotes(optionIndex) {
		if (votes[optionIndex] === undefined) {
			votes[optionIndex] = 1;
		} else {
			votes[optionIndex]++;
		}
		// votes = updateVotes(votes);
	}

	function getTotalVotes() {
		return Object.values(votes).reduce((acc, curr) => acc + curr, 0);
	}

	function getPercentage(optionIndex) {
		const totalVotes = getTotalVotes();
		return totalVotes === 0 ? 0 : ((votes[optionIndex] || 0) / totalVotes) * 100;
	}

	console.log('createPoll');
</script>

{#if !pollData}
	<div class="modal">
		<div class="modal-content">
			<h2>Create Poll</h2>
			<label for="pollName">Poll Name:</label>
			<input type="text" id="pollName" bind:value={title} />

			<h3>Options:</h3>
			{#each options as { }, index}
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
		<h2>{pollData.name}</h2>
		{#each pollData.options as option, index}
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

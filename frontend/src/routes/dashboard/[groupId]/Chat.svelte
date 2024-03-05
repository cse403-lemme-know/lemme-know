<script>
	// @ts-nocheck

	import PollCreationModal from './PollCreationModal.svelte';
	import { sendMessage } from '$lib/model';

	export let groupId;
	export let group;
	let content = '';
	export let isPoll = false;
	/**
	 * Send a message to the chat.
	 */
	function handleSendMessage() {
		if (content.trim() !== '') {
			sendMessage(groupId, content.trim());
			content = '';
		}
	}
	/**
	 * send message if hit enter
	 * @param event
	 */
	function handleKeyPress(event) {
		if (event.key === 'Enter') {
			handleSendMessage();
		}
	}
</script>

<div class="chatbox">
	<div class="messages">
		{#if group}
			{#each group.messages as message (message.timestamp)}
				<div class:message class:message.sender={message.sender}>
					{#if true}
						<strong class="user-message">{message.sender}:</strong> {message.content}
					{/if}
				</div>
			{/each}
		{/if}
	</div>
	<div class="poll">
		{#if isPoll}
			<PollCreationModal {groupId} {group} />
		{/if}
	</div>
	<div class="input-bar">
		<input
			class="input"
			bind:value={content}
			placeholder="Type your message..."
			on:keydown={handleKeyPress}
		/>
		<button on:click={handleSendMessage} on:keyup={handleSendMessage}>Send Message</button>
	</div>
</div>

<style>
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
		flex-direction: column;
		justify-content: end;
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

	.input-bar {
		display: flex;
		align-items: center;
		justify-content: space-between;
		margin-top: 10px;
	}
</style>

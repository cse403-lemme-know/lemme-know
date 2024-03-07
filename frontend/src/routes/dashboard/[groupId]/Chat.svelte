<script>
	// @ts-nocheck

	import PollCreationModal from './PollCreationModal.svelte';
	import { sendMessage, users } from '$lib/model';

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
						{#if $users[message.sender] && $users[message.sender].name !== ''}
							<strong class="user-message">{$users[message.sender].name}:</strong> {message.content}
						{:else}
							<strong class="user-message">{message.sender}:</strong> {message.content}
						{/if}
					{/if}
				</div>
			{/each}
		{/if}
	</div>
	<div class="poll">
		{#if isPoll || (group && group.poll)}
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
		padding: 1rem 1rem 1rem 1rem;
		width: 500px;
		height: 650px;
		border-radius: 16px;
		box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
		overflow-y: auto; /* Add scrollbar when content exceeds the height */
		margin-right: 0.5rem;
		position: relative; /* Add relative positioning */
	}

	.messages {
		flex-grow: 1;
		display: flex;
		flex-direction: column;
		justify-content: flex-end;
	}

	.message {
		margin: 8px 0;
		padding: 8px;
		background-color: #f0f0f0;
		border-radius: 4px;
		word-wrap: break-word;
		font-family: 'Playfair Display', serif;
	}

	.user-message {
		background-color: #e6f7ff;
		text-align: right;
		font-family: 'Poppins', sans-serif;
	}

	.input-bar {
		/* display: flex; */
		align-items: center;
		justify-content: space-between;
		margin-top: 10px;
	}

	.input {
		width: 500px;
		margin-right: 1px;
	}

	button {
		margin-left: 1x;
	}
</style>

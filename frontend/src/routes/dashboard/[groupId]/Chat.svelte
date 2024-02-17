<script>
	import PollCreationModal from "./PollCreationModal.svelte";

    export let group;
	let messages = [];
	let content = '';
    export let isPoll = false;
	/**
	 * Send a message to the chat.
	 */
	function handleSendMessage() {
		if (content.trim() !== '') {
			messages = [...messages, { text: content, sender: 'user' }];
			content = content.trim();
			// sendMessage($groupId, content);
			content = '';
			// need to handle Fetch messages
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
        {#if isPoll}
            <PollCreationModal />
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
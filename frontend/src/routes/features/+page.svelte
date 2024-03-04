<script>
</script>

<main class="container">
	<h1>LemmeKnow</h1>
	<p>Description</p>
	<article>
		<p>
			LemmeKnow helps small groups, such as groups of students or friends, organize meetings or
			activities. Each member may express their time constraints, self-assign tasks, reach a
			consensus by creating a poll, and communicate via chat. As an activity approaches, members
			receive reminders and can check the confirmation status of other members.
		</p>
	</article>

	<section class="access">
		<h2>Access</h2>
		<p>
			Navigate to <a href="https://lemmeknow.xyz" target="_blank" rel="noopener noreferrer"
				>https://lemmeknow.xyz</a
			> in your web browser.
		</p>
	</section>

	<section class="features">
		<h2>Features</h2>
		<ul>
			<li>Interactive web-based dashboard</li>
			<li>
				Calendar, shared between members, showing group activities and individual availabilities
			</li>
			<li>Polls, for reaching a consensus based on individual preferences</li>
			<li>Tasks, for assigning preparation steps in advance</li>
			<li>Messaging system, for general discussion and coordination</li>
			<li>Reminders, to increase the chance of attendance</li>
		</ul>
	</section>

	<section class="usage">
		<h2>Usage</h2>
		<div class="usage-section">
			<h3>Creating a group</h3>
			<p>
				To create a group, type the name of the group on the homepage. Use the start date picker to
				identify the earliest possible meeting date and the end date picker to identify the latest
				possible meeting date. Press the "Let me know!" button to finish creating the group.
			</p>
		</div>

		<div class="usage-section">
			<h3>Inviting group members</h3>
			<p>
				Press the "Copy Invite Link!" button on the sidebar to copy a group invitation link to your
				clipboard. You may share this link with other members to enable them to join the group.
			</p>
		</div>

		<div class="usage-section">
			<h3>Viewing Group Member's activity status</h3>
			<p>
				You can view your group members activity status by clicking on the "Members" Icon in the
				dashboard page. You can also update your activity status by clicking the option you want in
				the drop down menu
			</p>
		</div>

		<div class="usage-section">
			<h3>Joining a group</h3>
			<p>
				Navigate to an invite link in your web browser. You will automatically join the
				corresponding group.
			</p>
		</div>

		<div class="usage-section">
			<h3>Chat</h3>
			<p>This section is TBD.</p>
		</div>

		<div class="usage-section">
			<h3>Polls</h3>
			<p>This section is TBD.</p>
		</div>

		<div class="usage-section">
			<h3>Calendar</h3>
			<p>
				Select your availability based on your schedule. If another member of the group also selects
				the same availability, the intersecting availability will appear under "Common
				Availabilites"
			</p>
		</div>

		<div class="usage-section">
			<h3>Tasks</h3>
			<p>Create Tasks for your group, members can self assign tasks to themselves</p>
		</div>
	</section>

	<section class="bug-reports">
		<h2>Bug reports</h2>
		<p>
			Please use our <a
				href="https://github.com/cse403-lemmeknow/lemmeknow/issues"
				target="_blank"
				rel="noopener noreferrer">issue tracker</a
			> to report bugs and propose features. Include details and context about the issue, using the corresponding
			template.
		</p>
	</section>

	<section class="known-bugs">
		<h2>Known bugs</h2>
		<p>
			Please check our <a
				href="https://github.com/cse403-lemmeknow/lemmeknow/issues"
				target="_blank"
				rel="noopener noreferrer">issue tracker</a
			> to see known bugs and limitations.
		</p>
	</section>
</main>

<style>
	.container {
		max-width: 800px;
		margin: auto;
		padding: 20px;
	}

	h1,
	h2 {
		color: #333;
	}

	article,
	section {
		margin-bottom: 20px;
	}

	p,
	a,
	ul,
	li {
		font-size: 16px;
		color: #666;
	}

	a {
		color: #007bff;
	}

	a:hover {
		text-decoration: underline;
	}

	.usage-section {
		margin-bottom: 15px;
	}

	@media (max-width: 768px) {
		.container {
			padding: 10px;
		}

		h1,
		h2 {
			font-size: 1.5rem;
		}

		p,
		a,
		ul,
		li {
			font-size: 14px;
		}
	}
</style>

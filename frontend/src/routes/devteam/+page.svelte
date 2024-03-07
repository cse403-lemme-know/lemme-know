<script>
</script>

<main>
	<h1>CSE 403 Team Project: LemmeKnow</h1>
	<h2>Backend Team</h2>
	<div class="team-section">
		<div class="team-member">
			<span class="role">Finn Bear:</span> Backend Developer
			<img class="team-member-image" src="../finn.png" alt="Finn Bear" />
		</div>
		<div class="team-member">
			<span class="role">Daven Tejalaksana:</span> Backend Developer
			<img class="team-member-image" src="daven.png" alt="Daven Tejalaksana" />
		</div>
		<div class="team-member">
			<span class="role">Mohamad Abdi:</span> Backend Developer
			<img class="team-member-image" src="../mohamed.png" alt="Mohamad Abdi" />
		</div>
	</div>

	<h2>Frontend Team</h2>
	<div class="team-section">
		<div class="team-member">
			<span class="role">Brooklyn Zhu:</span> Frontend Developer
			<img class="team-member-image" src="../brooklyn.png" alt="Brooklyn Zhu" />
		</div>
		<div class="team-member">
			<span class="role">Zohar Le:</span> Frontend Developer
			<img class="team-member-image" src="../zohar.png" alt="Zohar Le" />
		</div>
	</div>
</main>

<style>
	main {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		min-height: 100vh;
		text-align: center;
	}

	h1,
	h2 {
		color: #333;
	}

	.team-section {
		display: flex;
		flex-direction: row;
		justify-content: center;
		margin-bottom: 20px;
	}

	.team-member {
		font-size: 16px;
		color: #666;
		margin: 10px 10px;
	}

	.role {
		font-weight: bold;
		color: #007bff;
	}

	.team-member-image {
		width: 300px;
		height: 200px;
		display: block;
		margin: 10px auto;
	}
</style>

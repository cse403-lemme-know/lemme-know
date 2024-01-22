<script>
	export let name;
</script>

<main>
	<h1>{name}</h1>
	<p>hi</p>
</main>

<style>
	main {
		text-align: left;
		/*padding: 0em;*/
		max-width: 230px;
		margin: 1em 2em;
	}

	h1 {
		color: #817f7f;
		/*text-transform: uppercase;*/
		font-size: 4em;
		font-family: Helvetica;
		/*font-weight: 100;*/
	}

	@media (min-width: 640px) {
		main {
			max-width: none;
		}
	}
</style>
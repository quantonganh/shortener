<script>
	let longURL
	let shortURL

	async function doShorten() {
        const res = await fetch('/create-short-url', {
            method: 'POST',
	        body: JSON.stringify({
		        'long_url': longURL
	        })
        })

		const json = await res.json()
		shortURL = json.short_url
	}
</script>

<form on:submit|preventDefault={doShorten}>
	<input required placeholder="Long URL" type="url" bind:value={longURL}>
	<button>Shorten</button>
</form>

<a href="{shortURL}">{shortURL}</a>
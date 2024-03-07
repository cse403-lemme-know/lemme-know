async function getVapidPublicKey() {
	const response = await fetch(`//${location.host}/api/push/`);
	const json = await response.json();
	return json.vapidPublicKey;
}

async function subscribe() {
	if (!navigator.serviceWorker) {
		return;
	}
	const vapidPublicKey = await getVapidPublicKey();
	navigator.serviceWorker
		.register('/serviceWorker.js', { scope: '/' })
		.then(function (registration) {
			return registration.pushManager.getSubscription().then(function (subscription) {
				if (subscription) {
					return subscription;
				}
				return registration.pushManager.subscribe({
					// required.
					userVisibleOnly: true,
					applicationServerKey: vapidPublicKey
				});
			});
		})
		.then(function (subscription) {
			return fetch(`//${location.host}/api/push/`, {
				method: 'PATCH',
				body: JSON.stringify(subscription)
			});
		});
}

if (window.Notification && Notification.permission === 'granted') {
	subscribe();
} else if (window.Notification && Notification.permission !== 'denied') {
	Notification.requestPermission(function (permission) {
		if (permission === 'granted') {
			subscribe();
		}
	});
}

self.addEventListener('push', function (event) {
	var data = event.data.json();

    if (data.message) {
        // https://developer.mozilla.org/en-US/docs/Web/API/ServiceWorkerRegistration/showNotification
        event.waitUntil(
            self.registration.showNotification(`${data.message.sender} in ${data.message.group}`, {
                body: data.message.content,
                timestamp: data.message.timestamp,
            })
        );
    }
});

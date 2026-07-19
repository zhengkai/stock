self.addEventListener("push", (event) => {
	const data = event.data.json();

	console.log('event.data', event.data);
	console.log('event.data.json', data);

	event.waitUntil(
		self.registration.showNotification(
			data.title,
			{
				body: data.body,
				icon: "/icon.png"
			}
		)
	);
});

console.log("Service Worker: Push notifications are enabled.");

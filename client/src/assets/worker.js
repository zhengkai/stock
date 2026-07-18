self.addEventListener("push", event => {
	const data = event.data.json();

	self.registration.showNotification(
		data.title,
		{
			body: data.body,
			icon: "/icon.png"
		}
	);
});

console.log("Service Worker: Push notifications are enabled.");

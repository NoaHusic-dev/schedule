export class ClockState {
	time = $state(new Date());

	constructor() {
		setInterval(() => {
			this.time = new Date();
		}, 1000);
	}

	get formatted() {
		return this.time.toLocaleTimeString('en-US', {
			hour12: false,
			hour: '2-digit',
			minute: '2-digit',
			second: '2-digit'
		});
	}
}

export const clock = new ClockState();
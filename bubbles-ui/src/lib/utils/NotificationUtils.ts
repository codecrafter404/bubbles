export function confirm(msg: string): boolean {
	let message = `${msg}\nPress 0 to confirm\nPress any key to cancel`;
	let x = prompt(message);
	return x == "0";
}

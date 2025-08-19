export enum KeyboardPressedEvent {
	MoveUp,
	MoveDown,
	MoveLeft,
	MoveRight,
	Confirm,
	Add,
	Remove,
}
const config: Map<KeyboardPressedEvent, Array<PressedKey>> = new Map();

const non = { control: false, shift: false, alt: false }

config.set(KeyboardPressedEvent.Confirm, [{ code: "Enter", ...non }])
config.set(KeyboardPressedEvent.MoveUp, [{ code: "ArrowUp", ...non }, { code: "KeyK", ...non }])
config.set(KeyboardPressedEvent.MoveDown, [{ code: "ArrowDown", ...non }, { code: "KeyJ", ...non }])
config.set(KeyboardPressedEvent.MoveLeft, [{ code: "ArrowLeft", ...non }, { code: "KeyH", ...non }])
config.set(KeyboardPressedEvent.MoveRight, [{ code: "ArrowRight", ...non }, { code: "KeyL", ...non }])
config.set(KeyboardPressedEvent.Add, [{ code: "BracketRight", ...non }, { code: "Space", ...non }])
config.set(KeyboardPressedEvent.Remove, [{ code: "Slash", ...non }, { code: "Backspace", ...non }])

export interface ShortcutConfig {
	MoveUp: Array<PressedKey>
	MoveDown: Array<PressedKey>
	MoveLeft: Array<PressedKey>
	MoveRight: Array<PressedKey>
	Confirm: Array<PressedKey>
}

export interface PressedKey {
	code: string
	control: boolean
	shift: boolean
	alt: boolean
}

export function ProcessEvent(event: KeyboardEvent): KeyboardPressedEvent | undefined {
	let ctl = event.metaKey || event.ctrlKey;
	let alt = event.altKey;
	let shift = event.shiftKey;

	let res: KeyboardPressedEvent | undefined;
	config.forEach((v, k) => {
		if (res != undefined) {
			return
		}
		let found = v.find((x) => x.alt == alt && x.control == ctl && x.shift == shift && x.code == event.code)
		if (found != undefined) {
			res = k
		}
	})
	return res
}

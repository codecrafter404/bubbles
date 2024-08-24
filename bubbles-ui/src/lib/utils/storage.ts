export function load_endpoint(): string {
	let res = localStorage.getItem("gql_endpoint");
	if (res != null) {
		return res;
	}

	return "http://localhost:8080/query";
}
export function load_ws_endpoint(): string {
	let res = localStorage.getItem("gql_ws_endpoint");
	if (res != null) {
		return res;
	}

	return "ws://localhost:8080/query";
}

export function update_endpoint(endpoint: string) {
	localStorage.setItem("gql_endpoint", endpoint)
	location.reload()
}
export function update_ws_endpoint(endpoint: string) {
	localStorage.setItem("gql_ws_endpoint", endpoint)
	location.reload()
}

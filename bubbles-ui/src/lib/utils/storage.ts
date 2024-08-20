export function load_endpoint(): string {
	let res = localStorage.getItem("gql_endpoint");
	if (res != null) {
		return res;
	}

	return "http://localhost:8080/query";
}

export function update_endpoint(endpoint: string) {
	localStorage.setItem("gql_endpoint", endpoint)
	location.reload()
}

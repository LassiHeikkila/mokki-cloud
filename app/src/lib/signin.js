import React from "react";

export async function signin(username, password) {
	fetch(
		`/api/authorize`, {
			method: "POST",
			mode: "cors",
			body: JSON.stringify({
				"username": username,
				"password": password
			})
		})
	.then(response => {
		return response.json();
	})
	.catch( (error) => {
		console.error("error logging in")
		return {};
	})
}

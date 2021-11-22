import React, { useState } from "react";
import Form from "react-bootstrap/Form";
import Button from "react-bootstrap/Button";
import { useHistory } from "react-router-dom";
import { useAppContext } from "../lib/contextLib";
import "./Login.css";

export default function Login() {
	const [password, setPassword] = useState("");

	const history = useHistory();
	var context = useAppContext();

	function validateForm() {
		return password.length > 0;
	}

	async function signIn(pw) {
		fetch(
			`/api/authorize`, {
				method: "POST",
				mode: "cors",
				body: JSON.stringify({
					"username": "generic-user",
					"password": pw
				})
			})
		.then(response => {
			return response.json();
		})
		.then(data => {
			if (data["ok"] !== true) {
				console.error("failed to log in!");
				throw new Error('login failed');
			}

			console.log("log in succeeded!");
			console.log("got token: ", data["token"]);
			context.isAuthenticated = true;
			context.storedToken = data["token"];
			console.log("saving token to local storage")
			localStorage.setItem('authToken', data["token"]);
			history.push("/");

		})
		.catch( (err) => {
			alert(`Failed to login with error: ` + err.message);
			context.isAuthenticated = false;
			context.storedToken = "";
		})
	}

	async function handleSubmit(event) {
		event.preventDefault();

		signIn(password)
	}

	return (
		<div className="Login">
			<Form onSubmit={handleSubmit} >
				<Form.Group size="lg" controlId="password">
					<Form.Label>Password</Form.Label>
					<Form.Control
						autoFocus
						type="password"
						value={password}
						onChange={(e) => setPassword(e.target.value)}
					/>
				</Form.Group>
				<Button block size="lg" type="submit" disabled={!validateForm()}>
				Login
				</Button>
			</Form>
		</div>
	);
}

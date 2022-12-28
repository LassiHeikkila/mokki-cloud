import React, { useState } from "react";
import { useDispatch } from 'react-redux';
import Form from "react-bootstrap/Form";
import Button from "react-bootstrap/Button";
import "./Login.css";
import { doAuthorize } from "../lib/api";
import config from '../config.json';
import { setToken, setIsAuthenticated } from "../state/AppState";

const Login = (props) => {
	const [password, setPassword] = useState("");

	const dispatch = useDispatch();

	function validateForm() {
		return password.length > config.minimumPasswordLength;
	}

	async function signIn(pw) {
		doAuthorize(pw)
		.then(token => {
			if (token === '') {
				dispatch(setIsAuthenticated(false));
				dispatch(setToken(''));
			} else {
				console.info('login succeeded');
				dispatch(setIsAuthenticated(true));
				dispatch(setToken(token));

				if (props.onLogin) {
					props.onLogin();
				}
			}
		})
	}

	async function handleSubmit(event) {
		event.preventDefault();

		signIn(password);
	}

	return (
		<Form onSubmit={handleSubmit} >
			<Form.Group size='lg' controlId='password' className='mb-3'>
				<Form.Label>Enter password:</Form.Label>
				<Form.Control
					autoFocus
					type='password'
					placeholder='password'
					value={password}
					onChange={(e) => setPassword(e.target.value)}
				/>
			</Form.Group>
			<Button size='lg' type='submit' disabled={!validateForm()}>Login</Button>
		</Form>
	);
};

export default Login;

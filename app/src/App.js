import React, { useState, useEffect } from 'react';
import './App.css';
import Routes from "./Routes";
import Navbar from "react-bootstrap/Navbar";
import Nav from "react-bootstrap/Nav";
import { LinkContainer } from "react-router-bootstrap";
import { AppContext } from "./lib/contextLib";

function App() {
	const [isAuthenticating, setIsAuthenticating] = useState(true);
	const [isAuthenticated, setIsAuthenticated] = useState(false);
	const [storedToken, setStoredToken] = useState(localStorage.getItem('authToken') || '');

	useEffect(() => {
		async function onLoad() {
			console.log("onLoad called, doing initialization stuff");
			if (storedToken === '') {
				console.log("no token in local storage")
				setIsAuthenticated(false);
				setIsAuthenticating(false);
				return;
			}
			await fetch(
				`/api/checkToken`, {
					method: "GET",
					mode: "cors",
					headers: {
						'X-API-KEY': storedToken
					}
				}
			).then((response) => {
				if (response.status === 200) {
					console.log("stored token is valid")
					setIsAuthenticated(true);
					setStoredToken(localStorage.getItem('authToken'));
				} else {
					console.error("failed to verify stored token as valid")
					setIsAuthenticated(false);
					localStorage.setItem('authToken', '');
				}
			}).catch((e) => {
				setIsAuthenticated(false);
			})
			setIsAuthenticating(false);
		}
		onLoad();
	}, [isAuthenticated, storedToken]);

	return (
		!isAuthenticating && (
			<div className="App container py-3">
			<Navbar collapseOnSelect bg="light" expand="md" className="mb-3">
				<LinkContainer to="/">
					<Navbar.Brand href="/" className="font-weight-bold text-muted">
						Mokki status
					</Navbar.Brand>
				</LinkContainer>
				<Navbar.Toggle />
				<Navbar.Collapse className="justify-content-end">
					<Nav activeKey={window.location.pathname}>
					{isAuthenticated ? (
						<LinkContainer to="/logout">
							<Nav.Link>Logout</Nav.Link>
						</LinkContainer>
					) : (
						<LinkContainer to="/login">
							<Nav.Link>Login</Nav.Link>
						</LinkContainer>
					)}
					</Nav>
				</Navbar.Collapse>
			</Navbar>
			<AppContext.Provider value={{ isAuthenticated, setIsAuthenticated, storedToken }}>
				<Routes />
			</AppContext.Provider>
			</div>
		)
	);
}

export default App;

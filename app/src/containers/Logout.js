import React, { useEffect }  from "react";
import { useHistory } from "react-router-dom";
import { useAppContext } from "../lib/contextLib";

export default function Logout() {
	const { setIsAuthenticated, storedToken } = useAppContext();
	const history = useHistory();

	useEffect(() => {
		function logout() {
			setIsAuthenticated(false);
			history.push("/login");
			localStorage.setItem('authToken', '');
		}
		logout();
	}, [history, setIsAuthenticated]);

	return (
		<div>
			<h2>
				Logged out!
			</h2>
		</div>
	);
}

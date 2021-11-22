import React from "react";
import  { Route, Switch } from "react-router-dom";
import Home from "./containers/Home";
import Login from "./containers/Login";
import Logout from "./containers/Logout";
import NotFound from "./containers/NotFound";

export default function Routes() {
	return (
		<Switch>
			<Route exact path="/">
				<Home />
			</Route>
			<Route exact path="/login">
				<Login />
			</Route>
			<Route exact path="/logout">
				<Logout />
			</Route>
			<Route>
				<NotFound />
			</Route>
		</Switch>
	);
}

import React from 'react';
import { useDispatch } from 'react-redux';
import Button from 'react-bootstrap/Button';

import { storeState } from '../state/Persist';
import { setToken, setIsAuthenticated } from '../state/AppState';

const Logout = (props) => {
	const dispatch = useDispatch();

	function logout() {
		dispatch(setIsAuthenticated(false));
		dispatch(setToken(''));
		storeState();

		if (props.onLogout) {
			props.onLogout();
		}
	}

	return (
		<Button block size='lg' type='button' onClick={() => { logout(); }}>Log out</Button>
	);
};

export default Logout;

import React, { useState, useEffect } from 'react';
import { useSelector, useDispatch } from 'react-redux';

import './App.css';
import Container from 'react-bootstrap/Container';
import Navbar from 'react-bootstrap/Navbar';
import Button from 'react-bootstrap/Button';
import Modal from 'react-bootstrap/Modal';

import { doCheckToken } from './lib/api';

import {
	selectToken,
	setToken,
	setIsAuthenticated,
	selectIsAuthenticated,
	selectDarkmode,
	selectSelectedMeasurement,
	selectSelectedSensor,
	selectSelectedTimePeriod,
	setDarkmode,
	setSelectedMeasurement,
	setSelectedSensor,
	setSelectedTimePeriod,
} from './state/AppState';

import { loadState, storeState } from './state/Persist';

import Home from './containers/Home';
import Login from './containers/Login';
import Logout from './containers/Logout';
import NotLoggedIn from './containers/NotLoggedIn';

function App() {
	const dispatch = useDispatch();

	const isAuthenticated = useSelector(selectIsAuthenticated);
	const token = useSelector(selectToken);
	const darkmode = useSelector(selectDarkmode);
	const selectedMeasurement = useSelector(selectSelectedMeasurement);
	const selectedSensor = useSelector(selectSelectedSensor);
	const selectedTimePeriod = useSelector(selectSelectedTimePeriod);

	const [storedStateLoaded, setStoredStateLoaded] = useState(false);
	const [initializationCompleted, setInitializationCompleted] = useState(false);
	const [loginCardVisible, setLoginCardVisible] = useState(false);
	const [logoutCardVisible, setLogoutCardVisible] = useState(false);

	const handleShowLoginCard = () => setLoginCardVisible(true);
	const handleShowLogoutCard = () => setLogoutCardVisible(true);
	const handleHideLoginCard = () => setLoginCardVisible(false);
	const handleHideLogoutCard = () => setLogoutCardVisible(false);

	const loadStoredState = () => {
		console.info('loading state');
		const state = loadState();

		if (!state) {
			console.warn('no state loaded');
			setStoredStateLoaded(true);
			return;
		}

		if (state.token) {
			dispatch(setToken(state.token));
		}

		if (state.darkmode) {
			dispatch(setDarkmode(state.darkmode));
		}

		if (state.selectedMeasurement) {
			dispatch(setSelectedMeasurement(state.selectedMeasurement));
		}

		if (state.selectedSensor) {
			dispatch(setSelectedSensor(state.selectedSensor));
		}

		if (state.selectedTimePeriod) {
			dispatch(setSelectedTimePeriod(state.selectedTimePeriod));
		}

		setStoredStateLoaded(true);
	};

	useEffect(() => {
		if (!storedStateLoaded) {
			// don't persist blank state before stored state is loaded first
			return;
		}
		console.info('persisting state');
		storeState({
			token: token,
			darkmode: darkmode,
			selectedMeasurement: selectedMeasurement,
			selectedSensor: selectedSensor,
			selectedTimePeriod: selectedTimePeriod,
		});
	}, [storedStateLoaded, token, darkmode, selectedMeasurement, selectedSensor, selectedTimePeriod]);

	// load stored state from local storage on initial load
	useEffect(() => {
		loadStoredState();
	}, []);

	useEffect(() => {
		if (initializationCompleted) {
			return;
		}
		if (!token || token === '') {
			return;
		}
		console.debug(`validating token ${token}`);

		doCheckToken(
			token
		).then((tokenOK) => {
			if (tokenOK) {
				console.log('stored token is valid');
				dispatch(setIsAuthenticated(true));
			} else {
				console.error('failed to verify stored token as valid')
				dispatch(setIsAuthenticated(false));
			}
			setInitializationCompleted(true);
		});
	}, [token, initializationCompleted]);

	return (
		<Container fluid='lg'>
			<Navbar bg='light' expand='lg'>
				<Navbar.Brand>Mökki status</Navbar.Brand>
				<Navbar.Toggle aria-controls="basic-navbar-nav" />
				<Navbar.Collapse className="justify-content-end">
					{isAuthenticated ?
						(<Button size='lg' onClick={handleShowLogoutCard}>Log out</Button>)
						:
						(<Button size='lg' onClick={handleShowLoginCard}>Log in</Button>)
					}
				</Navbar.Collapse>
			</Navbar>
			<Modal show={loginCardVisible} onHide={handleHideLoginCard}>
				<Modal.Header closeButton>Log in:</Modal.Header>
				<Modal.Body>
					<Login onLogin={handleHideLoginCard} />
				</Modal.Body>
			</Modal>
			<Modal show={logoutCardVisible} onHide={handleHideLogoutCard}>
				<Modal.Header closeButton>Are you sure you want to log out?</Modal.Header>
				<Modal.Body>
					<Logout onLogout={handleHideLogoutCard} />
				</Modal.Body>
			</Modal>
			{isAuthenticated ?
				<Home />
				:
				<NotLoggedIn />
			}
		</Container>
	);
}

export default App;

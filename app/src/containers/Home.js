import React, { useState, useEffect } from 'react';
import { useSelector, useDispatch } from 'react-redux';

import Container from 'react-bootstrap/Container';
import Stack from 'react-bootstrap/Stack';
import Dropdown from 'react-bootstrap/Dropdown';
import Table from 'react-bootstrap/Table';

import MyChart from "./Chart";
import './Home.css';

import {
	selectToken,
	selectIsAuthenticated,
	selectSelectedMeasurement,
	selectSelectedSensor,
	setSelectedMeasurement,
	setSelectedSensor
} from '../state/AppState';

import { doApiCall } from '../lib/api';
import config from '../config.json';

function unitFromQuery(query) {
	switch (query) {
		case "temperature":
			return "°C"
		case "pressure":
			return "Pa"
		case "humidity":
			return "%"
		default:
			return "unknown"
	}
}

const getSensorName = (settings, id) => {
	// range over sensors, return name if value === id
	for (var s of settings.sensors) {
		if (s.value === id) {
			return s.name;
		}
	}
	return '';
};

const Home = () => {
	const dispatch = useDispatch();

	var now = new Date();
	var rangeStart = new Date();
	rangeStart.setHours(rangeStart.getHours() - (24 * 7));
	const interval = 30 * 60;

	const [activeQuery, setActiveQuery] = useState(useSelector(selectSelectedMeasurement));
	const [measurementData, setMeasurementData] = useState({}); // map with sensor id as key, latest measurement as value
	const [activeSensor, setActiveSensor] = useState(useSelector(selectSelectedSensor));
	const [activeUnit, setActiveUnit] = useState("°C");
	const [rangeData, setRangeData] = useState([]);
	const [startTime] = useState(rangeStart.toISOString());
	const [stopTime] = useState(now.toISOString());
	const [gotSettings, setGotSettings] = useState(false);
	const [suggestedMin, setSuggestedMin] = useState(null);
	const [suggestedMax, setSuggestedMax] = useState(null);

	const token = useSelector(selectToken);
	const isAuthenticated = useSelector(selectIsAuthenticated);

	const [settings, setSettings] = useState({
		"suggestedMins": {
			"temperature": -10,
			"pressure": 95000,
			"humidity": 0
		},
		"suggestedMaxs": {
			"temperature": 30,
			"pressure": 105000,
			"humidity": 100
		},
		"sensors": []
	});

	useEffect(() => {
		// fetch settings.json on initial load
		fetch(
			`${config.backendUrl}/static/settings.json`
		).then(response => response.json()
		).then(data => {
			if (data.sensors != null && data.sensors.length > 0) {
				setSettings(data);
				setGotSettings(true);

				if (activeSensor === '') {
					setActiveSensor(data.sensors[0].value);
				}
			} else {
				throw new Error("settings don't contain any sensors")
			}
		}).catch((error) => {
			console.log("error fetching settings: ", error)
		});
	}, []);

	useEffect(() => {
		// fetch latest data point for each sensor
		if (!gotSettings || !isAuthenticated || !activeQuery) {
			return;
		}
		for (var i = 0; i < settings.sensors.length; i++) {
			const sensor = settings.sensors[i];
			doApiCall(
				token,
				'GET',
				`data/${activeQuery}/${sensor.value}/latest`
			).then(data => {
				var latestMeasurements = measurementData;
				latestMeasurements[sensor.value] = data[`${activeQuery}`];
				setMeasurementData(latestMeasurements);
			}).catch((err) => {
				console.error(`error fetching latest data for ${sensor.value}:`, err);
			})
		}
	}, [gotSettings, activeQuery, token, isAuthenticated]);

	useEffect(() => {
		// fetch range data
		if (!isAuthenticated || !activeQuery || !activeSensor) {
			return;
		}
		doApiCall(token, 'GET', `data/${activeQuery}/${activeSensor}/range?from=${startTime}&to=${stopTime}&interval=${interval}`)
			.then(data => {
				setRangeData(data);
			}).catch((error) => {
				console.log("error getting range data: ", error);
				setRangeData([]);
			});
	}, [gotSettings, activeQuery, activeSensor, startTime, stopTime, interval, token, isAuthenticated]);

	useEffect(() => {
		setActiveUnit(unitFromQuery(activeQuery));
		console.info(`dispatching selected measurement: ${activeQuery}`);
		dispatch(setSelectedMeasurement(activeQuery));
	}, [activeQuery, dispatch]);

	useEffect(() => {
		console.info(`dispatching selected sensor: ${activeSensor}`);
		dispatch(setSelectedSensor(activeSensor));
	}, [activeSensor, dispatch]);

	useEffect(() => {
		setSuggestedMin(settings.suggestedMins[`${activeQuery}`]);
		setSuggestedMax(settings.suggestedMaxs[`${activeQuery}`]);
	}, [settings, activeQuery]);

	return (
		<Container fluid='true'>
			<h3>Latest readings</h3>
			<Table>
				<thead>
					{settings.sensors.map((sensor) => (
						<th>{sensor.name}</th>
					))}
				</thead>
				<tbody>
					<tr>
						{settings.sensors.map((sensor) => (
							<td>{measurementData[sensor.value]} {activeUnit}</td>
						))}
					</tr>
				</tbody>
			</Table>
			<h3>Historical data</h3>
			<Stack direction='horizontal' gap={3}>
				<Dropdown>
					<Dropdown.Toggle variant='outline-primary'>Measurement: {activeQuery}</Dropdown.Toggle>
					<Dropdown.Menu>
						<Dropdown.Item onClick={() => { setActiveQuery('temperature') }}>temperature</Dropdown.Item>
						<Dropdown.Item onClick={() => { setActiveQuery('humidity') }}>humidity</Dropdown.Item>
						<Dropdown.Item onClick={() => { setActiveQuery('pressure') }}>air pressure</Dropdown.Item>
					</Dropdown.Menu>
				</Dropdown>
				<Dropdown>
					<Dropdown.Toggle variant='outline-primary'>Sensor: {getSensorName(settings, activeSensor)}</Dropdown.Toggle>
					<Dropdown.Menu>
						{settings.sensors.map((sensor) => (
							<Dropdown.Item onClick={() => { setActiveSensor(sensor.value) }}>{sensor.name}</Dropdown.Item>
						))}
					</Dropdown.Menu>
				</Dropdown>
			</Stack>

			<MyChart
				data={rangeData}
				measurement={activeQuery}
				period='week'
				suggestedMin={suggestedMin}
				suggestedMax={suggestedMax}
				unit={activeUnit}
			/>

		</Container>
	);
}

export default Home;

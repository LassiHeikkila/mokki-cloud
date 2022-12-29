import React, { useState, useEffect } from 'react';
import { useSelector, useDispatch } from 'react-redux';

import Container from 'react-bootstrap/Container';
import Dropdown from 'react-bootstrap/Dropdown';
import DropdownButton from 'react-bootstrap/DropdownButton';
import Table from 'react-bootstrap/Table';
import Navbar from 'react-bootstrap/Navbar';

import MyChart from "./Chart";
import './Home.css';

import {
	selectToken,
	selectIsAuthenticated,
	selectSelectedMeasurement,
	selectSelectedSensor,
	selectSelectedTimePeriod,
	setSelectedMeasurement,
	setSelectedSensor,
	setSelectedTimePeriod,
} from '../state/AppState';

import { doApiCall } from '../lib/api';
import { unitFromQuery, getPeriodString, getInterval } from '../lib/unitConversions';
import { calculateRangeStart } from '../lib/time';
import { oneDay, threeDays, oneWeek, oneMonth } from '../lib/units';

import config from '../config.json';

const getSensorName = (settings, id) => {
	// range over sensors, return name if value === id
	for (var s of settings.sensors) {
		if (s.value === id) {
			return s.name;
		}
	}
	return 'N/A';
};


const Home = () => {
	const dispatch = useDispatch();

	const [activeQuery, setActiveQuery] = useState(useSelector(selectSelectedMeasurement));
	const [latestTemperatureData, setLatestTemperatureData] = useState({}); // map with sensor id as key, latest measurement as value
	const [latestHumidityData, setLatestHumidityData] = useState({}); // map with sensor id as key, latest measurement as value
	const [latestPressureData, setLatestPressureData] = useState({}); // map with sensor id as key, latest measurement as value
	const [activeSensor, setActiveSensor] = useState(useSelector(selectSelectedSensor));
	const [activeUnit, setActiveUnit] = useState("°C");
	const [rangeData, setRangeData] = useState([]);
	const [startTime, setStartTime] = useState(calculateRangeStart(oneWeek));
	const [stopTime] = useState(new Date());
	const [timePeriod, setTimePeriod] = useState(useSelector(selectSelectedTimePeriod));
	const [gotSettings, setGotSettings] = useState(false);
	const [suggestedMin, setSuggestedMin] = useState(null);
	const [suggestedMax, setSuggestedMax] = useState(null);
	const [sensorDropdownButtonTitle, setSensorDropdownButtonTitle] = useState('Sensor');
	const [measurementDropdownButtonTitle, setMeasurementDropdownButtonTitle] = useState('Measurement');
	const [timePeriodButtonTitle, setTimePeriodButtonTitle] = useState('Period');

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

				if (activeQuery === '') {
					setActiveQuery('temperature');
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
		if (!gotSettings) {
			return;
		}
		for (var i = 0; i < settings.sensors.length; i++) {
			const sensor = settings.sensors[i];
			doApiCall(
				token,
				'GET',
				`data/temperature/${sensor.value}/latest`
			).then(data => {
				var d = latestTemperatureData;
				d[sensor.value] = data['temperature'];
				setLatestTemperatureData(d);
			}).catch((err) => {
				console.error(`error fetching latest temperature data for ${sensor.value}:`, err);
			})

			doApiCall(
				token,
				'GET',
				`data/humidity/${sensor.value}/latest`
			).then(data => {
				var d = latestHumidityData;
				d[sensor.value] = data['humidity'];
				setLatestHumidityData(d);
			}).catch((err) => {
				console.error(`error fetching latest humidity data for ${sensor.value}:`, err);
			})

			doApiCall(
				token,
				'GET',
				`data/pressure/${sensor.value}/latest`
			).then(data => {
				var d = latestPressureData;
				d[sensor.value] = data['pressure'];
				setLatestPressureData(d);
			}).catch((err) => {
				console.error(`error fetching latest pressure data for ${sensor.value}:`, err);
			})
		}
	}, [gotSettings]);

	useEffect(() => {
		// fetch range data
		if (!gotSettings || !isAuthenticated || !activeQuery || !activeSensor) {
			return;
		}
		doApiCall(
			token,
			'GET',
			`data/${activeQuery}/${activeSensor}/range?from=${startTime.toISOString()}&to=${stopTime.toISOString()}&interval=${getInterval(timePeriod)}`
		).then(data => {
			setRangeData(data);
		}).catch((error) => {
			console.log("error getting range data: ", error);
			setRangeData([]);
		});
	}, [gotSettings, activeQuery, activeSensor, startTime, stopTime, token, timePeriod, isAuthenticated]);

	useEffect(() => {
		setActiveUnit(unitFromQuery(activeQuery));
		dispatch(setSelectedMeasurement(activeQuery));
		setMeasurementDropdownButtonTitle(`Measurement: ${activeQuery}`);
	}, [activeQuery]);

	useEffect(() => {
		dispatch(setSelectedSensor(activeSensor));
		setSensorDropdownButtonTitle(`Sensor: ${getSensorName(settings, activeSensor)}`);
	}, [settings, activeSensor]);

	useEffect(() => {
		setSuggestedMin(settings.suggestedMins[`${activeQuery}`]);
		setSuggestedMax(settings.suggestedMaxs[`${activeQuery}`]);
	}, [settings, activeQuery]);

	useEffect(() => {
		dispatch(setSelectedTimePeriod(timePeriod));
		setTimePeriodButtonTitle(`Period: ${getPeriodString(timePeriod)}`);
		setStartTime(calculateRangeStart(timePeriod));
	}, [timePeriod]);

	const getLatestDataByQuery = (query) => {
		switch (query) {
			case 'temperature':
				return latestTemperatureData;
			case 'humidity':
				return latestHumidityData;
			case 'pressure':
				return latestPressureData;
			default:
				return {};
		}
	};

	return (
		<Container fluid='true' id='homeView'>
			<h4>Latest readings</h4>
			<Table>
				<thead><tr key='headerRow'>
					{settings.sensors.map((sensor) => (
						<th>{sensor.name}</th>
					))}
				</tr></thead>
				<tbody><tr key='bodyRow'>
					{settings.sensors.map((sensor) => (
						<td>{getLatestDataByQuery(activeQuery)[sensor.value]
							? `${getLatestDataByQuery(activeQuery)[sensor.value].toFixed(1)} ${activeUnit}`
							: "N/A"}</td>
					))}
				</tr></tbody>
			</Table>
			<h3>Historical data</h3>
			<Navbar bg='light' expand='md'>
				<Navbar.Toggle aria-controls='parameterSelectorNavbar'/>
				<Navbar.Collapse className="justify-content-begin">
					<DropdownButton
						id='measurementDropdownButton'
						variant='outline-primary'
						title={measurementDropdownButtonTitle}
						size='sm'
					>
						<Dropdown.Item onClick={() => { setActiveQuery('temperature') }}>Temperature</Dropdown.Item>
						<Dropdown.Item onClick={() => { setActiveQuery('humidity') }}>Humidity</Dropdown.Item>
						<Dropdown.Item onClick={() => { setActiveQuery('pressure') }}>Air pressure</Dropdown.Item>
					</DropdownButton>
					<DropdownButton
						id='sensorDropdownButton'
						variant='outline-primary'
						title={sensorDropdownButtonTitle}
						size='sm'
					>
						{settings.sensors.map((sensor) => (
							<Dropdown.Item onClick={() => { setActiveSensor(sensor.value) }}>{sensor.name}</Dropdown.Item>
						))}
					</DropdownButton>
					<DropdownButton
						id='timePeriodButton'
						variant='outline-primary'
						title={timePeriodButtonTitle}
						size='sm'
					>
						<Dropdown.Item onClick={() => { setTimePeriod(oneDay) }}>1 day</Dropdown.Item>
						<Dropdown.Item onClick={() => { setTimePeriod(threeDays) }}>3 days</Dropdown.Item>
						<Dropdown.Item onClick={() => { setTimePeriod(oneWeek) }}>1 week</Dropdown.Item>
						<Dropdown.Item onClick={() => { setTimePeriod(oneMonth) }}>1 month</Dropdown.Item>
					</DropdownButton>
				</Navbar.Collapse>
			</Navbar>

			<MyChart
				data={rangeData}
				measurement={activeQuery}
				period={getPeriodString(timePeriod)}
				suggestedMin={suggestedMin}
				suggestedMax={suggestedMax}
			/>

		</Container>
	);
}

export default Home;

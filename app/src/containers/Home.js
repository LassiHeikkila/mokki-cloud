import React, {useState, useEffect} from 'react';
import Chart from "./Chart";
import NotLoggedIn from "./NotLoggedIn";
import { useAppContext } from "../lib/contextLib";
import './Home.css';

function unitFromQuery(query) {
	switch(query) {
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

export default function Home() {
	var now = new Date();
	var rangeStart = new Date();
	rangeStart.setHours(rangeStart.getHours() - (24 * 7));
	const interval = 30 * 60;

	const [activeQuery, setActiveQuery] = useState("temperature");
	const [measurementData, setMeasurementData] = useState({});
	const [activeSensor, setActiveSensor] = useState("");
	const [activeUnit, setActiveUnit] = useState("Celsius");
	const [rangeData, setRangeData] = useState([]);
	const [startTime] = useState(rangeStart.toISOString());
	const [stopTime] = useState(now.toISOString());
	const [gotSettings, setGotSettings] = useState(false);
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

	var context = useAppContext();

	useEffect(() => {
		async function getSettings() {
			if (gotSettings) {
				// no need to get them again
				return;
			}
			fetch(
				'/static/settings.json'
			).then(response => response.json()
			).then(data => {
				if (data.sensors != null && data.sensors.length > 0) {
					setSettings(data);
					setGotSettings(true);
					setActiveSensor(data.sensors[0].value);
				} else {
					throw new Error("settings don't contain any sensors")
				}
			}).catch((error)=> {
				console.log("error fetching settings: ", error)
			});
		};
		async function getLatestData() {
			if (!context.isAuthenticated) {
				console.error("not authenticated, will not attempt to fetch data")
				return
			}
			fetch(
				`/api/data/${activeQuery}/${activeSensor}/latest`, {
					method: "GET",
					mode: "cors",
					headers: {
						'X-API-KEY': context.storedToken
					}
				}
			).then(response => response.json()
			).then(data => {
				setMeasurementData(data);
			}).catch((error) => {
				console.log("error getting latest measurement: ", error)
				setMeasurementData({});
			});
		};

		async function getRangeData() {
			if (!context.isAuthenticated) {
				console.error("not authenticated, will not attempt to fetch data")
				return
			}
			fetch(
				`/api/data/${activeQuery}/${activeSensor}/range?from=${startTime}&to=${stopTime}&interval=${interval}`, {
					method: "GET",
					mode: "cors",
					headers: {
						'X-API-KEY': context.storedToken
					}
				}
			).then(response => response.json()
			).then(data => {
				setRangeData(data);
				const key = `${activeQuery}`;
				setActiveUnit(unitFromQuery(key));
				if (!context.isAuthenticated) {
					return;
				}
			}).catch((error) => {
				console.log("error getting range data: ", error);
				setRangeData([]);
			});
		};
		getSettings();
		getLatestData();
		getRangeData();
	}, [activeQuery, activeSensor, startTime, stopTime, context, interval, settings, gotSettings]);

	console.log(measurementData)
	if (!context.isAuthenticated) {
		return <NotLoggedIn />
	}
	return (
		<div style={{textAlign:"center"}}>
			<p>
			<label>
				<select value={activeQuery} onChange={e => setActiveQuery(e.target.value)}>
					<option key="temperature" value="temperature">Temperature</option>
					<option key="humidity" value="humidity">Humidity</option>
					<option key="pressure" value="pressure">Pressure</option>
				</select>
			</label>
			<label>
				<select value={activeSensor} onChange={e => setActiveSensor(e.target.value)}>
					{settings.sensors.map((e) => {
						return <option key={e.value} value={e.value}>{e.name}</option>;
					})}
				</select>
			</label>
			</p>
			<p>
				Latest measured {activeQuery} is {parseFloat(measurementData[`${activeQuery}`]).toFixed(2)} {activeUnit}.
			</p>
			<p>
				Measured at {(new Date(measurementData[`time`])).toLocaleString()}.
			</p>
			<Chart data={rangeData} measurement={activeQuery} period='week' suggestedMin={settings.suggestedMins[`${activeQuery}`]} suggestedMax={settings.suggestedMaxs[`${activeQuery}`]} unit={unitFromQuery(`${activeQuery}`)}/>
		</div>
	);
}

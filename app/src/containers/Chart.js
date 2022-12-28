import React from 'react';

import { Chart as ChartJS, LineElement, PointElement, LinearScale, TimeScale, Tooltip, Decimation } from 'chart.js';
import { Line } from 'react-chartjs-2';
import 'chartjs-adapter-date-fns';

ChartJS.register(LineElement, PointElement, LinearScale, TimeScale, Tooltip, Decimation);

const options = {
	responsive: true,
	interaction: {
		mode: 'nearest',
		axis: 'x',
		intersect: false,
	},
	plugins: {
		decimation: {
			enabled: true,
			samples: 70,
			algorithm: 'lttb',
		},
	},
	scales: {
		x: {
			type: 'time',
			time: {
				unit: 'day'
			}
		}
	},
	datasets: {
		line: {
			borderColor: '#36a2eb'
		}
	},
};

function getTimestamps(data) {
	var ts = [];
	var i;
	for (i = 0; i < data.length; i++) {
		ts[i] = (new Date(data[i]["time"]));
	}

	return ts;
}

function getValues(data, key) {
	var vals = [];
	var i;
	for (i = 0; i < data.length; i++) {
		vals[i] = data[i][key];
	}

	return vals;
}

const MyChart = (props) => {
	const ts = getTimestamps(props.data);
	const vals = getValues(props.data, props.measurement);

	const d = {
		labels: ts,
		datasets: [
			{
				label: props.measurement,
				data: vals
			}
		]
	};
	var opts = options;
	opts.scales.y = {
		suggestedMin: props.suggestedMin,
		suggestedMax: props.suggestedMax,
		title: {
			text: props.unit,
			display: true,
			align: 'end',
		}
	}

	return (
		<div className='header'>
			<h2>Measurements from past {props.period}</h2>
			<Line data={d} options={opts} />
		</div>
	)
}

export default MyChart;

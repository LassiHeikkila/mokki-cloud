import React from 'react';
import { Line } from 'react-chartjs-2';
import 'chartjs-adapter-date-fns';

const options = {
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
	}
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

const Chart = (props) => {
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

export default Chart;

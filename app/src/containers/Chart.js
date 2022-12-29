import React from 'react';

import { Chart as ChartJS, LineElement, PointElement, LinearScale, TimeScale, Tooltip, Decimation } from 'chart.js';
import { Line } from 'react-chartjs-2';
import 'chartjs-adapter-date-fns';

import Container from 'react-bootstrap/Container';

import { unitFromQuery } from '../lib/unitConversions';

ChartJS.register(LineElement, PointElement, LinearScale, TimeScale, Tooltip, Decimation);

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
	return (
		<Container size='md'>
			<Line 
				data={{
					labels: getTimestamps(props.data),
					datasets: [{
						label: props.measurement,
						data: getValues(props.data, props.measurement),
					}]
				}}
				options={{
					responsive: true,
					interaction: {
						mode: 'nearest',
						axis: 'x',
						intersect: false,
					},
					scales: {
						x: {
							type: 'time',
							time: {
								unit: 'day'
							}
						},
						y: {
							suggestedMin: props.suggestedMin,
							suggestedMax: props.suggestedMax,
							title: {
								text: unitFromQuery(props.measurement),
								display: true,
								align: 'end',
							}
						}
					},
					datasets: {
						line: { borderColor: '#36a2eb' }
					},
				}} 
			/>
		</Container>
	)
}

export default MyChart;

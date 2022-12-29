import { oneDay, threeDays, oneWeek, oneMonth, second, minute, hour } from './units';

const unitFromQuery = (query) => {
	switch (query) {
    case "temperature":
        return "°C";
    case "pressure":
        return "Pa";
    case "humidity":
        return "%";
    default:
        return "unknown";
	}
};

const getInterval = (period) => {
	switch (period) {
		case oneDay:
			return 30 * minute;
		case threeDays:
			return 1 * hour;
		case oneWeek:
			return 2 * hour;
		case oneMonth:
			return 6 * hour;
	}
	// assume that anything besides those options will be longer than a month
	return 24 * hour;
};

const getPeriodString = (period) => {
	switch (period) {
		case oneDay:
			return "1 day";
		case threeDays:
			return "3 days";
		case oneWeek:
			return "1 week";
		case oneMonth:
			return "1 month";
	}
	// assume that anything besides those options will be longer than a month
	return "unknown";
};

export { unitFromQuery, getInterval, getPeriodString };
import { oneDay, threeDays, oneWeek, oneMonth } from './units';

const calculateRangeStart = (period) => {
	const start = new Date();
	switch (period) {
		case oneDay:
			start.setHours(start.getHours() - 24);
			break;
		case threeDays:
			start.setHours(start.getHours() - (24 * 3));
			break;
		case oneWeek:
			start.setHours(start.getHours() - (24 * 7));
			break;
		case oneMonth:
			start.setHours(start.getHours() - (24 * 31));
			break;
		default:
			// ? no good default with chosen way to implement => default to 1 week
			start.setHours(start.getHours() - (24 * 7));
	}
	return start;
};

export { calculateRangeStart };
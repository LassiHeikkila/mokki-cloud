package server

import (
	"context"
	"fmt"
	"log"
	"time"
)

const (
	fromBucket = `from(bucket: "%s")`

	queryLastRecord   = `|> last()`
	queryLastDuration = `|> range(start: -%v)`
	queryBetweenTimes = `|> range(start: %s, stop: %s)`

	filterSensorMAC   = `|> filter(fn: (r) => r["sensormac"] == "%s")`
	filterMeasurement = `|> filter(fn: (r) => r["_measurement"] == "%s")`

	filterField = `|> filter(fn: (r) => r["_field"] == "%s")`

	aggregate = `|> aggregateWindow(every: %s, fn: mean, createEmpty: false) |> yield(name: "mean")`
)

// QueryLastValue assumes there is some data in the past 3h
func (q *Querier) QueryLastValue(ctx context.Context, bucket, field, sensorID, measurement string) Measurement {
	query := fmt.Sprintf(fromBucket, bucket)
	query += fmt.Sprintf(queryLastDuration, 24*time.Hour)
	query += fmt.Sprintf(filterSensorMAC, sensorID)
	query += fmt.Sprintf(filterMeasurement, measurement)
	query += fmt.Sprintf(filterField, field)
	query += queryLastRecord

	log.Println("running query:", query)

	records, err := q.ExecuteQuery(ctx, query)
	if err != nil {
		log.Println("error running query:", err)
		return nil
	}

	if len(records) < 1 {
		log.Println("no records found")
		return nil
	}
	// since we only want the last record, assume there is at most one record returned
	for _, r := range records {
		m, err := MeasurementFromRecord(r)
		if err != nil {
			log.Printf("hep: %e", err)
			continue
		}
		if m.Measurement() == field {
			return m
		}
	}
	log.Printf("no %s measurement found\n", field)
	return nil
}

func (q *Querier) QueryLastDuration(
	ctx context.Context,
	bucket, field, sensorID, measurement string,
	duration time.Duration,
	interval time.Duration,
) []Measurement {
	query := fmt.Sprintf(fromBucket, bucket)
	query += fmt.Sprintf(queryLastDuration, duration)
	query += fmt.Sprintf(filterSensorMAC, sensorID)
	query += fmt.Sprintf(filterMeasurement, measurement)
	query += fmt.Sprintf(filterField, field)
	query += fmt.Sprintf(aggregate, interval)

	log.Println("running query:", query)

	records, err := q.ExecuteQuery(ctx, query)
	if err != nil {
		log.Println("error running query:", err)
		return nil
	}

	if len(records) < 1 {
		log.Println("no records found")
		return nil
	}

	var measurements []Measurement
	var errorCount int
	for _, record := range records {
		m, err := MeasurementFromRecord(record)
		if m == nil || err != nil {
			errorCount++
			continue
		}
		measurements = append(measurements, m)
	}

	log.Printf("%d records failed to be converted to measurements", errorCount)

	return measurements
}

func (q *Querier) QueryBetweenTimes(
	ctx context.Context,
	bucket, field, sensorID, measurement string,
	start, stop time.Time,
	interval time.Duration,
) []Measurement {
	query := fmt.Sprintf(fromBucket, bucket)
	query += fmt.Sprintf(queryBetweenTimes, start.Format(time.RFC3339Nano), stop.Format(time.RFC3339Nano))
	query += fmt.Sprintf(filterSensorMAC, sensorID)
	query += fmt.Sprintf(filterMeasurement, measurement)
	query += fmt.Sprintf(filterField, field)
	query += fmt.Sprintf(aggregate, interval)

	log.Println("running query:", query)

	records, err := q.ExecuteQuery(ctx, query)
	if err != nil {
		log.Println("error running query:", err)
		return nil
	}

	if len(records) < 1 {
		log.Println("no records found")
		return nil
	}

	var measurements []Measurement
	var errorCount int
	for _, record := range records {
		m, err := MeasurementFromRecord(record)
		if m == nil || err != nil {
			errorCount++
			continue
		}
		measurements = append(measurements, m)
	}

	log.Printf("%d records failed to be converted to measurements", errorCount)

	return measurements
}

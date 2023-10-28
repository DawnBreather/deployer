package deployment

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	HTTP_HEALTH_CHECKS_RANGES_CACHE = map[string]ValidStatusCodes{}
)

func (h *HealthChecks) StartValidatingHealthChecks(d *Deployment) {

	go h.TCP.StartValidatingHealthChecks(d)
	go h.Mssql.StartValidatingHealthChecks(d)
	go h.Redis.StartValidatingHealthChecks(d)
	go h.HTTP.StartValidatingHealthChecks(d)
	//TODO: implement h.Filesystem.StartValidatingHealthChecks(d)

	go func() {
		for {
			//PutIntoRedis(fmt.Sprintf(string(RDB_AGENT_HEARTBEAT_PATH), d.Name(), d.tierName(), agentId), GetTimeNowString(), 2*time.Second)
			GetNodeStatus().AgentHeartbeatAt = GetTimeNowString()

			time.Sleep(1 * time.Second)
		}
	}()

}

// TODO: Implement healthy and unhealth thresholds
func (t *Tcp) StartValidatingHealthChecks(d *Deployment) {

	for {

		if !IsSecretsDecryptionAccomplished {
			time.Sleep(1000 * time.Millisecond)
		}

		// Validate healthchecks
		var message string
		var status = "healthy"

		healthChecksStatus := []NetworkHealthCheckStatus{}

		if len(t.Targets) > 0 {

			for _, check := range t.Targets {

				host := TransformValuePlaceholderIntoValue(d.Secrets, check.Host)
				port := TransformValuePlaceholderIntoValue(d.Secrets, check.Port)

				go func() {
					target := fmt.Sprintf("%s:%s", host, port)

					d := net.Dialer{Timeout: time.Duration(t.Options.Timeout) * time.Second}
					_, err := d.Dial("tcp", target)
					if err != nil {
						logrus.Warnf("WARNING Failed passing TCP healthcheck { %s } with timeout { %d }: %v", target, t.Options.Timeout, err)
						message = err.Error()
						status = "unhealthy"
					}
					healthChecksStatus = append(healthChecksStatus, NetworkHealthCheckStatus{
						UpdatedAt: GetTimeNowString(),
						Target:    target,
						Status:    status,
						Message:   message,
					})
				}()
			}
		} else {
			status = ""
		}

		for len(healthChecksStatus) != len(t.Targets) {
			time.Sleep(200 * time.Millisecond)
		}

		// Send healthchecks statuses

		var interval = t.Options.Interval
		if t.Options.Interval == 0 {
			interval = 60
		}

		GetNodeStatus().HealthChecks.TCP = healthChecksStatus

		//if status != "" {
		//	for index, status := range healthChecksStatus {
		//		PutIntoRedis(fmt.Sprintf(string(RDB_TCP_HEALTH_CHECKS_STATUS_PATH), d.Name(), d.tierName(), agentId, index), status, time.DurationSeconds(interval*2)*time.Second)
		//	}
		//} else {
		//	PutIntoRedis(fmt.Sprintf(string(RDB_TCP_HEALTH_CHECKS_STATUS_PATH), d.Name(), d.tierName(), agentId, 0), NetworkHealthCheckStatus{
		//		Status: status,
		//	}, time.DurationSeconds(interval*2)*time.Second)
		//}
		time.Sleep(time.Duration(interval) * time.Second)
	}

}

// TODO: Implement healthy and unhealth thresholds
func (m *Mssql) StartValidatingHealthChecks(d *Deployment) {

	for {

		if !IsSecretsDecryptionAccomplished {
			time.Sleep(1000 * time.Millisecond)
		}

		// Validate healthchecks
		var message string
		var status = "healthy"

		healthChecksStatus := []NetworkHealthCheckStatus{}

		if len(m.Targets) > 0 {
			for _, check := range m.Targets {

				username := TransformValuePlaceholderIntoValue(d.Secrets, check.Username)
				password := TransformValuePlaceholderIntoValue(d.Secrets, check.Password)
				host := TransformValuePlaceholderIntoValue(d.Secrets, check.Host)
				port := TransformValuePlaceholderIntoValue(d.Secrets, check.Port)
				database := TransformValuePlaceholderIntoValue(d.Secrets, check.Database)

				connectionString := fmt.Sprintf(`sqlserver://%s:%s@%s:%s?database=%s&connection+timeout=%d`, username, password, host, port, database, m.Options.Timeout)
				connectionStringCleaned := strings.ReplaceAll(strings.ReplaceAll(connectionString, username, "***"), password, "***")
				db, err := sql.Open("mssql", connectionString)
				err = db.Ping()
				if err != nil {
					logrus.Warnf("WARNING Failed passing MSSQL healthcheck { %s } with timeout { %d }: %v", connectionStringCleaned, m.Options.Timeout, err)
					message = err.Error()
					status = "unhealthy"
				}
				healthChecksStatus = append(healthChecksStatus, NetworkHealthCheckStatus{
					UpdatedAt: GetTimeNowString(),
					Target:    connectionStringCleaned,
					Status:    status,
					Message:   message,
				})
				err = db.Close()
				if err != nil {
					logrus.Errorf("[E] closing database connection { %s }: %v", connectionStringCleaned, err)
				}
			}
		} else {
			//status = "N/A"
			status = ""
		}

		for len(healthChecksStatus) != len(m.Targets) {
			time.Sleep(200 * time.Millisecond)
		}

		// Send healthchecks statuses

		var interval = m.Options.Interval
		if m.Options.Interval == 0 {
			interval = 60
		}

		//if status != "N/A" {
		//	for index, status := range healthChecksStatus {
		//		PutIntoRedis(fmt.Sprintf(string(RDB_MSSQL_HEALTH_CHECKS_STATUS_PATH), d.Name(), d.tierName(), agentId, index), status, time.DurationSeconds(interval*2)*time.Second)
		//	}
		//} else {
		//	PutIntoRedis(fmt.Sprintf(string(RDB_MSSQL_HEALTH_CHECKS_STATUS_PATH), d.Name(), d.tierName(), agentId, 0), NetworkHealthCheckStatus{
		//		Status: status,
		//	}, time.DurationSeconds(interval*2)*time.Second)
		//}

		GetNodeStatus().HealthChecks.Mssql = healthChecksStatus

		time.Sleep(time.Duration(interval) * time.Second)

	}

}

// TODO: Implement healthy and unhealth thresholds
// TODO: Implement TLS support
func (r *Redis) StartValidatingHealthChecks(d *Deployment) {

	for {

		if !IsSecretsDecryptionAccomplished {
			time.Sleep(1000 * time.Millisecond)
		}

		// Validate healthchecks
		var message string
		var status = "healthy"

		healthChecksStatus := []NetworkHealthCheckStatus{}

		if len(r.Targets) > 0 {
			for _, check := range r.Targets {

				username := TransformValuePlaceholderIntoValue(d.Secrets, check.Username)
				password := TransformValuePlaceholderIntoValue(d.Secrets, check.Password)
				host := TransformValuePlaceholderIntoValue(d.Secrets, check.Host)
				port := TransformValuePlaceholderIntoValue(d.Secrets, check.Port)

				socket := fmt.Sprintf("%s:%s", host, port)
				rdb := redis.NewClient(&redis.Options{
					Addr:        socket,
					Username:    username,
					Password:    password,
					DialTimeout: time.Duration(r.Options.Timeout) * time.Second,
				})
				redisConnectivityStatus := rdb.Ping(context.TODO())
				err := redisConnectivityStatus.Err()
				if err != nil {
					logrus.Warnf("WARNING Failed passing Redis healthcheck { %s } with timeout { %d }: %v", socket, r.Options.Timeout, err)
					message = err.Error()
					status = "unhealthy"
				}
				healthChecksStatus = append(healthChecksStatus, NetworkHealthCheckStatus{
					UpdatedAt: GetTimeNowString(),
					Target:    socket,
					Status:    status,
					Message:   message,
				})
				err = rdb.Close()
				if err != nil {
					logrus.Errorf("[E] closing redis connection { %s }: %v", socket, err)
				}
			}
		} else {
			//status = "N/A"
			status = ""
		}

		for len(healthChecksStatus) != len(r.Targets) {
			time.Sleep(200 * time.Millisecond)
		}

		// Send healthchecks statuses

		var interval = r.Options.Interval
		if r.Options.Interval == 0 {
			interval = 60
		}

		//if status != "N/A" {
		//	for index, status := range healthChecksStatus {
		//		PutIntoRedis(fmt.Sprintf(string(RDB_REDIS_HEALTH_CHECKS_STATUS_PATH), d.Name(), d.tierName(), agentId, index), status, time.DurationSeconds(interval*2)*time.Second)
		//	}
		//} else {
		//	PutIntoRedis(fmt.Sprintf(string(RDB_REDIS_HEALTH_CHECKS_STATUS_PATH), d.Name(), d.tierName(), agentId, 0), NetworkHealthCheckStatus{Status: status}, time.DurationSeconds(interval*2)*time.Second)
		//}

		GetNodeStatus().HealthChecks.Redis = healthChecksStatus

		time.Sleep(time.Duration(interval) * time.Second)
	}
}

// TODO: Implement healthy and unhealth thresholds
// TODO: Implement validation of HTTP status codes according to the parameters
func (h *Http) StartValidatingHealthChecks(d *Deployment) {

	if !IsSecretsDecryptionAccomplished {
		time.Sleep(1000 * time.Millisecond)
	}

	for {

		// Validate healthchecks
		var message string
		var status = "healthy"

		healthChecksStatus := []NetworkHealthCheckStatus{}

		if len(h.Targets) > 0 {
			for _, check := range h.Targets {
				checkURL := TransformValuePlaceholderIntoValue(d.Secrets, check.URL)
				req, err := http.NewRequest(http.MethodGet, checkURL, nil)
				if err != nil {
					logrus.Errorf("[E] creating HTTP request for { %+v }: %v", check, err)
				}
				httpClient := http.Client{
					Timeout: time.Duration(h.Options.Timeout) * time.Second,
				}

				resp, err := httpClient.Do(req)

				//TODO: add support of Status validation based on provided value (not hardcoded as in the code line below)
				if err != nil || getValidStatusCodesRange(check.Status).IsValid(resp.StatusCode) != true {
					if err == nil {
						err = fmt.Errorf("status code = %d", resp.StatusCode)
					}
					logrus.Warnf("WARNING Failed passing HTTP healthcheck for { %s } expecting status codes { %s } with timeout { %d }: %v", checkURL, check.Status, h.Options.Timeout, err)
					message = err.Error()
					status = "unhealthy"
				}
				healthChecksStatus = append(healthChecksStatus, NetworkHealthCheckStatus{
					UpdatedAt: GetTimeNowString(),
					Target:    checkURL,
					Status:    status,
					Message:   message,
				})
			}
		} else {
			//status = "N/A"
			status = ""
		}

		for len(healthChecksStatus) != len(h.Targets) {
			time.Sleep(200 * time.Millisecond)
		}

		// Send healthchecks statuses

		var interval = h.Options.Interval
		if h.Options.Interval == 0 {
			interval = 60
		}

		//if status != "N/A" {
		//	for index, status := range healthChecksStatus {
		//		PutIntoRedis(fmt.Sprintf(string(RDB_HTTP_HEALTH_CHECKS_STATUS_PATH), d.Name(), d.tierName(), agentId, index), status, time.DurationSeconds(interval*2)*time.Second)
		//	}
		//} else {
		//	PutIntoRedis(fmt.Sprintf(string(RDB_HTTP_HEALTH_CHECKS_STATUS_PATH), d.Name(), d.tierName(), agentId, 0), NetworkHealthCheckStatus{
		//		Status: status,
		//	}, time.DurationSeconds(interval*2)*time.Second)
		//}

		GetNodeStatus().HealthChecks.HTTP = healthChecksStatus

		time.Sleep(time.Duration(interval) * time.Second)

	}
}

type ValidStatusCodes struct {
	AtomicRanges []ValidStatusCodesAtomicRange
	AtomicValues []int
}

type ValidStatusCodesAtomicRange struct {
	Min int
	Max int
}

func (vs ValidStatusCodes) IsValid(probe int) bool {

	regex := "^([1-5][0-9][0-9])$"
	match, err := regexp.MatchString(regex, fmt.Sprintf("%d", probe))
	if match != true && err != nil {
		return false
	}

	for _, atomicRange := range vs.AtomicRanges {
		if probe >= atomicRange.Min && probe <= atomicRange.Max {
			return true
		}
	}

	for _, atomicValue := range vs.AtomicValues {
		if probe == atomicValue {
			return true
		}
	}

	return false
}

func (vs ValidStatusCodes) ToString() (res string) {
	for _, atomicRange := range vs.AtomicRanges {
		res += fmt.Sprintf("%d-%d,", atomicRange.Min, atomicRange.Max)
	}
	for _, atomicValue := range vs.AtomicValues {
		res += fmt.Sprintf("%d,", atomicValue)
	}
	return strings.TrimSuffix(res, ",")
}

func getValidStatusCodesRange(rangeString string) (res ValidStatusCodes) {
	if cachedResultingRange, isMapContainsKey := HTTP_HEALTH_CHECKS_RANGES_CACHE[rangeString]; isMapContainsKey {
		return cachedResultingRange
	} else {
		ParseHttpStatusCode(rangeString, &res, "", -1)
	}
	return res
}

func ParseHttpStatusCode(rangeString string, resultingRange *ValidStatusCodes, lastDelimiter string, index int) error {

	regex := "^(([1-5][0-9][0-9])|([1-5][0-9][0-9]-[1-5][0-9][0-9]))(,(([1-5][0-9][0-9])|([1-5][0-9][0-9]-[1-5][0-9][0-9])))*$"
	defaultCodes := ValidStatusCodes{
		AtomicRanges: []ValidStatusCodesAtomicRange{
			{
				Min: 200,
				Max: 399,
			},
		},
	}

	match, err := regexp.MatchString(regex, rangeString)
	if err != nil || match != true {
		*resultingRange = defaultCodes
		switch {
		case err != nil:
			logrus.Warnf("[W] Error parsing HTTP status codes range { %s } for health checks (falling back to deafult { %s } status codes): err { %v }", rangeString, defaultCodes.ToString(), err)
			break
		case err == nil:
			logrus.Warnf("[W] Error parsing HTTP status codes range { %s } for health checks (falling back to deafult { %s } status codes)", rangeString, defaultCodes.ToString())
			break
		}
		return nil
	}

	switch {
	case strings.Contains(rangeString, ","):
		parts := strings.SplitN(rangeString, ",", 2)
		if e := ParseHttpStatusCode(parts[0], resultingRange, ",", 0); e != nil {
			return e
		}
		if e := ParseHttpStatusCode(parts[1], resultingRange, ",", 1); e != nil {
			return e
		}
		break
	case strings.Contains(rangeString, "-"):
		parts := strings.SplitN(rangeString, "-", 2)
		if e := ParseHttpStatusCode(parts[0], resultingRange, "-", 0); e != nil {
			return e
		}
		if e := ParseHttpStatusCode(parts[1], resultingRange, "-", 1); e != nil {
			return e
		}
		break
	case rangeString == "" && lastDelimiter == "":
		//if len(resultingRange.AtomicRanges) == 0 && len(resultingRange.AtomicValues) == 0 {
		*resultingRange = defaultCodes
		//}
		break
	default:
		val, err := strconv.Atoi(rangeString)
		if err != nil {
			logrus.Errorf("[E] parsing HTTP status code string value { %s } to int: %v", rangeString, err)
		}

		switch {
		case lastDelimiter == ",":
			resultingRange.AtomicValues = append(resultingRange.AtomicValues, val)
			break
		case lastDelimiter == "-":
			switch index {
			case 0:
				resultingRange.AtomicRanges = append(resultingRange.AtomicRanges, ValidStatusCodesAtomicRange{
					Min: val,
					Max: -1,
				})
				break
			case 1:
				resultingRange.AtomicRanges[len(resultingRange.AtomicRanges)-1].Max = val
			}
			break
		case lastDelimiter == "":
			resultingRange.AtomicValues = append(resultingRange.AtomicValues, val)
			break
		}
	}

	return nil
}

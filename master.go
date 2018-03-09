package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

func newMasterCollector(httpClient *httpClient) prometheus.Collector {
	metrics := map[prometheus.Collector]func(metricMap, prometheus.Collector) error{
		// CPU/Disk/Mem resources in free/used
		gauge("master", "cpus", "Current CPU resources in cluster.", "type"): func(m metricMap, c prometheus.Collector) error {
			total, ok := m["master/cpus_total"]
			used, ok := m["master/cpus_used"]
			if !ok {
				return errNotFoundInMap
			}
			c.(*prometheus.GaugeVec).WithLabelValues("free").Set(total - used)
			c.(*prometheus.GaugeVec).WithLabelValues("used").Set(used)
			return nil
		},
		gauge("master", "cpus_revocable", "Current revocable CPU resources in cluster.", "type"): func(m metricMap, c prometheus.Collector) error {
			total, ok := m["master/cpus_revocable_total"]
			used, ok := m["master/cpus_revocable_used"]
			if !ok {
				return errNotFoundInMap
			}
			c.(*prometheus.GaugeVec).WithLabelValues("free").Set(total - used)
			c.(*prometheus.GaugeVec).WithLabelValues("used").Set(used)
			return nil
		},
		gauge("master", "mem", "Current memory resources in cluster.", "type"): func(m metricMap, c prometheus.Collector) error {
			total, ok := m["master/mem_total"]
			used, ok := m["master/mem_used"]
			if !ok {
				return errNotFoundInMap
			}
			c.(*prometheus.GaugeVec).WithLabelValues("free").Set(total - used)
			c.(*prometheus.GaugeVec).WithLabelValues("used").Set(used)
			return nil
		},
		gauge("master", "mem_revocable", "Current revocable memory resources in cluster.", "type"): func(m metricMap, c prometheus.Collector) error {
			total, ok := m["master/mem_revocable_total"]
			used, ok := m["master/mem_revocable_used"]
			if !ok {
				return errNotFoundInMap
			}
			c.(*prometheus.GaugeVec).WithLabelValues("free").Set(total - used)
			c.(*prometheus.GaugeVec).WithLabelValues("used").Set(used)
			return nil
		},
		gauge("master", "disk", "Current disk resources in cluster.", "type"): func(m metricMap, c prometheus.Collector) error {
			total, ok := m["master/disk_total"]
			used, ok := m["master/disk_used"]
			if !ok {
				return errNotFoundInMap
			}
			c.(*prometheus.GaugeVec).WithLabelValues("free").Set(total - used)
			c.(*prometheus.GaugeVec).WithLabelValues("used").Set(used)
			return nil
		},
		gauge("master", "disk_revocable", "Current disk resources in cluster.", "type"): func(m metricMap, c prometheus.Collector) error {
			total, ok := m["master/disk_revocable_total"]
			used, ok := m["master/disk_revocable_used"]
			if !ok {
				return errNotFoundInMap
			}
			c.(*prometheus.GaugeVec).WithLabelValues("free").Set(total - used)
			c.(*prometheus.GaugeVec).WithLabelValues("used").Set(used)
			return nil
		},

		// Master stats about uptime and election state
		prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "mesos",
			Subsystem: "master",
			Name:      "elected",
			Help:      "1 if master is elected leader, 0 if not",
		}): func(m metricMap, c prometheus.Collector) error {
			elected, ok := m["master/elected"]
			if !ok {
				return errNotFoundInMap
			}
			c.(prometheus.Gauge).Set(elected)
			return nil
		},
		prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "mesos",
			Subsystem: "master",
			Name:      "uptime_seconds",
			Help:      "Number of seconds the master process is running.",
		}): func(m metricMap, c prometheus.Collector) error {
			uptime, ok := m["master/uptime_secs"]
			if !ok {
				return errNotFoundInMap
			}
			c.(prometheus.Gauge).Set(uptime)
			return nil
		},
		// Master stats about agents
		counter("master", "slave_registration_events_total", "Total number of registration events on this master since it booted.", "event"): func(m metricMap, c prometheus.Collector) error {
			registrations, ok := m["master/slave_registrations"]
			reregistrations, ok := m["master/slave_reregistrations"]
			if !ok {
				return errNotFoundInMap
			}
			c.(*settableCounterVec).Set(registrations, "register")
			c.(*settableCounterVec).Set(reregistrations, "reregister")
			return nil
		},

		counter("master", "slave_removal_events_total", "Total number of removal events on this master since it booted.", "event"): func(m metricMap, c prometheus.Collector) error {
			scheduled, ok := m["master/slave_shutdowns_scheduled"]
			canceled, ok := m["master/slave_shutdowns_canceled"]
			completed, ok := m["master/slave_shutdowns_completed"]
			removals, ok := m["master/slave_removals"]
			if !ok {
				return errNotFoundInMap
			}
			c.(*settableCounterVec).Set(scheduled, "scheduled")
			c.(*settableCounterVec).Set(canceled, "canceled")
			c.(*settableCounterVec).Set(completed, "completed")
			c.(*settableCounterVec).Set(removals-completed, "died")
			return nil
		},
		gauge("master", "slaves_state", "Current number of slaves known to the master per connection and registration state.", "state"): func(m metricMap, c prometheus.Collector) error {
			active, ok := m["master/slaves_active"]
			inactive, ok := m["master/slaves_inactive"]
			disconnected, ok := m["master/slaves_disconnected"]
			if !ok {
				return errNotFoundInMap
			}

			// FIXME: Make sure those assumptions are right
			// Every "active" node is connected to the master
			c.(*prometheus.GaugeVec).WithLabelValues("connected", "active").Set(active)
			// Every "inactive" node is connected but node sending offers
			c.(*prometheus.GaugeVec).WithLabelValues("connected", "inactive").Set(inactive)
			// Every "disconnected" node is "inactive"
			c.(*prometheus.GaugeVec).WithLabelValues("disconnected", "inactive").Set(disconnected)
			// Every "connected" node is either active or inactive
			return nil
		},

		// Master stats about frameworks
		gauge("master", "frameworks_state", "Current number of frames known to the master per connection and registration state.", "state"): func(m metricMap, c prometheus.Collector) error {
			active, ok := m["master/frameworks_active"]
			inactive, ok := m["master/frameworks_inactive"]
			disconnected, ok := m["master/frameworks_disconnected"]
			if !ok {
				return errNotFoundInMap
			}
			// FIXME: Make sure those assumptions are right
			// Every "active" framework is connected to the master
			c.(*prometheus.GaugeVec).WithLabelValues("connected", "active").Set(active)
			// Every "inactive" framework is connected but framework sending offers
			c.(*prometheus.GaugeVec).WithLabelValues("connected", "inactive").Set(inactive)
			// Every "disconnected" framework is "inactive"
			c.(*prometheus.GaugeVec).WithLabelValues("disconnected", "inactive").Set(disconnected)
			// Every "connected" framework is either active or inactive
			return nil
		},
		prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "mesos",
			Subsystem: "master",
			Name:      "offers_pending",
			Help:      "Current number of offers made by the master which aren't yet accepted or declined by frameworks.",
		}): func(m metricMap, c prometheus.Collector) error {
			offers, ok := m["master/outstanding_offers"]
			if !ok {
				return errNotFoundInMap
			}
			c.(prometheus.Gauge).Set(offers)
			return nil
		},
		// Master stats about tasks
		counter("master", "task_states_exit_total", "Total number of tasks processed by exit state.", "state"): func(m metricMap, c prometheus.Collector) error {
			errored, ok := m["master/tasks_error"]
			failed, ok := m["master/tasks_failed"]
			finished, ok := m["master/tasks_finished"]
			killed, ok := m["master/tasks_killed"]
			lost, ok := m["master/tasks_lost"]
			if !ok {
				return errNotFoundInMap
			}
			c.(*settableCounterVec).Set(errored, "errored")
			c.(*settableCounterVec).Set(failed, "failed")
			c.(*settableCounterVec).Set(finished, "finished")
			c.(*settableCounterVec).Set(killed, "killed")
			c.(*settableCounterVec).Set(lost, "lost")
			return nil
		},
		counter("master", "task_states_current", "Current number of tasks by state.", "state"): func(m metricMap, c prometheus.Collector) error {
			running, ok := m["master/tasks_running"]
			staging, ok := m["master/tasks_staging"]
			starting, ok := m["master/tasks_starting"]
			if !ok {
				return errNotFoundInMap
			}
			c.(*settableCounterVec).Set(running, "running")
			c.(*settableCounterVec).Set(staging, "staging")
			c.(*settableCounterVec).Set(starting, "starting")
			return nil
		},

		// Master stats about messages
		counter("master", "messages_outcomes_total",
			"Total number of messages by outcome of operation and direction.",
			"source", "destination", "type", "outcome"): func(m metricMap, c prometheus.Collector) error {
			frameworkToExecutorValid, ok := m["master/valid_framework_to_executor_messages"]
			frameworkToExecutorInvalid, ok := m["master/invalid_framework_to_executor_messages"]
			executorToFrameworkValid, ok := m["master/valid_executor_to_framework_messages"]
			executorToFrameworkInvalid, ok := m["master/invalid_executor_to_framework_messages"]

			// status updates are sent from framework?(FIXME) to slave
			// status update acks are sent from slave to framework?
			statusUpdateAckValid, ok := m["master/valid_status_update_acknowledgements"]
			statusUpdateAckInvalid, ok := m["master/invalid_status_update_acknowledgements"]
			statusUpdateValid, ok := m["master/valid_status_updates"]
			statusUpdateInvalid, ok := m["master/invalid_status_updates"]
			if !ok {
				return errNotFoundInMap
			}
			c.(*settableCounterVec).Set(frameworkToExecutorValid, "framework", "executor", "", "valid")
			c.(*settableCounterVec).Set(frameworkToExecutorInvalid, "framework", "executor", "", "invalid")

			c.(*settableCounterVec).Set(executorToFrameworkValid, "executor", "framework", "", "valid")
			c.(*settableCounterVec).Set(executorToFrameworkInvalid, "executor", "framework", "", "invalid")

			// We consider a ack message simply as a message from slave to framework
			c.(*settableCounterVec).Set(statusUpdateValid, "framework", "slave", "status_update", "valid")
			c.(*settableCounterVec).Set(statusUpdateInvalid, "framework", "slave", "status_update", "invalid")
			c.(*settableCounterVec).Set(statusUpdateAckValid, "slave", "framework", "status_update", "valid")
			c.(*settableCounterVec).Set(statusUpdateAckInvalid, "slave", "framework", "status_update", "invalid")
			return nil
		},
		counter("master", "messages_type_total", "Total number of valid messages by type.", "type"): func(m metricMap, c prometheus.Collector) error {
			for k, v := range m {
				i := strings.Index("master/messages_", k)
				if i == -1 {
					continue
				}
				// FIXME: We expose things like messages_framework_to_executor twice
				c.(*settableCounterVec).Set(v, k[i:])
			}
			return nil
		},

		// Master stats about events
		gauge("master", "event_queue_length", "Current number of elements in event queue by type", "type"): func(m metricMap, c prometheus.Collector) error {
			dispatches, ok := m["master/event_queue_dispatches"]
			httpRequests, ok := m["master/event_queue_http_requests"]
			messages, ok := m["master/event_queue_messages"]
			if !ok {
				return errNotFoundInMap
			}
			c.(*prometheus.GaugeVec).WithLabelValues("message").Set(messages)
			c.(*prometheus.GaugeVec).WithLabelValues("http_request").Set(httpRequests)
			c.(*prometheus.GaugeVec).WithLabelValues("dispatches").Set(dispatches)
			return nil
		},

		// Master stats about allocations

		prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "mesos",
			Subsystem: "master",
			Name:      "allocation_run_ms_count",
			Help:      "Number of allocation algorithm time measurements in the window",
		}): func(m metricMap, c prometheus.Collector) error {
			count := m["allocator/mesos/allocation_runs"]
			c.(prometheus.Gauge).Set(count)
			return nil
		},

		gauge("master", "allocation_run_ms", "Time spent in allocation algorithm in ms.", "type"): func(m metricMap, c prometheus.Collector) error {
			mean, ok := m["allocator/mesos/allocation_run_ms"]
			min, ok := m["allocator/mesos/allocation_run_ms/min"]
			max, ok := m["allocator/mesos/allocation_run_ms/max"]
			p50, ok := m["allocator/mesos/allocation_run_ms/p50"]
			p90, ok := m["allocator/mesos/allocation_run_ms/p90"]
			p95, ok := m["allocator/mesos/allocation_run_ms/p95"]
			p99, ok := m["allocator/mesos/allocation_run_ms/p99"]
			p999, ok := m["allocator/mesos/allocation_run_ms/p999"]
			p9999, ok := m["allocator/mesos/allocation_run_ms/p9999"]
			if !ok {
				return errNotFoundInMap
			}

			c.(*prometheus.GaugeVec).WithLabelValues("mean").Set(mean)
			c.(*prometheus.GaugeVec).WithLabelValues("min").Set(min)
			c.(*prometheus.GaugeVec).WithLabelValues("max").Set(max)
			c.(*prometheus.GaugeVec).WithLabelValues("p50").Set(p50)
			c.(*prometheus.GaugeVec).WithLabelValues("p90").Set(p90)
			c.(*prometheus.GaugeVec).WithLabelValues("p95").Set(p95)
			c.(*prometheus.GaugeVec).WithLabelValues("p99").Set(p99)
			c.(*prometheus.GaugeVec).WithLabelValues("p999").Set(p999)
			c.(*prometheus.GaugeVec).WithLabelValues("p9999").Set(p9999)
			return nil
		},

		prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "mesos",
			Subsystem: "master",
			Name:      "allocation_runs",
			Help:      "Number of times the allocation algorithm has run",
		}): func(m metricMap, c prometheus.Collector) error {
			count := m["allocator/mesos/allocation_runs"]
			c.(prometheus.Counter).Add(count)
			return nil
		},

		prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "mesos",
			Subsystem: "master",
			Name:      "allocation_run_latency_ms_count",
			Help:      "Number of allocation batch latency measurements in the window",
		}): func(m metricMap, c prometheus.Collector) error {
			count := m["allocator/mesos/allocation_runs"]
			c.(prometheus.Counter).Add(count)
			return nil
		},

		gauge("master", "allocation_run_latency_ms", "Allocation batch latency in ms.", "type"): func(m metricMap, c prometheus.Collector) error {
			mean, ok := m["allocator/mesos/allocation_run_latency_ms"]
			min, ok := m["allocator/mesos/allocation_run_latency_ms/min"]
			max, ok := m["allocator/mesos/allocation_run_latency_ms/max"]
			p50, ok := m["allocator/mesos/allocation_run_latency_ms/p50"]
			p90, ok := m["allocator/mesos/allocation_run_latency_ms/p90"]
			p95, ok := m["allocator/mesos/allocation_run_latency_ms/p95"]
			p99, ok := m["allocator/mesos/allocation_run_latency_ms/p99"]
			p999, ok := m["allocator/mesos/allocation_run_latency_ms/p999"]
			p9999, ok := m["allocator/mesos/allocation_run_latency_ms/p9999"]
			if !ok {
				return errNotFoundInMap
			}

			c.(*prometheus.GaugeVec).WithLabelValues("mean").Set(mean)
			c.(*prometheus.GaugeVec).WithLabelValues("min").Set(min)
			c.(*prometheus.GaugeVec).WithLabelValues("max").Set(max)
			c.(*prometheus.GaugeVec).WithLabelValues("p50").Set(p50)
			c.(*prometheus.GaugeVec).WithLabelValues("p90").Set(p90)
			c.(*prometheus.GaugeVec).WithLabelValues("p95").Set(p95)
			c.(*prometheus.GaugeVec).WithLabelValues("p99").Set(p99)
			c.(*prometheus.GaugeVec).WithLabelValues("p999").Set(p999)
			c.(*prometheus.GaugeVec).WithLabelValues("p9999").Set(p9999)
			return nil
		},

		prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "mesos",
			Subsystem: "master",
			Name:      "event_queue_dispatches",
			Help:      "Number of dispatch events in the event queue.",
		}): func(m metricMap, c prometheus.Collector) error {
			count := m["allocator/mesos/event_queue_dispatches"]
			c.(prometheus.Gauge).Set(count)
			return nil
		},

		gauge("master", "allocator_offer_filters_active", "Number of active offer filters for all frameworks within the role", "role"): func(m metricMap, c prometheus.Collector) error {
			re, err := regexp.Compile("allocator/mesos/offer_filters/roles/(.*?)/active")
			if err != nil {
				return fmt.Errorf("could not compile allocator_offer_filters_active regex: %s", err)
			}
			for metric, value := range m {
				matches := re.FindStringSubmatch(metric)
				if len(matches) != 2 {
					continue
				}
				role := matches[1]
				c.(*prometheus.GaugeVec).WithLabelValues(role).Set(value)
			}
			return nil
		},

		gauge("master", "allocator_role_quota_offered_or_allocated", "Amount of resources considered offered or allocated towards a role's quota guarantee.", "role", "resource"): func(m metricMap, c prometheus.Collector) error {
			re, err := regexp.Compile("allocator/mesos/quota/roles/(.*?)/resources/(.*?)/offered_or_allocated")
			if err != nil {
				return fmt.Errorf("could not compile allocator_role_quota_offered_or_allocated regex: %s", err)
			}
			for metric, value := range m {
				matches := re.FindStringSubmatch(metric)
				if len(matches) != 3 {
					continue
				}
				role := matches[1]
				resource := matches[2]
				c.(*prometheus.GaugeVec).WithLabelValues(role, resource).Set(value)
			}
			return nil
		},

		gauge("master", "allocator_role_quota_guarantee", "Amount of resources guaranteed for a role via quota", "role", "resource"): func(m metricMap, c prometheus.Collector) error {
			re, err := regexp.Compile("allocator/mesos/quota/roles/(.*?)/resources/(.*?)/guarantee")
			if err != nil {
				return fmt.Errorf("could not compile allocator_role_quota_guarantee regex: %s", err)
			}
			for metric, value := range m {
				matches := re.FindStringSubmatch(metric)
				if len(matches) != 3 {
					continue
				}
				role := matches[1]
				resource := matches[2]
				c.(*prometheus.GaugeVec).WithLabelValues(role, resource).Set(value)
			}
			return nil
		},

		gauge("master", "allocator_resources_cpus", "Number of CPUs offered or allocated", "type"): func(m metricMap, c prometheus.Collector) error {
			total, ok := m["allocator/mesos/resources/cpus/total"]
			offeredOrAllocated, ok := m["allocator/mesos/resources/cpus/offered_or_allocated"]
			if !ok {
				return errNotFoundInMap
			}

			c.(*prometheus.GaugeVec).WithLabelValues("total").Set(total)
			c.(*prometheus.GaugeVec).WithLabelValues("offered_or_allocated").Set(offeredOrAllocated)
			return nil
		},

		gauge("master", "allocator_resources_disk", "Allocated or offered disk space in MB", "type"): func(m metricMap, c prometheus.Collector) error {
			total, ok := m["allocator/mesos/resources/disk/total"]
			offeredOrAllocated, ok := m["allocator/mesos/resources/disk/offered_or_allocated"]
			if !ok {
				return errNotFoundInMap
			}

			c.(*prometheus.GaugeVec).WithLabelValues("total").Set(total)
			c.(*prometheus.GaugeVec).WithLabelValues("offered_or_allocated").Set(offeredOrAllocated)
			return nil
		},

		gauge("master", "allocator_resources_mem", "Allocated or offered memory in MB", "type"): func(m metricMap, c prometheus.Collector) error {
			total, ok := m["allocator/mesos/resources/mem/total"]
			offeredOrAllocated, ok := m["allocator/mesos/resources/mem/offered_or_allocated"]
			if !ok {
				return errNotFoundInMap
			}

			c.(*prometheus.GaugeVec).WithLabelValues("total").Set(total)
			c.(*prometheus.GaugeVec).WithLabelValues("offered_or_allocated").Set(offeredOrAllocated)
			return nil
		},
	}
	return newMetricCollector(httpClient, metrics)
}

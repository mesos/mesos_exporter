package main

import (
	"fmt"
	"regexp"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

func newMasterCollector(httpClient *httpClient, hostname string) prometheus.Collector {
	metrics := map[prometheus.Collector]func(metricMap, prometheus.Collector) error{
		// CPU/Disk/Mem resources in free/used
		gauge("master", "cpus", "Current CPU resources in cluster.", "type", "hostname"): func(m metricMap, c prometheus.Collector) error {
			percent, ok := m["master/cpus_percent"]
			if !ok {
				log.WithField("metric", "master/cpus_percent").Warn(LogErrNotFoundInMap)
			}
			total, ok := m["master/cpus_total"]
			if !ok {
				log.WithField("metric", "master/cpus_total").Warn(LogErrNotFoundInMap)
			}
			used, ok := m["master/cpus_used"]
			if !ok {
				log.WithField("metric", "master/cpus_used").Warn(LogErrNotFoundInMap)
			}
			c.(*prometheus.GaugeVec).WithLabelValues("percent", hostname).Set(percent)
			c.(*prometheus.GaugeVec).WithLabelValues("total", hostname).Set(total)
			c.(*prometheus.GaugeVec).WithLabelValues("free", hostname).Set(total - used)
			c.(*prometheus.GaugeVec).WithLabelValues("used", hostname).Set(used)
			return nil
		},
		gauge("master", "cpus_revocable", "Current revocable CPU resources in cluster.", "type", "hostname"): func(m metricMap, c prometheus.Collector) error {
			percent, ok := m["master/cpus_revocable_percent"]
			if !ok {
				log.WithField("metric", "master/cpus_revocable_percent").Warn(LogErrNotFoundInMap)
			}
			total, ok := m["master/cpus_revocable_total"]
			if !ok {
				log.WithField("metric", "master/cpus_revocable_total").Warn(LogErrNotFoundInMap)
			}
			used, ok := m["master/cpus_revocable_used"]
			if !ok {
				log.WithField("metric", "master/cpus_revocable_used").Warn(LogErrNotFoundInMap)
			}
			c.(*prometheus.GaugeVec).WithLabelValues("percent", hostname).Set(percent)
			c.(*prometheus.GaugeVec).WithLabelValues("total", hostname).Set(total)
			c.(*prometheus.GaugeVec).WithLabelValues("free", hostname).Set(total - used)
			c.(*prometheus.GaugeVec).WithLabelValues("used", hostname).Set(used)
			return nil
		},
		gauge("master", "gpus", "Current GPU resources in cluster.", "type", "hostname"): func(m metricMap, c prometheus.Collector) error {
			percent, ok := m["master/gpus_percent"]
			if !ok {
				log.WithField("metric", "master/gpus_percent").Warn(LogErrNotFoundInMap)
			}
			total, ok := m["master/gpus_total"]
			if !ok {
				log.WithField("metric", "master/gpus_total").Warn(LogErrNotFoundInMap)
			}
			used, ok := m["master/gpus_used"]
			if !ok {
				log.WithField("metric", "master/gpus_used").Warn(LogErrNotFoundInMap)
			}
			c.(*prometheus.GaugeVec).WithLabelValues("percent", hostname).Set(percent)
			c.(*prometheus.GaugeVec).WithLabelValues("total", hostname).Set(total)
			c.(*prometheus.GaugeVec).WithLabelValues("free", hostname).Set(total - used)
			c.(*prometheus.GaugeVec).WithLabelValues("used", hostname).Set(used)
			return nil
		},
		gauge("master", "gpus_revocable", "Current revocable GPU resources in cluster.", "type", "hostname"): func(m metricMap, c prometheus.Collector) error {
			percent, ok := m["master/gpus_revocable_percent"]
			if !ok {
				log.WithField("metric", "master/gpus_revocable_percent").Warn(LogErrNotFoundInMap)
			}
			total, ok := m["master/gpus_revocable_total"]
			if !ok {
				log.WithField("metric", "master/gpus_revocable_total").Warn(LogErrNotFoundInMap)
			}
			used, ok := m["master/gpus_revocable_used"]
			if !ok {
				log.WithField("metric", "master/gpus_revocable_used").Warn(LogErrNotFoundInMap)
			}
			c.(*prometheus.GaugeVec).WithLabelValues("percent", hostname).Set(percent)
			c.(*prometheus.GaugeVec).WithLabelValues("total", hostname).Set(total)
			c.(*prometheus.GaugeVec).WithLabelValues("free", hostname).Set(total - used)
			c.(*prometheus.GaugeVec).WithLabelValues("used", hostname).Set(used)
			return nil
		},
		gauge("master", "mem", "Current memory resources in cluster.", "type", "hostname"): func(m metricMap, c prometheus.Collector) error {
			percent, ok := m["master/mem_percent"]
			if !ok {
				log.WithField("metric", "master/mem_percent").Warn(LogErrNotFoundInMap)
			}
			total, ok := m["master/mem_total"]
			if !ok {
				log.WithField("metric", "master/mem_total").Warn(LogErrNotFoundInMap)
			}
			used, ok := m["master/mem_used"]
			if !ok {
				log.WithField("metric", "master/mem_used").Warn(LogErrNotFoundInMap)
			}
			c.(*prometheus.GaugeVec).WithLabelValues("percent", hostname).Set(percent)
			c.(*prometheus.GaugeVec).WithLabelValues("total", hostname).Set(total)
			c.(*prometheus.GaugeVec).WithLabelValues("free", hostname).Set(total - used)
			c.(*prometheus.GaugeVec).WithLabelValues("used", hostname).Set(used)
			return nil
		},
		gauge("master", "mem_revocable", "Current revocable memory resources in cluster.", "type", "hostname"): func(m metricMap, c prometheus.Collector) error {
			percent, ok := m["master/mem_revocable_percent"]
			if !ok {
				log.WithField("metric", "master/mem_revocable_percent").Warn(LogErrNotFoundInMap)
			}
			total, ok := m["master/mem_revocable_total"]
			if !ok {
				log.WithField("metric", "master/mem_revocable_total").Warn(LogErrNotFoundInMap)
			}
			used, ok := m["master/mem_revocable_used"]
			if !ok {
				log.WithField("metric", "master/mem_revocable_used").Warn(LogErrNotFoundInMap)
			}
			c.(*prometheus.GaugeVec).WithLabelValues("percent", hostname).Set(percent)
			c.(*prometheus.GaugeVec).WithLabelValues("total", hostname).Set(total)
			c.(*prometheus.GaugeVec).WithLabelValues("free", hostname).Set(total - used)
			c.(*prometheus.GaugeVec).WithLabelValues("used", hostname).Set(used)
			return nil
		},
		gauge("master", "disk", "Current disk resources in cluster.", "type", "hostname"): func(m metricMap, c prometheus.Collector) error {
			percent, ok := m["master/disk_percent"]
			if !ok {
				log.WithField("metric", "master/disk_percent").Warn(LogErrNotFoundInMap)
			}
			total, ok := m["master/disk_total"]
			if !ok {
				log.WithField("metric", "master/disk_total").Warn(LogErrNotFoundInMap)
			}
			used, ok := m["master/disk_used"]
			if !ok {
				log.WithField("metric", "master/disk_used").Warn(LogErrNotFoundInMap)
			}
			c.(*prometheus.GaugeVec).WithLabelValues("percent", hostname).Set(percent)
			c.(*prometheus.GaugeVec).WithLabelValues("total", hostname).Set(total)
			c.(*prometheus.GaugeVec).WithLabelValues("free", hostname).Set(total - used)
			c.(*prometheus.GaugeVec).WithLabelValues("used", hostname).Set(used)
			return nil
		},
		gauge("master", "disk_revocable", "Current disk resources in cluster.", "type", "hostname"): func(m metricMap, c prometheus.Collector) error {
			percent, ok := m["master/disk_revocable_percent"]
			if !ok {
				log.WithField("metric", "master/disk_revocable_percent").Warn(LogErrNotFoundInMap)
			}
			total, ok := m["master/disk_revocable_total"]
			if !ok {
				log.WithField("metric", "master/disk_revocable_total").Warn(LogErrNotFoundInMap)
			}
			used, ok := m["master/disk_revocable_used"]
			if !ok {
				log.WithField("metric", "master/disk_revocable_used").Warn(LogErrNotFoundInMap)
			}
			c.(*prometheus.GaugeVec).WithLabelValues("percent", hostname).Set(percent)
			c.(*prometheus.GaugeVec).WithLabelValues("total", hostname).Set(total)
			c.(*prometheus.GaugeVec).WithLabelValues("free", hostname).Set(total - used)
			c.(*prometheus.GaugeVec).WithLabelValues("used", hostname).Set(used)
			return nil
		},
		// Master stats about uptime and election state
		prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "mesos",
			Subsystem: "master",
			Name:      "elected",
			Help:      "1 if master is elected leader, 0 if not",
		}, []string{"hostname"}): func(m metricMap, c prometheus.Collector) error {
			elected, ok := m["master/elected"]
			if !ok {
				log.WithField("metric", "master/elected").Warn(LogErrNotFoundInMap)
			}
			c.(*prometheus.GaugeVec).WithLabelValues(hostname).Set(elected)
			return nil
		},
		prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "mesos",
			Subsystem: "master",
			Name:      "uptime_seconds",
			Help:      "Number of seconds the master process is running.",
		}, []string{"hostname"}): func(m metricMap, c prometheus.Collector) error {
			uptime, ok := m["master/uptime_secs"]
			if !ok {
				log.WithField("metric", "master/uptime_secs").Warn(LogErrNotFoundInMap)
			}
			c.(*prometheus.GaugeVec).WithLabelValues(hostname).Set(uptime)
			return nil
		},
		// Master stats about agents
		counter("master", "slave_registration_events_total", "Total number of registration events on this master since it booted.", "event", "hostname"): func(m metricMap, c prometheus.Collector) error {
			registrations, ok := m["master/slave_registrations"]
			if !ok {
				log.WithField("metric", "master/slave_registrations").Warn(LogErrNotFoundInMap)
			}
			reregistrations, ok := m["master/slave_reregistrations"]
			if !ok {
				log.WithField("metric", "master/slave_reregistrations").Warn(LogErrNotFoundInMap)
			}
			c.(*settableCounterVec).Set(registrations, "register", hostname)
			c.(*settableCounterVec).Set(reregistrations, "reregister", hostname)
			return nil
		},

		counter("master", "recovery_slave_removal_events_total", "Total number of recovery removal events on this master since it booted.", "event", "hostname"): func(m metricMap, c prometheus.Collector) error {
			removals, ok := m["master/recovery_slave_removals"]
			if !ok {
				log.WithField("metric", "master/recovery_slave_removals").Warn(LogErrNotFoundInMap)
			}
			c.(*settableCounterVec).Set(removals, "removal", hostname)
			return nil
		},

		counter("master", "slave_removal_events_total", "Total number of removal events on this master since it booted.", "event", "hostname"): func(m metricMap, c prometheus.Collector) error {
			scheduled, ok := m["master/slave_shutdowns_scheduled"]
			if !ok {
				log.WithField("metric", "master/slave_shutdowns_scheduled").Warn(LogErrNotFoundInMap)
			}
			canceled, ok := m["master/slave_shutdowns_canceled"]
			if !ok {
				log.WithField("metric", "master/slave_shutdowns_canceled").Warn(LogErrNotFoundInMap)
			}
			completed, ok := m["master/slave_shutdowns_completed"]
			if !ok {
				log.WithField("metric", "master/slave_shutdowns_completed").Warn(LogErrNotFoundInMap)
			}
			removals, ok := m["master/slave_removals"]
			if !ok {
				log.WithField("metric", "master/slave_removals").Warn(LogErrNotFoundInMap)
			}
			c.(*settableCounterVec).Set(scheduled, "scheduled", hostname)
			c.(*settableCounterVec).Set(canceled, "canceled", hostname)
			c.(*settableCounterVec).Set(completed, "completed", hostname)
			// set this explicitly to be more obvious
			died := removals - completed
			c.(*settableCounterVec).Set(died, "died", hostname)
			return nil
		},
		counter("master", "slave_removal_events_reasons", "Total number of slave removal events by reason on this master since it booted.", "reason", "hostname"): func(m metricMap, c prometheus.Collector) error {
			re, err := regexp.Compile("master/slave_removals/reason_(.*?)$")
			if err != nil {
				log.WithFields(log.Fields{
					"regex":  "master/slave_removals/reason_(.*?)$",
					"metric": "master_slave_removal_events_reasons",
					"error":  err,
				}).Error("could not compile regex")
				return fmt.Errorf("could not compile slave_removal_events_reasons regex: %s", err)
			}
			for metric, value := range m {
				matches := re.FindStringSubmatch(metric)
				if len(matches) != 2 {
					continue
				}
				reason := matches[1]
				c.(*settableCounterVec).Set(value, reason, hostname)
			}
			return nil
		},
		counter("master", "slave_unreachable_events_total", "Total number of slave unreachable events on this master since it booted.", "event", "hostname"): func(m metricMap, c prometheus.Collector) error {
			canceled, ok := m["master/slave_unreachable_canceled"]
			if !ok {
				log.WithField("metric", "master/slave_unreachable_canceled").Warn(LogErrNotFoundInMap)
			}
			completed, ok := m["master/slave_unreachable_completed"]
			if !ok {
				log.WithField("metric", "master/slave_unreachable_completed").Warn(LogErrNotFoundInMap)
			}
			scheduled, ok := m["master/slave_unreachable_scheduled"]
			if !ok {
				log.WithField("metric", "master/slave_unreachable_scheduled").Warn(LogErrNotFoundInMap)
			}
			c.(*settableCounterVec).Set(canceled, "canceled", hostname)
			c.(*settableCounterVec).Set(completed, "completed", hostname)
			c.(*settableCounterVec).Set(scheduled, "scheduled", hostname)
			return nil
		},

		gauge("master", "slaves_state", "Current number of slaves known to the master per connection and registration state.", "state", "hostname"): func(m metricMap, c prometheus.Collector) error {
			active, ok := m["master/slaves_active"]
			if !ok {
				log.WithField("metric", "master/slaves_active").Warn(LogErrNotFoundInMap)
			}
			inactive, ok := m["master/slaves_inactive"]
			if !ok {
				log.WithField("metric", "master/slaves_inactive").Warn(LogErrNotFoundInMap)
			}
			disconnected, ok := m["master/slaves_disconnected"]
			if !ok {
				log.WithField("metric", "master/slaves_disconnected").Warn(LogErrNotFoundInMap)
			}
			unreachable, ok := m["master/slaves_unreachable"]
			if !ok {
				log.WithField("metric", "master/slaves_unreachable").Warn(LogErrNotFoundInMap)
			}

			// FIXME: Make sure those assumptions are right
			// Every "active" node is connected to the master
			c.(*prometheus.GaugeVec).WithLabelValues("connected_active", hostname).Set(active)
			// Every "inactive" node is connected but node sending offers
			c.(*prometheus.GaugeVec).WithLabelValues("connected_inactive", hostname).Set(inactive)
			// Every "disconnected" node is "inactive"
			c.(*prometheus.GaugeVec).WithLabelValues("disconnected_inactive", hostname).Set(disconnected)
			// Every "connected" node is either active or inactive
			c.(*prometheus.GaugeVec).WithLabelValues("unreachable", hostname).Set(unreachable)
			return nil
		},

		// Master stats about frameworks
		gauge("master", "frameworks_state", "Current number of frames known to the master per connection and registration state.", "state", "hostname"): func(m metricMap, c prometheus.Collector) error {
			active, ok := m["master/frameworks_active"]
			if !ok {
				log.WithField("metric", "master/frameworks_active").Warn(LogErrNotFoundInMap)
			}
			inactive, ok := m["master/frameworks_inactive"]
			if !ok {
				log.WithField("metric", "master/frameworks_inactive").Warn(LogErrNotFoundInMap)
			}
			disconnected, ok := m["master/frameworks_disconnected"]
			if !ok {
				log.WithField("metric", "master/frameworks_disconnected").Warn(LogErrNotFoundInMap)
			}
			// FIXME: Make sure those assumptions are right
			// Every "active" framework is connected to the master
			c.(*prometheus.GaugeVec).WithLabelValues("connected_active", hostname).Set(active)
			// Every "inactive" framework is connected but framework sending offers
			c.(*prometheus.GaugeVec).WithLabelValues("connected_inactive", hostname).Set(inactive)
			// Every "disconnected" framework is "inactive"
			c.(*prometheus.GaugeVec).WithLabelValues("disconnected_inactive", hostname).Set(disconnected)
			// Every "connected" framework is either active or inactive
			return nil
		},
		prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "mesos",
			Subsystem: "master",
			Name:      "offers_pending",
			Help:      "Current number of offers made by the master which aren't yet accepted or declined by frameworks.",
		}, []string{"hostname"}): func(m metricMap, c prometheus.Collector) error {
			offers, ok := m["master/outstanding_offers"]
			if !ok {
				log.WithField("metric", "master/outstanding_offers").Warn(LogErrNotFoundInMap)
			}
			c.(*prometheus.GaugeVec).WithLabelValues(hostname).Set(offers)
			return nil
		},

		// Master stats about tasks
		counter("master", "task_states_exit_total", "Total number of tasks processed by exit state.", "state", "hostname"): func(m metricMap, c prometheus.Collector) error {
			dropped, ok := m["master/tasks_dropped"]
			if !ok {
				log.WithField("metric", "master/tasks_dropped").Warn(LogErrNotFoundInMap)
			}
			errored, ok := m["master/tasks_error"]
			if !ok {
				log.WithField("metric", "master/tasks_error").Warn(LogErrNotFoundInMap)
			}
			failed, ok := m["master/tasks_failed"]
			if !ok {
				log.WithField("metric", "master/tasks_failed").Warn(LogErrNotFoundInMap)
			}
			finished, ok := m["master/tasks_finished"]
			if !ok {
				log.WithField("metric", "master/tasks_finished").Warn(LogErrNotFoundInMap)
			}
			gone, ok := m["master/tasks_gone"]
			if !ok {
				log.WithField("metric", "master/tasks_gone").Warn(LogErrNotFoundInMap)
			}
			goneByOperator, ok := m["master/tasks_gone_by_operator"]
			if !ok {
				log.WithField("metric", "master/tasks_gone_by_operator").Warn(LogErrNotFoundInMap)
			}
			killed, ok := m["master/tasks_killed"]
			if !ok {
				log.WithField("metric", "master/tasks_killed").Warn(LogErrNotFoundInMap)
			}
			killing, ok := m["master/tasks_killing"]
			if !ok {
				log.WithField("metric", "master/tasks_killing").Warn(LogErrNotFoundInMap)
			}
			lost, ok := m["master/tasks_lost"]
			if !ok {
				log.WithField("metric", "master/tasks_lost").Warn(LogErrNotFoundInMap)
			}
			c.(*settableCounterVec).Set(dropped, "dropped", hostname)
			c.(*settableCounterVec).Set(errored, "errored", hostname)
			c.(*settableCounterVec).Set(failed, "failed", hostname)
			c.(*settableCounterVec).Set(finished, "finished", hostname)
			c.(*settableCounterVec).Set(gone, "gone", hostname)
			c.(*settableCounterVec).Set(goneByOperator, "gone_by_operator", hostname)
			c.(*settableCounterVec).Set(killed, "killed", hostname)
			c.(*settableCounterVec).Set(killing, "killing", hostname)
			c.(*settableCounterVec).Set(lost, "lost", hostname)
			return nil
		},

		counter("master", "task_states_current", "Current number of tasks by state.", "state", "hostname"): func(m metricMap, c prometheus.Collector) error {
			running, ok := m["master/tasks_running"]
			if !ok {
				log.WithField("metric", "master/tasks_running").Warn(LogErrNotFoundInMap)
			}
			staging, ok := m["master/tasks_staging"]
			if !ok {
				log.WithField("metric", "master/tasks_staging").Warn(LogErrNotFoundInMap)
			}
			starting, ok := m["master/tasks_starting"]
			if !ok {
				log.WithField("metric", "master/tasks_starting").Warn(LogErrNotFoundInMap)
			}
			unreachable, ok := m["master/tasks_unreachable"]
			if !ok {
				log.WithField("metric", "master/tasks_unreachable").Warn(LogErrNotFoundInMap)
			}
			c.(*settableCounterVec).Set(running, "running", hostname)
			c.(*settableCounterVec).Set(staging, "staging", hostname)
			c.(*settableCounterVec).Set(starting, "starting", hostname)
			c.(*settableCounterVec).Set(unreachable, "unreachable", hostname)
			return nil
		},

		counter("master", "task_state_counts_by_source_reason", "Number of task states by source and reason", "state", "source", "reason", "hostname"): func(m metricMap, c prometheus.Collector) error {
			re, err := regexp.Compile("master/task_(.*?)/source_(.*?)/reason_(.*?)$")
			if err != nil {
				log.WithFields(log.Fields{
					"regex":  "master/task_(.*?)/source_(.*?)/reason_(.*?)$",
					"metric": "master_task_state_counts_by_source_reason",
					"error":  err,
				}).Error("could not compile regex")
				return fmt.Errorf("could not compile task_state_counts regex: %s", err)
			}
			for metric, value := range m {
				matches := re.FindStringSubmatch(metric)
				if len(matches) != 4 {
					continue
				}
				state := matches[1]
				source := matches[2]
				reason := matches[3]
				c.(*settableCounterVec).Set(value, state, source, reason, hostname)
			}
			return nil
		},

		// Master stats about messages
		counter("master", "messages", "Number of messages by the master by state", "type", "hostname"): func(m metricMap, c prometheus.Collector) error {
			droppedMessages, ok := m["master/dropped_messages"]
			if !ok {
				log.WithField("metric", "master/dropped_messages").Warn(LogErrNotFoundInMap)
			}
			authenticateMessages, ok := m["master/messages_authenticate"]
			if !ok {
				log.WithField("metric", "master/messages_authenticate").Warn(LogErrNotFoundInMap)
			}
			deactivateFrameworkMessages, ok := m["master/messages_deactivate_framework"]
			if !ok {
				log.WithField("metric", "master/messages_deactivate_framework").Warn(LogErrNotFoundInMap)
			}
			declineOfferMessages, ok := m["master/messages_decline_offers"]
			if !ok {
				log.WithField("metric", "master/messages_decline_offers").Warn(LogErrNotFoundInMap)
			}
			executorToFrameworkMessages, ok := m["master/messages_executor_to_framework"]
			if !ok {
				log.WithField("metric", "master/messages_executor_to_framework").Warn(LogErrNotFoundInMap)
			}
			exitedExecutor, ok := m["master/messages_exited_executor"]
			if !ok {
				log.WithField("metric", "master/messages_exited_executor").Warn(LogErrNotFoundInMap)
			}
			frameworkToExecutor, ok := m["master/messages_framework_to_executor"]
			if !ok {
				log.WithField("metric", "master/messages_framework_to_executor").Warn(LogErrNotFoundInMap)
			}
			killTask, ok := m["master/messages_kill_task"]
			if !ok {
				log.WithField("metric", "master/messages_kill_task").Warn(LogErrNotFoundInMap)
			}
			launchTasks, ok := m["master/messages_launch_tasks"]
			if !ok {
				log.WithField("metric", "master/messages_launch_tasks").Warn(LogErrNotFoundInMap)
			}
			reconcileTasks, ok := m["master/messages_reconcile_tasks"]
			if !ok {
				log.WithField("metric", "master/messages_reconcile_tasks").Warn(LogErrNotFoundInMap)
			}
			registerFramework, ok := m["master/messages_register_framework"]
			if !ok {
				log.WithField("metric", "master/messages_register_framework").Warn(LogErrNotFoundInMap)
			}
			registerSlave, ok := m["master/messages_register_slave"]
			if !ok {
				log.WithField("metric", "master/messages_register_slave").Warn(LogErrNotFoundInMap)
			}
			reregisterFramework, ok := m["master/messages_reregister_framework"]
			if !ok {
				log.WithField("metric", "master/messages_reregister_framework").Warn(LogErrNotFoundInMap)
			}
			reregisterSlave, ok := m["master/messages_reregister_slave"]
			if !ok {
				log.WithField("metric", "master/messages_reregister_slave").Warn(LogErrNotFoundInMap)
			}
			resourceRequest, ok := m["master/messages_resource_request"]
			if !ok {
				log.WithField("metric", "master/messages_resource_request").Warn(LogErrNotFoundInMap)
			}
			reviveOffers, ok := m["master/messages_revive_offers"]
			if !ok {
				log.WithField("metric", "master/messages_revive_offers").Warn(LogErrNotFoundInMap)
			}
			statusUpdate, ok := m["master/messages_status_update"]
			if !ok {
				log.WithField("metric", "master/messages_status_update").Warn(LogErrNotFoundInMap)
			}
			statusUpdateAck, ok := m["master/messages_status_update_acknowledgement"]
			if !ok {
				log.WithField("metric", "master/messages_status_update_acknowledgement").Warn(LogErrNotFoundInMap)
			}
			suppressOffers, ok := m["master/messages_suppress_offers"]
			if !ok {
				log.WithField("metric", "master/messages_suppress_offers").Warn(LogErrNotFoundInMap)
			}
			unregisterFramework, ok := m["master/messages_unregister_framework"]
			if !ok {
				log.WithField("metric", "master/messages_unregister_framework").Warn(LogErrNotFoundInMap)
			}
			unregisterSlave, ok := m["master/messages_unregister_slave"]
			if !ok {
				log.WithField("metric", "master/messages_unregister_slave").Warn(LogErrNotFoundInMap)
			}
			updateSlave, ok := m["master/messages_update_slave"]
			if !ok {
				log.WithField("metric", "master/messages_update_slave").Warn(LogErrNotFoundInMap)
			}

			c.(*settableCounterVec).Set(authenticateMessages, "authenticate_messages", hostname)
			c.(*settableCounterVec).Set(droppedMessages, "dropped_messages", hostname)
			c.(*settableCounterVec).Set(deactivateFrameworkMessages, "deactivate_framework", hostname)
			c.(*settableCounterVec).Set(declineOfferMessages, "decline_offers", hostname)
			c.(*settableCounterVec).Set(executorToFrameworkMessages, "executor_to_framework", hostname)
			c.(*settableCounterVec).Set(exitedExecutor, "exited_executor", hostname)
			c.(*settableCounterVec).Set(frameworkToExecutor, "framework_to_executor", hostname)
			c.(*settableCounterVec).Set(killTask, "kill_task", hostname)
			c.(*settableCounterVec).Set(launchTasks, "launch_tasks", hostname)
			c.(*settableCounterVec).Set(reconcileTasks, "reconcile_tasks", hostname)
			c.(*settableCounterVec).Set(registerFramework, "register_framework", hostname)
			c.(*settableCounterVec).Set(registerSlave, "register_slave", hostname)
			c.(*settableCounterVec).Set(reregisterFramework, "reregister_framework", hostname)
			c.(*settableCounterVec).Set(reregisterSlave, "reregister_slave", hostname)
			c.(*settableCounterVec).Set(resourceRequest, "resource_request", hostname)
			c.(*settableCounterVec).Set(reviveOffers, "revive_offers", hostname)
			c.(*settableCounterVec).Set(statusUpdate, "status_update", hostname)
			c.(*settableCounterVec).Set(statusUpdateAck, "status_update_acknowledgement", hostname)
			c.(*settableCounterVec).Set(suppressOffers, "suppress_offers", hostname)
			c.(*settableCounterVec).Set(unregisterFramework, "unregister_framework", hostname)
			c.(*settableCounterVec).Set(unregisterSlave, "unregister_slave", hostname)
			c.(*settableCounterVec).Set(updateSlave, "update_slave", hostname)
			return nil
		},

		counter("master", "messages_outcomes_total",
			"Total number of messages by outcome of operation and direction.",
			"source", "destination", "type", "outcome", "hostname"): func(m metricMap, c prometheus.Collector) error {
			frameworkToExecutorValid, ok := m["master/valid_framework_to_executor_messages"]
			if !ok {
				log.WithField("metric", "master/valid_framework_to_executor_messages").Warn(LogErrNotFoundInMap)
			}
			frameworkToExecutorInvalid, ok := m["master/invalid_framework_to_executor_messages"]
			if !ok {
				log.WithField("metric", "master/invalid_framework_to_executor_messages").Warn(LogErrNotFoundInMap)
			}
			executorToFrameworkValid, ok := m["master/valid_executor_to_framework_messages"]
			if !ok {
				log.WithField("metric", "master/valid_executor_to_framework_messages").Warn(LogErrNotFoundInMap)
			}
			executorToFrameworkInvalid, ok := m["master/invalid_executor_to_framework_messages"]
			if !ok {
				log.WithField("metric", "master/invalid_executor_to_framework_messages").Warn(LogErrNotFoundInMap)
			}

			// status updates are sent from framework?(FIXME) to slave
			// status update acks are sent from slave to framework?
			statusUpdateAckValid, ok := m["master/valid_status_update_acknowledgements"]
			if !ok {
				log.WithField("metric", "master/valid_status_update_acknowledgements").Warn(LogErrNotFoundInMap)
			}
			statusUpdateAckInvalid, ok := m["master/invalid_status_update_acknowledgements"]
			if !ok {
				log.WithField("metric", "master/invalid_status_update_acknowledgements").Warn(LogErrNotFoundInMap)
			}
			statusUpdateValid, ok := m["master/valid_status_updates"]
			if !ok {
				log.WithField("metric", "master/valid_status_updates").Warn(LogErrNotFoundInMap)
			}
			statusUpdateInvalid, ok := m["master/invalid_status_updates"]
			if !ok {
				log.WithField("metric", "master/invalid_status_updates").Warn(LogErrNotFoundInMap)
			}
			c.(*settableCounterVec).Set(frameworkToExecutorValid, "framework", "executor", "", "valid", hostname)
			c.(*settableCounterVec).Set(frameworkToExecutorInvalid, "framework", "executor", "", "invalid", hostname)

			c.(*settableCounterVec).Set(executorToFrameworkValid, "executor", "framework", "", "valid", hostname)
			c.(*settableCounterVec).Set(executorToFrameworkInvalid, "executor", "framework", "", "invalid", hostname)

			// We consider a ack message simply as a message from slave to framework
			c.(*settableCounterVec).Set(statusUpdateValid, "framework", "slave", "status_update", "valid", hostname)
			c.(*settableCounterVec).Set(statusUpdateInvalid, "framework", "slave", "status_update", "invalid", hostname)
			c.(*settableCounterVec).Set(statusUpdateAckValid, "slave", "framework", "status_update", "valid", hostname)
			c.(*settableCounterVec).Set(statusUpdateAckInvalid, "slave", "framework", "status_update", "invalid", hostname)
			return nil
		},
		// Master stats about events
		gauge("master", "event_queue_length", "Current number of elements in event queue by type", "type", "hostname"): func(m metricMap, c prometheus.Collector) error {
			dispatches, ok := m["master/event_queue_dispatches"]
			if !ok {
				log.WithField("metric", "master/event_queue_dispatches").Warn(LogErrNotFoundInMap)
			}
			httpRequests, ok := m["master/event_queue_http_requests"]
			if !ok {
				log.WithField("metric", "master/event_queue_http_requests").Warn(LogErrNotFoundInMap)
			}
			messages, ok := m["master/event_queue_messages"]
			if !ok {
				log.WithField("metric", "master/event_queue_messages").Warn(LogErrNotFoundInMap)
			}
			c.(*prometheus.GaugeVec).WithLabelValues("message", hostname).Set(messages)
			c.(*prometheus.GaugeVec).WithLabelValues("http_request", hostname).Set(httpRequests)
			c.(*prometheus.GaugeVec).WithLabelValues("dispatches", hostname).Set(dispatches)
			return nil
		},

		// Master stats about allocations
		prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "mesos",
			Subsystem: "master",
			Name:      "allocator_event_queue_dispatches",
			Help:      "Number of dispatch events in the allocator event queue.",
		}, []string{"hostname"}): func(m metricMap, c prometheus.Collector) error {
			count, ok := m["allocator/event_queue_dispatches"]
			if !ok {
				log.WithField("metric", "allocator/event_queue_dispatches").Warn(LogErrNotFoundInMap)
			}
			c.(*prometheus.GaugeVec).WithLabelValues(hostname).Set(count)
			return nil
		},

		prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "mesos",
			Subsystem: "master",
			Name:      "allocation_run_ms_count",
			Help:      "Number of allocation algorithm time measurements in the window",
		}, []string{"hostname"}): func(m metricMap, c prometheus.Collector) error {
			count, ok := m["allocator/mesos/allocation_runs"]
			if !ok {
				log.WithField("metric", "allocator/mesos/allocation_runs").Warn(LogErrNotFoundInMap)
			}
			c.(*prometheus.GaugeVec).WithLabelValues(hostname).Set(count)
			return nil
		},

		gauge("master", "allocation_run_ms", "Time spent in allocation algorithm in ms.", "type", "hostname"): func(m metricMap, c prometheus.Collector) error {
			mean, ok := m["allocator/mesos/allocation_run_ms"]
			if !ok {
				log.WithField("metric", "allocator/mesos/allocation_run_ms").Warn(LogErrNotFoundInMap)
			}
			min, ok := m["allocator/mesos/allocation_run_ms/min"]
			if !ok {
				log.WithField("metric", "allocator/mesos/allocation_run_ms/min").Warn(LogErrNotFoundInMap)
			}
			max, ok := m["allocator/mesos/allocation_run_ms/max"]
			if !ok {
				log.WithField("metric", "allocator/mesos/allocation_run_ms/max").Warn(LogErrNotFoundInMap)
			}
			p50, ok := m["allocator/mesos/allocation_run_ms/p50"]
			if !ok {
				log.WithField("metric", "allocator/mesos/allocation_run_ms/p50").Warn(LogErrNotFoundInMap)
			}
			p90, ok := m["allocator/mesos/allocation_run_ms/p90"]
			if !ok {
				log.WithField("metric", "allocator/mesos/allocation_run_ms/p90").Warn(LogErrNotFoundInMap)
			}
			p95, ok := m["allocator/mesos/allocation_run_ms/p95"]
			if !ok {
				log.WithField("metric", "allocator/mesos/allocation_run_ms/p95").Warn(LogErrNotFoundInMap)
			}
			p99, ok := m["allocator/mesos/allocation_run_ms/p99"]
			if !ok {
				log.WithField("metric", "allocator/mesos/allocation_run_ms/p99").Warn(LogErrNotFoundInMap)
			}
			p999, ok := m["allocator/mesos/allocation_run_ms/p999"]
			if !ok {
				log.WithField("metric", "allocator/mesos/allocation_run_ms/p999").Warn(LogErrNotFoundInMap)
			}
			p9999, ok := m["allocator/mesos/allocation_run_ms/p9999"]
			if !ok {
				log.WithField("metric", "allocator/mesos/allocation_run_ms/p9999").Warn(LogErrNotFoundInMap)
			}
			c.(*prometheus.GaugeVec).WithLabelValues("mean", hostname).Set(mean)
			c.(*prometheus.GaugeVec).WithLabelValues("min", hostname).Set(min)
			c.(*prometheus.GaugeVec).WithLabelValues("max", hostname).Set(max)
			c.(*prometheus.GaugeVec).WithLabelValues("p50", hostname).Set(p50)
			c.(*prometheus.GaugeVec).WithLabelValues("p90", hostname).Set(p90)
			c.(*prometheus.GaugeVec).WithLabelValues("p95", hostname).Set(p95)
			c.(*prometheus.GaugeVec).WithLabelValues("p99", hostname).Set(p99)
			c.(*prometheus.GaugeVec).WithLabelValues("p999", hostname).Set(p999)
			c.(*prometheus.GaugeVec).WithLabelValues("p9999", hostname).Set(p9999)
			return nil
		},

		counter("master", "allocation_runs", "Number of times the allocation alorithm has run", "event", "hostname"): func(m metricMap, c prometheus.Collector) error {
			runs, ok := m["allocator/mesos/allocation_runs"]
			if !ok {
				log.WithField("metric", "allocator/mesos/allocation_runs").Warn(LogErrNotFoundInMap)
			}
			c.(*settableCounterVec).Set(runs, "allocation", hostname)
			return nil
		},

		counter("master", "allocation_run_latency_ms_count", "Number of allocation batch latency measurements", "event", "hostname"): func(m metricMap, c prometheus.Collector) error {
			count, ok := m["allocator/mesos/allocation_run_latency_ms/count"]
			if !ok {
				log.WithField("metric", "allocator/mesos/allocation_run_latency_ms/count").Warn(LogErrNotFoundInMap)
			}
			c.(*settableCounterVec).Set(count, "allocation", hostname)
			return nil
		},

		gauge("master", "allocation_run_latency_ms", "Allocation batch latency in ms.", "type", "hostname"): func(m metricMap, c prometheus.Collector) error {
			mean, ok := m["allocator/mesos/allocation_run_latency_ms"]
			if !ok {
				log.WithField("metric", "allocator/mesos/allocation_run_latency_ms").Warn(LogErrNotFoundInMap)
			}
			min, ok := m["allocator/mesos/allocation_run_latency_ms/min"]
			if !ok {
				log.WithField("metric", "allocator/mesos/allocation_run_latency_ms/min").Warn(LogErrNotFoundInMap)
			}
			max, ok := m["allocator/mesos/allocation_run_latency_ms/max"]
			if !ok {
				log.WithField("metric", "allocator/mesos/allocation_run_latency_ms/max").Warn(LogErrNotFoundInMap)
			}
			p50, ok := m["allocator/mesos/allocation_run_latency_ms/p50"]
			if !ok {
				log.WithField("metric", "allocator/mesos/allocation_run_latency_ms/p50").Warn(LogErrNotFoundInMap)
			}
			p90, ok := m["allocator/mesos/allocation_run_latency_ms/p90"]
			if !ok {
				log.WithField("metric", "allocator/mesos/allocation_run_latency_ms/p90").Warn(LogErrNotFoundInMap)
			}
			p95, ok := m["allocator/mesos/allocation_run_latency_ms/p95"]
			if !ok {
				log.WithField("metric", "allocator/mesos/allocation_run_latency_ms/p95").Warn(LogErrNotFoundInMap)
			}
			p99, ok := m["allocator/mesos/allocation_run_latency_ms/p99"]
			if !ok {
				log.WithField("metric", "allocator/mesos/allocation_run_latency_ms/p99").Warn(LogErrNotFoundInMap)
			}
			p999, ok := m["allocator/mesos/allocation_run_latency_ms/p999"]
			if !ok {
				log.WithField("metric", "allocator/mesos/allocation_run_latency_ms/p999").Warn(LogErrNotFoundInMap)
			}
			p9999, ok := m["allocator/mesos/allocation_run_latency_ms/p9999"]
			if !ok {
				log.WithField("metric", "allocator/mesos/allocation_run_latency_ms/p9999").Warn(LogErrNotFoundInMap)
			}
			c.(*prometheus.GaugeVec).WithLabelValues("mean", hostname).Set(mean)
			c.(*prometheus.GaugeVec).WithLabelValues("min", hostname).Set(min)
			c.(*prometheus.GaugeVec).WithLabelValues("max", hostname).Set(max)
			c.(*prometheus.GaugeVec).WithLabelValues("p50", hostname).Set(p50)
			c.(*prometheus.GaugeVec).WithLabelValues("p90", hostname).Set(p90)
			c.(*prometheus.GaugeVec).WithLabelValues("p95", hostname).Set(p95)
			c.(*prometheus.GaugeVec).WithLabelValues("p99", hostname).Set(p99)
			c.(*prometheus.GaugeVec).WithLabelValues("p999", hostname).Set(p999)
			c.(*prometheus.GaugeVec).WithLabelValues("p9999", hostname).Set(p9999)
			return nil
		},

		prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "mesos",
			Subsystem: "master",
			Name:      "event_queue_dispatches",
			Help:      "Number of dispatch events in the allocator mesos event queue.",
		}, []string{"hostname"}): func(m metricMap, c prometheus.Collector) error {
			count, ok := m["allocator/mesos/event_queue_dispatches"]
			if !ok {
				log.WithField("metric", "allocator/mesos/event_queue_dispatches").Warn(LogErrNotFoundInMap)
			}
			c.(*prometheus.GaugeVec).WithLabelValues(hostname).Set(count)
			return nil
		},
		gauge("master", "allocator_offer_filters_active", "Number of active offer filters for all frameworks within the role", "role", "hostname"): func(m metricMap, c prometheus.Collector) error {
			re, err := regexp.Compile("allocator/mesos/offer_filters/roles/(.*?)/active")
			if err != nil {
				log.WithFields(log.Fields{
					"regex":  "allocator/mesos/offer_filters/roles/(.*?)/active",
					"metric": "master_allocator_offer_filters_active",
					"error":  err,
				}).Error("could not compile regex")
				return fmt.Errorf("could not compile allocator_offer_filters_active regex: %s", err)
			}
			for metric, value := range m {
				matches := re.FindStringSubmatch(metric)
				if len(matches) != 2 {
					continue
				}
				role := matches[1]
				c.(*prometheus.GaugeVec).WithLabelValues(role, hostname).Set(value)
			}
			return nil
		},

		gauge("master", "allocator_role_quota_offered_or_allocated", "Amount of resources considered offered or allocated towards a role's quota guarantee.", "role", "resource", "hostname"): func(m metricMap, c prometheus.Collector) error {
			re, err := regexp.Compile("allocator/mesos/quota/roles/(.*?)/resources/(.*?)/offered_or_allocated")
			if err != nil {
				log.WithFields(log.Fields{
					"regex":  "allocator/mesos/quota/roles/(.*?)/resources/(.*?)/offered_or_allocated",
					"metric": "master_allocator_role_quota_offered_or_allocated",
					"error":  err,
				}).Error("could not compile regex")
				return fmt.Errorf("could not compile allocator_role_quota_offered_or_allocated regex: %s", err)
			}
			for metric, value := range m {
				matches := re.FindStringSubmatch(metric)
				if len(matches) != 3 {
					continue
				}
				role := matches[1]
				resource := matches[2]
				c.(*prometheus.GaugeVec).WithLabelValues(role, resource, hostname).Set(value)
			}
			return nil
		},

		gauge("master", "allocator_role_shares_dominant", "Dominance factor for a role", "role", "hostname"): func(m metricMap, c prometheus.Collector) error {
			re, err := regexp.Compile("allocator/mesos/roles/(.*?)/shares/dominant")
			if err != nil {
				log.WithFields(log.Fields{
					"regex":  "allocator/mesos/roles/(.*?)/shares/dominant",
					"metric": "master_allocator_role_shares_dominant",
					"error":  err,
				}).Error("could not compile regex")
				return fmt.Errorf("could not compile allocator_role_shares_dominant regex: %s", err)
			}
			for metric, value := range m {
				matches := re.FindStringSubmatch(metric)
				if len(matches) != 2 {
					continue
				}
				role := matches[1]
				c.(*prometheus.GaugeVec).WithLabelValues(role, hostname).Set(value)
			}
			return nil
		},

		gauge("master", "allocator_role_quota_guarantee", "Amount of resources guaranteed for a role via quota", "role", "resource", "hostname"): func(m metricMap, c prometheus.Collector) error {
			re, err := regexp.Compile("allocator/mesos/quota/roles/(.*?)/resources/(.*?)/guarantee")
			if err != nil {
				log.WithFields(log.Fields{
					"regex":  "allocator/mesos/quota/roles/(.*?)/resources/(.*?)/guarantee",
					"metric": "master_allocator_role_quota_guarantee",
					"error":  err,
				}).Error("could not compile regex")
				return fmt.Errorf("could not compile allocator_role_quota_guarantee regex: %s", err)
			}
			for metric, value := range m {
				matches := re.FindStringSubmatch(metric)
				if len(matches) != 3 {
					continue
				}
				role := matches[1]
				resource := matches[2]
				c.(*prometheus.GaugeVec).WithLabelValues(role, resource, hostname).Set(value)
			}
			return nil
		},

		gauge("master", "allocator_resources_cpus", "Number of CPUs offered or allocated", "type", "hostname"): func(m metricMap, c prometheus.Collector) error {
			total, ok := m["allocator/mesos/resources/cpus/total"]
			if !ok {
				log.WithField("metric", "allocator/mesos/resources/cpus/total").Warn(LogErrNotFoundInMap)
			}
			offeredOrAllocated, ok := m["allocator/mesos/resources/cpus/offered_or_allocated"]
			if !ok {
				log.WithField("metric", "allocator/mesos/resources/cpus/offered_or_allocated").Warn(LogErrNotFoundInMap)
			}

			c.(*prometheus.GaugeVec).WithLabelValues("total", hostname).Set(total)
			c.(*prometheus.GaugeVec).WithLabelValues("offered_or_allocated", hostname).Set(offeredOrAllocated)
			return nil
		},

		gauge("master", "allocator_resources_disk", "Allocated or offered disk space in MB", "type", "hostname"): func(m metricMap, c prometheus.Collector) error {
			total, ok := m["allocator/mesos/resources/disk/total"]
			if !ok {
				log.WithField("metric", "allocator/mesos/resources/disk/total").Warn(LogErrNotFoundInMap)
			}
			offeredOrAllocated, ok := m["allocator/mesos/resources/disk/offered_or_allocated"]
			if !ok {
				log.WithField("metric", "allocator/mesos/resources/disk/offered_or_allocated").Warn(LogErrNotFoundInMap)
			}

			c.(*prometheus.GaugeVec).WithLabelValues("total", hostname).Set(total)
			c.(*prometheus.GaugeVec).WithLabelValues("offered_or_allocated", hostname).Set(offeredOrAllocated)
			return nil
		},

		gauge("master", "allocator_resources_mem", "Allocated or offered memory in MB", "type", "hostname"): func(m metricMap, c prometheus.Collector) error {
			total, ok := m["allocator/mesos/resources/mem/total"]
			if !ok {
				log.WithField("metric", "allocator/mesos/resources/mem/total").Warn(LogErrNotFoundInMap)
			}
			offeredOrAllocated, ok := m["allocator/mesos/resources/mem/offered_or_allocated"]
			if !ok {
				log.WithField("metric", "allocator/mesos/resources/mem/offered_or_allocated").Warn(LogErrNotFoundInMap)
			}

			c.(*prometheus.GaugeVec).WithLabelValues("total", hostname).Set(total)
			c.(*prometheus.GaugeVec).WithLabelValues("offered_or_allocated", hostname).Set(offeredOrAllocated)
			return nil
		},

		// Framework call counts (total)
		counter("master", "framework_calls_total", "Counts of API calls per framework", "framework_name", "framework_id", "hostname"): func(m metricMap, c prometheus.Collector) error {
			re, _ := regexp.Compile(`master/frameworks/([^.]*)\.([^/]+)/calls`)
			for metric, value := range m {
				matches := re.FindStringSubmatch(metric)
				if len(matches) != 3 {
					continue
				}
				name := matches[1]
				id := matches[2]
				c.(*settableCounterVec).Set(value, name, id, hostname)
			}
			return nil
		},

		// Framework call counts (by type)
		counter("master", "framework_calls", "Counts of API calls per framework", "framework_name", "framework_id", "type", "hostname"): func(m metricMap, c prometheus.Collector) error {
			re, _ := regexp.Compile(`master/frameworks/([^.]*)\.([^/]+)/calls/(.+)`)
			for metric, value := range m {
				matches := re.FindStringSubmatch(metric)
				if len(matches) != 4 {
					continue
				}
				name := matches[1]
				id := matches[2]
				typ := matches[3]
				c.(*settableCounterVec).Set(value, name, id, typ, hostname)
			}
			return nil
		},

		// Framework offer operation counts (total)
		counter("master", "framework_operations_total", "Counts of offer operations per framework", "framework_name", "framework_id", "hostname"): func(m metricMap, c prometheus.Collector) error {
			re, _ := regexp.Compile(`master/frameworks/([^.]*)\.([^/]+)/operations`)
			for metric, value := range m {
				matches := re.FindStringSubmatch(metric)
				if len(matches) != 3 {
					continue
				}
				name := matches[1]
				id := matches[2]
				c.(*settableCounterVec).Set(value, name, id, hostname)
			}
			return nil
		},

		// Framework offer operation counts (by type)
		counter("master", "framework_operations", "Counts of offer operations per framework", "framework_name", "framework_id", "type", "hostname"): func(m metricMap, c prometheus.Collector) error {
			re, _ := regexp.Compile(`master/frameworks/([^.]*)\.([^/]+)/operations/(.+)`)
			for metric, value := range m {
				matches := re.FindStringSubmatch(metric)
				if len(matches) != 4 {
					continue
				}
				name := matches[1]
				id := matches[2]
				typ := matches[3]
				c.(*settableCounterVec).Set(value, name, id, typ, hostname)
			}
			return nil
		},

		// Framework subscribed
		gauge("master", "framework_subscribed", "Boolean: is this framework subscribed?", "framework_name", "framework_id", "hostname"): func(m metricMap, c prometheus.Collector) error {
			re, _ := regexp.Compile(`master/frameworks/([^.]*)\.([^/]+)/subscribed`)
			for metric, value := range m {
				matches := re.FindStringSubmatch(metric)
				if len(matches) != 2 {
					continue
				}
				name := matches[1]
				id := matches[2]
				c.(*prometheus.GaugeVec).WithLabelValues(name, id, hostname).Set(value)
			}
			return nil
		},

		// Framework role state
		gauge("master", "framework_suppressed", "Boolean: are offers for this role suppressed?", "framework_name", "framework_id", "role", "hostname"): func(m metricMap, c prometheus.Collector) error {
			re, _ := regexp.Compile(`master/frameworks/([^.]*)\.([^/]+)/roles/([^/]+)/suppressed`)
			for metric, value := range m {
				matches := re.FindStringSubmatch(metric)
				if len(matches) != 4 {
					continue
				}
				name := matches[1]
				id := matches[2]
				role := matches[3]
				c.(*prometheus.GaugeVec).WithLabelValues(name, id, role, hostname).Set(value)
			}
			return nil
		},

		// Framework events (total)
		counter("master", "framework_events_total", "Counts of events per framework", "framework_name", "framework_id", "hostname"): func(m metricMap, c prometheus.Collector) error {
			re, _ := regexp.Compile(`master/frameworks/([^.]*)\.([^/]+)/events`)
			for metric, value := range m {
				matches := re.FindStringSubmatch(metric)
				if len(matches) != 3 {
					continue
				}
				name := matches[1]
				id := matches[2]
				c.(*settableCounterVec).Set(value, name, id, hostname)
			}
			return nil
		},

		// Framework events (by type)
		counter("master", "framework_events", "Counts of events per framework", "framework_name", "framework_id", "type", "hostname"): func(m metricMap, c prometheus.Collector) error {
			re, _ := regexp.Compile(`master/frameworks/([^.]*)\.([^/]+)/events/([^/]+)`)
			for metric, value := range m {
				matches := re.FindStringSubmatch(metric)
				if len(matches) != 4 {
					continue
				}
				name := matches[1]
				id := matches[2]
				typ := matches[3]
				c.(*settableCounterVec).Set(value, name, id, typ, hostname)
			}
			return nil
		},

		// Framework update events (by state)
		counter("master", "framework_update_events", "Counts of update events per framework", "framework_name", "framework_id", "state", "hostname"): func(m metricMap, c prometheus.Collector) error {
			re, _ := regexp.Compile(`master/frameworks/([^.]*)\.([^/]+)/events/update/(.+)`)
			for metric, value := range m {
				matches := re.FindStringSubmatch(metric)
				if len(matches) != 4 {
					continue
				}
				name := matches[1]
				id := matches[2]
				state := matches[3]
				c.(*settableCounterVec).Set(value, name, id, state, hostname)
			}
			return nil
		},

		// Framework task active states
		gauge("master", "framework_task_states", "Task states per framework", "framework_name", "framework_id", "state", "hostname"): func(m metricMap, c prometheus.Collector) error {
			re, _ := regexp.Compile(`master/frameworks/([^.]*)\.([^/]+)/tasks/([^/]+ing|task_lost)`)
			for metric, value := range m {
				matches := re.FindStringSubmatch(metric)
				if len(matches) != 3 {
					continue
				}
				name := matches[1]
				id := matches[2]
				state := matches[3]

				c.(*prometheus.GaugeVec).WithLabelValues(name, id, state, hostname).Set(value)
			}
			return nil
		},

		// Framework terminal task states
		counter("master", "framework_terminal_task_states", "Terminal task states per framework", "framework_name", "framework_id", "state", "hostname"): func(m metricMap, c prometheus.Collector) error {
			// TODO: don't use two regexes here
			re, _ := regexp.Compile(`master/frameworks/([^.]*)\.([^/]+)/tasks/([^/]+)`)
			notre, _ := regexp.Compile(`master/frameworks/([^.]*)\.([^/]+)/tasks/([^/]+ing|task_lost)`)

			for metric, value := range m {
				matches := re.FindStringSubmatch(metric)
				notmatches := notre.FindStringSubmatch(metric)
				if len(matches) != 4 || len(notmatches) == 4 {
					continue
				}
				name := matches[1]
				id := matches[2]
				state := matches[3]

				c.(*settableCounterVec).Set(value, name, id, state, hostname)
			}
			return nil
		},

		// Framework task failures by source reason
		counter("master", "framework_task_failures", "Framework task failures by source reason", "framework_name", "framework_id", "source", "reason", "hostname"): func(m metricMap, c prometheus.Collector) error {
			re, _ := regexp.Compile(`master/frameworks/([^.]*)\.([^/]+)/tasks/task_failed/([^/]+)/([^/]+)`)
			for metric, value := range m {
				matches := re.FindStringSubmatch(metric)
				if len(matches) != 5 {
					continue
				}
				name := matches[1]
				id := matches[2]
				source := matches[3]
				reason := matches[4]

				c.(*settableCounterVec).Set(value, name, id, source, reason, hostname)
			}
			return nil
		},

		// Framework offers
		counter("master", "framework_offers", "Number of offers by type per framework", "framework_name", "framework_id", "type", "hostname"): func(m metricMap, c prometheus.Collector) error {
			re, _ := regexp.Compile(`master/frameworks/([^.]*)\.([^/]+)/offers/(.+)`)
			for metric, value := range m {
				matches := re.FindStringSubmatch(metric)
				if len(matches) != 4 {
					continue
				}
				name := matches[1]
				id := matches[2]
				typ := matches[3]
				c.(*settableCounterVec).Set(value, name, id, typ, hostname)
			}
			return nil
		},

		// Per-framework resource allocation: offer filters
		counter("master", "framework_offer_filters", "Number of filters set per framework over period", "framework_name", "framework_id", "period", "hostname"): func(m metricMap, c prometheus.Collector) error {
			re, _ := regexp.Compile(`master/frameworks/([^.]*)\.([^/]+)/allocation/offer_filters/refused_seconds/(.+)`)
			for metric, value := range m {
				matches := re.FindStringSubmatch(metric)
				if len(matches) != 4 {
					continue
				}
				name := matches[1]
				id := matches[2]
				pd := matches[3]
				c.(*settableCounterVec).Set(value, name, id, pd, hostname)
			}
			return nil
		},

		// Per-framework resource filters
		counter("master", "framework_resource_filters", "Number of times resources were filtered per framework", "framework_name", "framework_id", "reason", "hostname"): func(m metricMap, c prometheus.Collector) error {
			re, _ := regexp.Compile(`master/frameworks/([^.]*)\.([^/]+)/allocation/resources_filtered/(.+)`)
			for metric, value := range m {
				matches := re.FindStringSubmatch(metric)
				if len(matches) != 4 {
					continue
				}
				name := matches[1]
				id := matches[2]
				reason := matches[3]
				c.(*settableCounterVec).Set(value, name, id, reason, hostname)
			}
			return nil
		},

		// Latest per-role DRF position
		gauge("master", "frameworks_drf_position", "Latest per-role DRF position", "framework_name", "framework_id", "role", "type", "hostname"): func(m metricMap, c prometheus.Collector) error {
			re, _ := regexp.Compile(`master/frameworks/([^.]*)\.([^/]+)/allocation/roles/([^/]+)/latest_position/(.+)`)
			for metric, value := range m {
				matches := re.FindStringSubmatch(metric)
				if len(matches) != 5 {
					continue
				}
				name := matches[1]
				id := matches[2]
				role := matches[3]
				typ := matches[4]
				c.(*prometheus.GaugeVec).WithLabelValues(name, id, role, typ, hostname).Set(value)
			}
			return nil
		},

		// Framework message metrics
		counter("master", "frameworks_messages", "Messages passed around with the frameworks", "framework", "type", "hostname"): func(m metricMap, c prometheus.Collector) error {
			re, err := regexp.Compile("frameworks/(.*?)/messages_(.*?)$")
			if err != nil {
				log.WithFields(log.Fields{
					"regex":  "frameworks/(.*?)/messages_(.*?)$",
					"metric": "master_frameworks_messages",
					"error":  err,
				}).Error("could not compile regex")
				return fmt.Errorf("could not compile frameworks_messages regex: %s", err)
			}
			for metric, messageCount := range m {
				matches := re.FindStringSubmatch(metric)
				if len(matches) != 3 {
					continue
				}
				framework := matches[1]
				messageStatus := matches[2]
				if len(framework) > 0 && len(messageStatus) > 0 {
					c.(*settableCounterVec).Set(messageCount, framework, messageStatus, hostname)
				}
			}
			return nil
		},

		// Registrar stats
		prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "mesos",
			Subsystem: "registrar",
			Name:      "registry_size_bytes",
			Help:      "Size of the registry in bytes",
		}, []string{"hostname"}): func(m metricMap, c prometheus.Collector) error {
			size, ok := m["registrar/registry_size_bytes"]
			if !ok {
				log.WithField("metric", "registrar/registry_size_bytes").Warn(LogErrNotFoundInMap)
			}
			c.(*prometheus.GaugeVec).WithLabelValues(hostname).Set(size)
			return nil
		},
		prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "mesos",
			Subsystem: "registrar",
			Name:      "queued_operations",
			Help:      "Number of operations in the registry queue",
		}, []string{"hostname"}): func(m metricMap, c prometheus.Collector) error {
			ops, ok := m["registrar/queued_operations"]
			if !ok {
				log.WithField("metric", "registrar/queued_operations").Warn(LogErrNotFoundInMap)
			}
			c.(*prometheus.GaugeVec).WithLabelValues(hostname).Set(ops)
			return nil
		},
		prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "mesos",
			Subsystem: "registrar",
			Name:      "state_fetch_ms",
			Help:      "Duration of state JSON fetch in ms",
		}, []string{"hostname"}): func(m metricMap, c prometheus.Collector) error {
			ms, ok := m["registrar/state_fetch_ms"]
			if !ok {
				log.WithField("metric", "registrar/state_fetch_ms").Warn(LogErrNotFoundInMap)
			}
			c.(*prometheus.GaugeVec).WithLabelValues(hostname).Set(ms)
			return nil
		},
		gauge("registrar", "state_store_ms", "Duration of state json store in ms.", "type", "hostname"): func(m metricMap, c prometheus.Collector) error {
			mean, ok := m["registrar/state_store_ms"]
			if !ok {
				log.WithFields(log.Fields{
					"name": "registrar/state_store_ms",
				}).Warn(LogErrNotFoundInMap)
			}
			min, ok := m["registrar/state_store_ms/min"]
			if !ok {
				log.WithFields(log.Fields{
					"name": "registrar/state_store_ms/min",
				}).Warn(LogErrNotFoundInMap)
			}
			max, ok := m["registrar/state_store_ms/max"]
			if !ok {
				log.WithFields(log.Fields{
					"name": "registrar/state_store_ms/max",
				}).Warn(LogErrNotFoundInMap)
			}
			p50, ok := m["registrar/state_store_ms/p50"]
			if !ok {
				log.WithFields(log.Fields{
					"name": "registrar/state_store_ms/p50",
				}).Warn(LogErrNotFoundInMap)
			}
			p90, ok := m["registrar/state_store_ms/p90"]
			if !ok {
				log.WithFields(log.Fields{
					"name": "registrar/state_store_ms/p90",
				}).Warn(LogErrNotFoundInMap)
			}
			p95, ok := m["registrar/state_store_ms/p95"]
			if !ok {
				log.WithFields(log.Fields{
					"name": "registrar/state_store_ms/p95",
				}).Warn(LogErrNotFoundInMap)
			}
			p99, ok := m["registrar/state_store_ms/p99"]
			if !ok {
				log.WithFields(log.Fields{
					"name": "registrar/state_store_ms/p99",
				}).Warn(LogErrNotFoundInMap)
			}
			p999, ok := m["registrar/state_store_ms/p999"]
			if !ok {
				log.WithFields(log.Fields{
					"name": "registrar/state_store_ms/p999",
				}).Warn(LogErrNotFoundInMap)
			}
			p9999, ok := m["registrar/state_store_ms/p9999"]
			if !ok {
				log.WithFields(log.Fields{
					"name": "registrar/state_store_ms/p9999",
				}).Warn(LogErrNotFoundInMap)
			}
			c.(*prometheus.GaugeVec).WithLabelValues("mean", hostname).Set(mean)
			c.(*prometheus.GaugeVec).WithLabelValues("min", hostname).Set(min)
			c.(*prometheus.GaugeVec).WithLabelValues("max", hostname).Set(max)
			c.(*prometheus.GaugeVec).WithLabelValues("p50", hostname).Set(p50)
			c.(*prometheus.GaugeVec).WithLabelValues("p90", hostname).Set(p90)
			c.(*prometheus.GaugeVec).WithLabelValues("p95", hostname).Set(p95)
			c.(*prometheus.GaugeVec).WithLabelValues("p99", hostname).Set(p99)
			c.(*prometheus.GaugeVec).WithLabelValues("p999", hostname).Set(p999)
			c.(*prometheus.GaugeVec).WithLabelValues("p9999", hostname).Set(p9999)
			return nil
		},
		prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "mesos",
			Subsystem: "registrar",
			Name:      "log_recovered",
			Help:      "Recovered status of the registrar log",
		}, []string{"hostname"}): func(m metricMap, c prometheus.Collector) error {
			recovered, ok := m["registrar/log/recovered"]
			if !ok {
				log.WithField("metric", "registrar/log/recovered").Warn(LogErrNotFoundInMap)
			}
			c.(*prometheus.GaugeVec).WithLabelValues(hostname).Set(recovered)
			return nil
		},
		prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "mesos",
			Subsystem: "registrar",
			Name:      "log_ensemble_size",
			Help:      "Ensemble size of the registrar log",
		}, []string{"hostname"}): func(m metricMap, c prometheus.Collector) error {
			size, ok := m["registrar/log/ensemble_size"]
			if !ok {
				log.WithField("metric", "registrar/log/ensemble_size").Warn(LogErrNotFoundInMap)
			}
			c.(*prometheus.GaugeVec).WithLabelValues(hostname).Set(size)
			return nil
		},

		// Overlay log
		prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "mesos",
			Subsystem: "overlay",
			Name:      "log_recovered",
			Help:      "Recovered status of the overlay log",
		}, []string{"hostname"}): func(m metricMap, c prometheus.Collector) error {
			recovered, ok := m["overlay/log/recovered"]
			if !ok {
				log.WithField("metric", "overlay/log/recovered").Warn(LogErrNotFoundInMap)
			}
			c.(*prometheus.GaugeVec).WithLabelValues(hostname).Set(recovered)
			return nil
		},
		prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "mesos",
			Subsystem: "overlay",
			Name:      "log_ensemble_size",
			Help:      "Ensemble size of the overlay log",
		}, []string{"hostname"}): func(m metricMap, c prometheus.Collector) error {
			size, ok := m["overlay/log/ensemble_size"]
			if !ok {
				log.WithFields(log.Fields{
					"name": "overlay/log_ensemble_size",
				}).Warn(LogErrNotFoundInMap)
			}
			c.(*prometheus.GaugeVec).WithLabelValues(hostname).Set(size)
			return nil
		},

		// END
	}
	return newMetricCollector(httpClient, metrics)
}

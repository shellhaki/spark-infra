# spark-infra

SparkCell is a Linux database cell for running isolated database workloads behind a small client/daemon control plane.

The first target is PostgreSQL. Cells get their own filesystem, process space, network namespace, and resource limits, with requests routed through a Unix socket to a persistent daemon.

See [plan.MD](plan.MD) for the long-term roadmap and implementation checklist.

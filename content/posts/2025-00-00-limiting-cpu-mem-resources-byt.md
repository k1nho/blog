---
title: "Byte Size Series - Limiting CPU and Memory usage with Cgroups"
pubDate: 2025-10-08
Categories: ["DevOps", "Platform Engineering", "Byte Wisdom"]
Tags: ["Byte Wisdom", "Cgroups", "Linux", "DevOps", "Learning"]
cover: "what_is_cidr_cover.png"
mermaid: true
---

The rise of containers and its adoption across the industry with technologies like Docker and Kubernetes has been great. Below those abstractions it sits the fundamental idea that
enables to start a Linux process while limiting its CPU and Memory resources with the use of [**Control Groups**](https://man7.org/linux/man-pages/man7/cgroups.7.html), also know as **Cgroups**. In this entry of the byte size series, we'll see how can we use [systemd](https://systemd.io/)
one of the most used init systems in many Linux distros to configure and add a process to a Cgroup.

What is a Cgroup? From the [man pages](https://man7.org/linux/man-pages/man7/cgroups.7.html), we get the following:

> Control groups, usually referred to as cgroups, are a Linux kernel
> feature which allow processes to be organized into hierarchical
> groups whose usage of various types of resources can then be
> limited and monitored. The kernel's cgroup interface is provided
> through a pseudo-filesystem called cgroupfs. Grouping is
> implemented in the core cgroup kernel code, while resource
> tracking and limits are implemented in a set of per-resource-type
> subsystems (memory, CPU, and so on).

From the definition, we can see that the way to work with the cgroup interface is via the `pseudo-filesystem cgroupfs` which is usually mounted at `/sys/fs/cgroup`. If we wanted to create a cgroup,
we would simply just create a sub-directory which then gets populated with files used to manipulate the cgroup configuration itself. Let's go one level above this abstraction into using [systemd](https://systemd.io/) to
control a cgroup to make the work easier, in fact, systemd performs operations under the pseudo-filesystem as if you where using it through shell commands itself. As an example creating the cgroup
would be running the `mkdir` command on `/sys/fs/cgroup`.

To demonstrate Cgroups, we will use `spin_loop.py` this is a simple program that spins forever adding more memory on each iteration.

```python {filename="spin_loop.py"}
data = []
while true:
    data.append([0] * 10**6)
```

Let us know add a this process to a Cgroup via `systemd-run`

```bash
systemd-run -u eatmem -p CPUQuota=20% -p MemoryMax=100M ~/spin_loop.py
```

In the above, we create a systemd unit called **eatmem**, and add the **CPU quota to be 20% and the Max Memory usage to be 100 megabytes**, lastly we specify the process that we want to add to the cgroup to be
our `spin_loop.py` program. Whenever, the program breaks those limits (such as in this example where we add memory infinitely) the **out-of-memory killer** gets called. Just like that we where able to limit
the process from going wild, very cool!

There is one more thing, if we ever wanted to save this configuration, that is a Cgroup that monitors and limits the CPU quota and memory max to those that we specify, we would need to define a `slice`
as otherwise the cgroup would be tied to the process that was invoked in it. We can achieve that with the following.

```plaintext {filename="sliceconfig"}
[Slice]
CPUQuota=20%
MemoryMax=100M
```

```bash
cat sliceconfig > /etc/systemd/system/eatmem.slice
```

Now we can place any process into our `eatmem.slice` as we did before.

```bash
systemd-run -u eatmem --slice=eatmem.slice ~/spin_loop.py
```

If we add more processes, then the cumulative CPU and Memory usage of those will be limited to the ones defined in our configuration.

**Cgroups** underpins one of the most, if not the most, powerful and important technology in containerization. The management of `cgroupfs` via **systemd** explained is yet another
lower abstraction, albeit important to understand, as container managers such as [Containerd](https://containerd.io/) leverage **systemd** a cgroup **driver**
to control container resource usage allowing a system to be protected against hungry processes and distribute compute fairly.

Did you know about **Cgroups**?

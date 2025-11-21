---
title: "Kinho's Homelab Series - GitOps and First Application (ArgoCD)"
pubDate: 2026-01-12
Description: "Let's build a mini homelab! In this entry, we move into the GitOps workflow with ArgoCD and run our first app!"
Categories: ["DevOps", "Networking", "Platform Engineering", "Homelab Series"]
Tags: ["DevOps", "Homelab Series", "Networking"]
cover: "homelabs_cover2.png"
mermaid: true
draft: true
---

Welcome to another entry in the **Kinho's Homelab Series**, in the last [entry]() we setup our orchestration platform with K3s and solidified our network
stack with Cilium. However, so far we have done cluster setup but there is no apps running, more importantly our installations, such as Cilium, have been manual
with no organization. It's time to adopt the GitOps workflow!

In this entry, I will introduce [Argo CD]() as our continuous delivery solution for our Kubernetes cluster, and while we are at it
install our first application in the cluster, take a drink and let's get going!

---

# What is Argo CD ?

From the official [Argo CD documentation ](), we get a very simple definition:

> Argo CD is a declarative, GitOps continuous delivery tool for Kubernetes.

The key part here is in **GitOps**, so we must define what it is. From [Gitlab's post on GitOps]:

> GitOps is an operational framework that applies DevOps practices like version control, collaboration, and CI/CD to infrastructure automation, ensuring consistent, repeatable deployments.
> Teams put GitOps into practice by using Git repositories as the single source of truth, automating deployments, and enforcing changes through merge requests or pull requests.
> Any configuration drift, such as manual changes or errors, is overwritten by GitOps automation so the environment converges on the desired state defined in Git.

In Layman's term, just as how we use version control for application code and merge changes that fix bugs or implement a feature, we can do so too with our infrastructure.
In this case, every pull request represents a change in our infrastructure, and then a _tool_ runs in a **reconciliation loop** to make sure the cluster state matches the desired state defined in **Git**.

The tool that we will use to achieve our **GitOps workflow** for the homelab will indeed be **Argo**!

---

# Installing Argo CD

---

# Lots of Research, Lots of Links

---

# Wrapping up

---

# Next in Kinho's Homelab Series

**TBD**

# Resources

- [Argo CD](https://docs.k3s.io/)
- [What is GitOps ?](https://about.gitlab.com/topics/gitops/)
- [Cilium](https://cilium.io/)

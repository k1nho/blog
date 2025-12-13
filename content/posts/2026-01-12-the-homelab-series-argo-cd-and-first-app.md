---
title: "Kinho's Homelab Series - GitOps and First Application"
pubDate: 2026-01-12
Description: "Let's build a mini homelab! In this entry, we move into the GitOps workflow with ArgoCD and deploy our first app!"
Categories: ["DevOps", "Networking", "Platform Engineering", "Homelab Series"]
Tags: ["DevOps", "Homelab Series", "Networking"]
cover: "homelabs_cover3.png"
mermaid: true
draft: true
---

Welcome to another entry in the **Kinho's Homelab Series**, in the last [entry]() we setup our orchestration platform with K3s and solidified our network
stack with Cilium. However, so far we have done cluster setup but there is no apps running, more importantly our installations, such as Cilium, have been manual
with no organization. Overtime, centralizing our configuration and knowing exactly what we have running in the cluster becomes crucial to tame the unwieldy complexity of **yaml hell**.
So, It's finally time to adopt the GitOps workflow!

In this entry, I will introduce [Argo CD](https://argo-cd.readthedocs.io/en/stable/) as our continuous delivery solution for our Kubernetes cluster. We will also setup our **external secrets operator (ESO)**
and the **Tailscale operator** to end up deploying our first app. Go ahead and take a drink, this will be a fun one!

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

Installing Argo is as simple as applying its CRD's with the following command:

```bash

```

We need to pick up the initial admin password created, so that we can login into the dashboard, and change it.

```bash

```

From then on, we can access the argocd-server for now lets expose it locally via port-forward.

![ArgoCD UI Login]()

## Bootstrapping the Cluster

The Argo UI is indeed very nice, however, using it defeats the purpose of setting up a declarative GitOps workflow, that is, we want to use Git as the centralized source
of truth of the cluster, and such that any changes to the cluster is logically related to a commit. We let then Argo sync the cluster using our repository to apply the
necessary resources.

---

# Strategy for Managing Secrets

Many of the applications we will run require [secrets](), however as you might have already noticed we are managing our cluster publicly! and we will
like not to be another one in the list of the [39 million secrets leaked on Github](https://resources.github.com/enterprise/understanding-secret-leak-exposure/) 😅.
There are two popular choices to manage secrets within a GitOps workflow either with a Secrets Operations Encryption [(SOPS)](), or with an External Secret Operator [(ESO)]().

The **SOPS** approach encrypts secrets that you can push into the repository such that decryption happens only within the cluster. This avoids the common pitfall of plaintext leaked secrets. Some of
the best choices for this approach are [Sealed Secrets](http://github.com/bitnami-labs/sealed-secrets), and [age](https://github.com/FiloSottile/age). On the other hand, we have the **ESO** approach
in which we pull secrets from an external manager **Azure Vault**, **AWS Secret Manager**, or **GCP Secret Manager** and sync them into Kubernetes secrets. The main idea here is to usually, store
[CRDs]() references to actual secrets, which once again avoids plaintext secrets and keeps the repository secret-free. So which one should we choose?

## Infisical Secrets Operator

Initially, I considered using Bitnami's sealed secrets which is simple enough to setup; however, there's a few things that made me actually choose a **ESO**, namely
**secret rotation** and API based secret management. With **sealed secrets**, we would need to re-encrypt our secrets every single time we want to rotate the current secrets
this is not too bad, but it becomes a bit manual. Moreover, having an API to create, fetch, and renovate secrets becomes incredibly important in CI/CD pipelines. While,
we have usual suspects to choose from, I discovered [Infisical](https://infisical.com/) a poweful open source all in one secret management platform.

Through this series, I like to consider this three as my bastions for choosing a particular software: **the project is open source, has a generous free tier, and can, if one chooses to, be self-hosted!**.
Infisical meets this criteria, and when I consider that they have both a [Kubernetes Operator](https://infisical.com/docs/integrations/platforms/kubernetes/overview) and an [SDK](https://infisical.com/docs/sdks/overview)
it ends up fullfilling the other requirements: sync secrets into our cluster, and manage them with an API. Let's define our Argo App to deploy the infisical secret operator.

```yaml {filename="argo/infisical.yaml"}
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: infisical-secrets-operator
  namespace: argocd
  annotations:
    argocd.argoproj.io/sync-wave: "0" # cluster level priority on sync (https://argo-cd.readthedocs.io/en/stable/user-guide/sync-waves/)
spec:
  project: default
  sources:
    - repoURL: https://dl.cloudsmith.io/public/infisical/helm-charts/helm/charts/
      chart: secrets-operator
      targetRevision: 0.10.3
      helm:
        valueFiles:
          - $values/cluster/infisical-operator/infisical-values.yaml
    - repoURL: https://github.com/k1nho/homelab
      targetRevision: main
      ref: values
  destination:
    server: https://kubernetes.default.svc
    namespace: infisical-secrets-operator
  syncPolicy:
    automated:
      prune: true
      selfHeal: true
    syncOptions:
      - CreateNamespace=true
```

Let's test this by setting up the **InfisicalSecret CRD** forthe oauth secret that we will need in order to setup our **tailscale operator**.

---

# Tailscale Operator

The [tailscale kubernetes operator](https://tailscale.com/kb/1236/kubernetes-operator) enables:

- Securing access to the Kubernetes control plane
- Exposing cluster workloads to the tailnet (Ingress)
- Exposing a tailnet service to the Kubernetes cluster (Egress)

There are many more possibilities as listed in [the official documentation](https://tailscale.com/kb/1236/kubernetes-operator), for this case we are interested
in **exposing our cluster workloads to the tailnet**, that is, using it as ingress. As usual, we will configure our tailscale operator as an ArgoCD app to deploy
the Helm chart into the cluster with some custom values.

```yaml {filename="argo/tailscale.yaml"}
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: tailscale-operator
  namespace: argocd
  annotations:
    argocd.argoproj.io/sync-wave: "1" # app level priority on sync (https://argo-cd.readthedocs.io/en/stable/user-guide/sync-waves/)
spec:
  project: default
  sources:
    # Kustomize
    - repoURL: https://github.com/k1nho/homelab
      targetRevision: main
      path: apps/tailscale-operator

    - repoURL: https://pkgs.tailscale.com/helmcharts
      chart: tailscale-operator
      targetRevision: 1.90.9
      helm:
        valueFiles:
          - $values/apps/tailscale-operator/tailscale-operator-values.yaml

    - repoURL: https://github.com/k1nho/homelab
      targetRevision: main
      ref: values

  destination:
    server: https://kubernetes.default.svc
    namespace: tailscale-operator
  syncPolicy:
    automated:
      prune: true
      selfHeal: true
    syncOptions:
      - CreateNamespace=true
```

## Adding the Oauth Secret

Notice we define the path `apps/tailscale-operator` this will pick up our `kustomization.yaml` definition, that packs the Infisical Secret CRD `operator-oauth-secret.yaml`:

```yaml
apiVersion: secrets.infisical.com/v1alpha1
kind: InfisicalSecret
metadata:
  name: operator-oauth
spec:
  authentication:
    universalAuth:
      secretsScope:
        projectSlug: homelab-d-s7-g

        envSlug: "prod"
        secretsPath: "/tailscale"
      credentialsRef:
        secretName: universal-auth-credentials
        secretNamespace: infisical-secrets-operator

  managedKubeSecretReferences:
    - secretName: operator-oauth
      secretNamespace: tailscale-operator
```

---

# The First Application

We are ready to deploy the first application into our cluster. For the first app, we are not going with anything crazy, but something simple enough
to demonstrate that our cluster is running properly to deploy services. We will go for a classic [stateless application](https://kubernetes.io/docs/tutorials/stateless-application/),
and what a better one than this blog itself!

## Setting Up a CI Pipeline with Github Actions and Dagger

To be able to define declarative our workload, we will need to create some Kubernetes resources that we will then push into our repository for Argo to sync. More importantly,
we do not have an image for the blog, so let's get started by creating a pipeline that will build our image when we push a tag in the blog repository.

### Dagger Module

## Exposing the Blog

---

# Wrapping up

That's it for this entry! we started by migrating our existing Cilium deployment into a declarative GitOps workflow which led us to the introduction of ArgoCD. From there
we adopted Argo's App of Apps pattern and setup our Infisical as our secret management solution, and the tailscale operator. Lastly, we deployed the blog as the first
application, and expose it with the Tailscale Ingress!

The cluster is now alive with the blog running! but there are many more improvements needed. First, we have deployed multiple applications but we have to monitor effectively
the resource compsumption and traffic of the different services in our cluster. Moreover, often times applications need a database to **persistently store information**, so a solution
for provisioning storage, backup, and a disaster recovery strategy becomes important to keep data safe. In the next entry, we'll explore a few of these!

---

# Next in Kinho's Homelab Series

**TBD**

# Resources

- [What is GitOps ?](https://about.gitlab.com/topics/gitops/)
- [Argo CD Documentation](https://docs.k3s.io/)
- [Infisical Kubernetes Operator](https://infisical.com/docs/integrations/platforms/kubernetes/overview)
- [Tailscale Kubernetes Operator](https://tailscale.com/kb/1236/kubernetes-operator)
- [Cilium](https://cilium.io/)

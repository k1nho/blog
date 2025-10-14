# The Kinho's Homelab Series (Securing Your Network with Tailscale)

This year I want to really up my DevOps skills to the next level, as it is an area I am very interested in. After having such a great experience, with my dive in Kubernetes services for the fantastic [NSDF Intersect](https://nationalsciencedatafabric.org/nsdf-intersect) project,
I decided that is now time to finally go full throttle into a mini homelab. Usually when people think of a homelab, they sort of envision this huge server rack with boxes upon boxes upon boxes... and with that a money tree 😁. However, If you are like me, chances are you want
to build a test lab, one that can serve you to experiment different technologies, and not a production grade cluster. The plan is to start from just a barebones Linux distro to running my own Kubernetes cluster. Along the way, I might look into
replacing Google photos with something like [Immich](https://immich.app/) or hosting a personal media server with [Plex](https://www.plex.tv/personal-media-server/) possibilities are truly endless in the world of homelab.

In this series, I'll document my hardware, decision, and progress of the **Kinho's Homelab** adding piece by piece my knowledge gains to it. This sure will be fun! (and we shall know [pain](https://www.youtube.com/shorts/xu7X5-5U-b0)).
For this first article in the series, I go about setting up my first node, and securing my devices under a tailnet with wireguard using [Tailscale](https://tailscale.com/).

# Grandma's Laptop Has a Use!

As we alluded to earlier, the point of this is to build a homelab that will work as an environment for us to experiment, that means your old grandma's laptop will work too. For me that one is an extremely old **Sony VAIO VPCEJ** i3 4GB RAM and 500GB SSD.
This machine is looking at the light ready to go barely taking the bloatware of windows. So, in order to give new life to this poor soul, we will revitalize this with Linux.

# Ubuntu btw ?

As deciding a distro to flash into the VAIO, I of course consider [Arch](https://archlinux.org/) btw, but it is fine, I'll take on the long bootcut pants for a moment and go for Ubuntu. Simple fact is that I actually do not mind Ubuntu, in fact any Linux but Windows is 10^6 better, and it is stable
the theme here is that we want this environment to be solid, even though our main objective is experimentation, we do not want to compromise for breaks during different setups. Don't get me wrong I will give Arch its shot in the future, but for now is Ubuntu btw.

If you have never flashed an OS in your life, you can follow this simple steps to do [install Ubuntu](https://ubuntu.com/tutorials/install-ubuntu-desktop#1-overview), keep in mind that the image used in the tutorial is for a Desktop environment and that is fine if you intend to use the laptop
as your daily driver. However, you can go for the [Ubuntu server] image which is more akin for future setups.

# Laptop Heavy, Knees Weak, Arms Heavy

Ok, so the first issue I wanted to solve is this idea of accessing my VAIO through my Macbook via [SSH](https://en.wikipedia.org/wiki/Secure_Shell) since I will be out of my local network somedays, and just setting up the whole ordeal with authorized keys would be a hassle.
Investigating more about how can I go about achieving, I found the fantastic [Tailscale](https://tailscale.com/) and just wow I was amazed to what it had in store for this.

# All This For Free

**Tailscale** works as a way to make secure networking for point to point connectivity between devices using [wireguard](https://www.wireguard.com/) in what is know as a **tailnet**.
That means, that I can connect from my Macbook to my VAIO with ssh completely private and most importantly securily. What is even better is that the [personal](https://tailscale.com/pricing?plan=personal) tier
is incredibly generous with up to 100 devices and 3 users. Not to glaze, but the setup was incredibly easy which is always something that I like to sing praises
whenever I try a new product.

For setting up my Macbook, I simply download the [tailscale client](https://tailscale.com/download/mac) from the official site, and continue the setup from there.
For my VAIO running linux to join the mesh it was as simple as running the following:

```
curl -fsSL https://tailscale.com/install.sh | sh
```

I am very exited to test out other features that Tailscale enables such as [serve](https://tailscale.com/kb/1312/serve) for routing traffic from my devices
to a local service running on my tailnet, and also [funnels](https://tailscale.com/kb/1223/funnel) to route traffic from the internet to a local service on my tailnet.
For now, I have achieve my goal, I can now connect from my Macbook to my VAIO via ssh privately and securely. There is one more thing to take care of to finish
securing the machine.

# Raise the wall

Finally, we will do a few things to raise the wall. First, let's setup our sshd_config to only listen to its private IP
defined in Tailscale.

```plaintext filename={"sshd_config"}
ListenAddress <your-tailscale-ip>
```

Also, we will configure `ufw` to allow access to SSH only from Tailscale IP [CIDR](https://k1nho.github.io/blog/posts/what-is-cidr/) range with the following
command.

```bash
sudo ufw allow from '100.64.0.0/10' to any port 22
```

Configuring for IPV6 as well.

```bash
sudo ufw allow from 'fda7:115c:a1e0::/48' to any port 22
```

Lastly, we deny access from any other IP to SSH.

```bash
sudo ufw deny 22
```

```bash
sudo ufw deny OpenSSH
```

Lastly, let's enable the firewall.

```bash
sudo ufw enable
```

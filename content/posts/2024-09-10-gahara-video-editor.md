---
title: "Gahara: a Vim Video Editor"
pubDate: 2024-09-10
Description: "I had done websites, backends, but there was one alley unexplored desktop apps. In this article, I talk about Gahara my small video editor created with Wails"
Categories: ["Desktop", "Go", "Video Editor", "ffmpeg"]
Tags: ["Story", "Go", "Svelte", "Desktop Dev", "WailsApp"]
cover: "gahara_cover.png"
mermaid: true
draft: true
---

In the summer of 2023, I did a small [contribution](https://github.com/wailsapp/wails/pull/2812) to a very interesting project by the name of [Wails](https://wails.io/).
This open source project boasting over [21,000](https://github.com/wailsapp/wails) stars, had me captivated with a simple fact, you could combine Go and any frontend technology to create
a cross platform desktop application. It's unified [eventing system](https://wails.io/docs/reference/runtime/events) sounded great for this crazy idea I had, call it a bit of a challenge.
**Could I create a simple video editor?**.

There are many full blown video editors that are open source out there such as [Kdenlive](https://kdenlive.org/) and [OpenShot](https://www.openshot.org/), but what I wanted to aim for
was more in the vein of the fantastic [Losslesscut](https://mifi.no/losslesscut/) project by Mikael Finstad, aka [mifi](https://github.com/mifi). What is the main difference of a full blown video editor vs a lossless cut editor?
the key is in re-encoding.

```mermaid
flowchart LR
    c1[" file.mp4"]
    subgraph Transformation
    c2["Trim Video Operation <br> No Re-encode"]
    end
    subgraph Result
    c3["cut_file.mp4"]
    end

    c1 --> c2
    c2 --> c3
```

```mermaid
flowchart LR
    c1[" file.mp4"]
    subgraph Transformation
    c2["Add Text Operation <br> Re-encoding"]
    end
    subgraph Result
    c3["text_added_file.mp4"]
    end

    c1 --> c2
    c2 --> c3
```

Links/Resources:

- [Gahara](https://github.com/Gahara-Editor/gahara)
- [Wails](https://wails.io/)

---
title: "Why Your CPU Is Fast But Your Program Is Slow: Understanding the Memory Wall"
date: "2026-04-18"
tags: [systems, performance, memory, cache, architecture, low-level]
description: "An exploration of why fast CPUs still run slow programs, uncovering the memory wall through experiments, cache behavior, and data movement."
permalink: posts/{{ title | slug }}/index.html
author_name: M Pranav
author_link: "https://github.com/pranav0x0112"
---

My laptop's CPU can do billions of operations per second. I know this because 
the spec sheet told me, and I believed it, because I am a trusting person.

So when I wrote a program to scan a 1GB array and it took 400 milliseconds, 
I was confused. That's not billions of anything. That's just... slow. 
Embarrassingly slow. The kind of slow that makes you question your life choices.

The CPU wasn't the problem. It was sitting there, starving, waiting for data 
that memory couldn't deliver fast enough. This gap between how fast your CPU 
can work and how fast memory can feed it has a name: the Memory Wall. And once 
you see it, you can't unsee it.

So I built a small framework called [Aletheia](https://github.com/pranav0x0112/Aletheia) 
to understand this properly. I ran some experiments, expected boring gradual 
results, and instead got a performance cliff so sharp it looked like a bug. 
It was not a bug. It was the hardware telling me exactly how it works - I just 
had not been listening.

---

## The Illusion of Compute

Here is something nobody tells you when you first learn programming: your CPU 
is almost never the reason your program is slow.

Modern processors are genuinely hard to wrap your head around. Your CPU is not 
sitting there executing one instruction, then the next, then the next, like a 
student working through a problem set. It is simultaneously predicting which 
instructions are coming several steps ahead, executing multiple instructions 
in parallel, and quietly reordering operations in the background to avoid 
sitting idle - all without you writing a single line to ask it to. This 
happens billions of times a second, every second, without you ever thinking 
about it.

A single integer addition on modern hardware takes roughly 1 nanosecond. In 
that same nanosecond, light travels about 30 centimeters. Your CPU is doing 
math at a speed that makes the laws of physics mildly uncomfortable.

So when your program is slow, the CPU has in all likelihood already done its 
part and is now waiting. The real question is not how to make your processor 
compute faster - it is why your processor is not getting the data it needs 
fast enough to keep up. That question is what the rest of this blog is trying 
to answer.

---

## The Memory Wall

CPUs and DRAM have been improving since the 1980s, but they have not been 
improving at the same rate, and that difference has quietly become one of 
the biggest problems in systems performance.

Processors got dramatically faster over the decades - smaller transistors 
meant higher clock speeds, and smarter microarchitecture meant each clock 
cycle did more useful work. DRAM improved too, but mostly in terms of how 
much data it could store rather than how quickly it could hand that data 
over. The underlying physics of how DRAM is built puts a ceiling on how fast 
it can respond to a memory request, and that ceiling has not moved anywhere 
near as fast as CPU speeds have.

![Memory Wall](https://i.postimg.cc/sXks9vqD/Prawns-2026-04-18-16-42-48.png)
> Figure: CPU performance vs DRAM performance over time.  
> Source: Computer Architecture: A Quantitative Approach, Hennessy & Patterson.

By the mid-90s, researchers were already writing about this and calling it 
the Memory Wall. The concern was not complicated - if the time it takes to 
fetch data from memory keeps growing relative to how fast the CPU runs, 
then it does not matter how many transistors you add to the processor side 
because the processor will just spend more and more of its time sitting idle, 
waiting. That chart above shows fairly clearly that the concern was justified.

Today the gap is somewhere between 50x to 100x in the worst case. The CPU 
is fast. Getting data to the CPU is not.

But saying "DRAM is slow" is a bit unsatisfying without understanding *why* 
it is slow. So before we talk about how hardware tries to work around this 
problem, it is worth taking a few minutes to look at what is physically 
happening inside a DRAM chip every time your program asks for a value.

---

## How DRAM Actually Works

DRAM stores each bit of data as a charge in a tiny capacitor - a charged 
capacitor represents a logic `1`, a discharged one represents a logic `0`. Billions of these 
capacitors are arranged in a grid of rows and columns inside each DRAM bank, 
and reading even a single byte from this grid involves more steps than you 
might expect.

![DRAM Bank Structure](https://i.postimg.cc/mgGrp0HS/dram-block-diagram.png)
> Figure: DRAM bank organization showing row activation into sense amplifiers and subsequent column selection for data access.  
> Source: [Branch Education](https://youtu.be/7J7X7aZvMXQ?si=-_zbUvDVD-Avn3aR).

When your CPU requests a memory address, the memory controller first sends a 
row address to the DRAM. This triggers row activation - the entire row 
corresponding to that address gets read out onto a set of sense amplifier 
lines. Think of it like pulling an entire filing cabinet drawer open just 
to get one piece of paper. Only after the whole row is sitting in the sense 
amplifiers can the column address come in to select the specific chunk of 
data you actually asked for.

![dram.png](https://i.postimg.cc/jdqTWtPF/dram.png)
> Figure: An entire DRAM row being read into sense amplifiers before a 
> single byte can be accessed.  
> Source: [Branch Education](https://youtu.be/7J7X7aZvMXQ?si=-_zbUvDVD-Avn3aR).

The slow part is what happens between two row accesses. Before a new row 
can be activated, the sense amplifier lines have to be reset back to a 
neutral voltage - this step is called precharge, and it takes a fixed 
amount of time that you simply cannot skip. Sequential reads that stay 
within the same row avoid this cost because the row stays active. But 
the moment your access pattern jumps to a different row, you pay the full 
precharge and activation penalty all over again.

![DRAM Cycle](https://i.postimg.cc/GpmTzZWq/t-RAS-full.png)
> Figure: DRAM read cycle showing precharge, row activation and column 
> selection.  
> Source: [Branch Education](https://youtu.be/7J7X7aZvMXQ?si=-_zbUvDVD-Avn3aR).

Random memory access patterns are slow largely because of this - every 
jump to a new address potentially means a new row activation, which means 
another precharge cycle, which means your CPU waits. This is not something 
you can fix in software. It is a physical constraint baked into how DRAM 
is built, and it is exactly the kind of behavior that showed up in the 
Aletheia experiments.

---

## The Cache Hierarchy

So if DRAM is this slow, how does anything run at a reasonable speed at all? 
The answer is that hardware engineers knew about this problem long before it 
became a crisis, and they built a workaround directly into the processor: 
a set of smaller, faster memory banks that sit between the CPU and DRAM, 
called caches.

The idea is simple. Instead of going all the way to DRAM every time the CPU 
needs data, keep a copy of recently used data closer to the processor where 
it can be accessed much faster. If the CPU asks for something and it is already 
in the cache, great - no trip to DRAM needed. If it is not, you go fetch it 
from DRAM and bring a copy back into the cache for next time.

Most modern processors have three levels of this, called L1, L2, and L3 (L here stands for Level).

![Cache Hierarchy](https://i.postimg.cc/RF7FmmnX/cache.png)
> Figure: Memory hierarchy showing cache levels, sizes and access speeds.  
> Source: Computer Architecture: A Quantitative Approach, Hennessy & Patterson.

L1 is the smallest and fastest - on my Ryzen 9 5900HX it is 32KB per core, 
and it sits so close to the execution units that it can respond in about 4 
clock cycles. L2 is larger at 512KB per core but a bit slower. L3 is shared 
across all cores, 16MB on this chip, and slower still. And then beyond L3 
is main memory - DRAM, which is enormous but comes with all the latency 
baggage we just talked about.

<table style="border-collapse: collapse; width: 100%; font-size: 0.95em;">
  <thead>
    <tr style="border-bottom: 1px solid #444;">
      <th align="left">Memory Level</th>
      <th align="left">Size (5900HX)</th>
      <th align="left">Latency</th>
    </tr>
  </thead>
  <tbody>
    <tr><td>L1 Cache</td><td>32 KB/core</td><td>~4 cycles</td></tr>
    <tr><td>L2 Cache</td><td>512 KB/core</td><td>~12 cycles</td></tr>
    <tr><td>L3 Cache</td><td>16 MB shared</td><td>~40 cycles</td></tr>
    <tr><td>DRAM</td><td>Main memory</td><td>~200+ cycles</td></tr>
  </tbody>
</table>

That jump from L3 to DRAM is not a typo. You are looking at a 5x latency 
penalty the moment your data is not in any cache level and has to be fetched 
from main memory. This is the cliff we keep talking about - not a gradual 
slope but a sudden drop the moment you fall out of cache.

The cache hierarchy buys back a lot of the performance that the memory wall 
takes away, but only if your program is actually using data in a way that 
lets the cache do its job. And that is where things get interesting - and 
where Aletheia comes in.

---

## The Experiment

Aletheia started as a simple question: if the memory wall is real and 
measurable, can I actually see it happen on my own machine with my own 
experiments rather than just taking a textbook's word for it?

The framework is straightforward. There is a host - a CLI tool that runs 
experiments and collects results. And there is a node - a server that holds 
datasets in memory and executes operations on them. Right now both run on the 
same machine, but the eventual goal is to run the node on a Raspberry Pi to 
simulate what actual near-memory compute looks like. 

> [!Fun Fact]  
> The name Aletheia comes from the Greek word for truth or disclosure - which felt appropriate for a 
> project whose entire purpose is to make the hardware tell you what it is 
> actually doing.

I ran several experiments -- a basic dataset scan, a vector addition workload, 
and the one that ended up being the most interesting: a stride scan.

The stride scan works like this. Instead of reading through an array 
sequentially, you jump through it at a fixed interval called a stride. A 
stride of 1 means you read every element. A stride of 4 means you read every 
fourth element. A stride of 4096 means you are making very large jumps through 
memory. The access pattern looks like this:

```rust
index = (k * stride) % N
```

The dataset was fixed at 256MB throughout. I tested strides of 1, 4, 16, 64, 
256, and 4096, and measured how long each took. My expectation going in was 
that performance would degrade gradually and smoothly as the stride increased 
-- bigger jumps, slightly worse performance, nothing dramatic.

That is not what happened.
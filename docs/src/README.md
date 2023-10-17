# Introduction

Guardian is designed to solve many of the problems with running a suite of security tools in a traditional CI/CD
pipeline like GitHub, GitLab, Jenkins, etc.
These tools are fantastic for what is typically defined as CI/CD.
They are flexible enough to run almost anything using complex execution flow concepts like stages, rules, workflows, 
etc.

In platform engineering where there are multiple applications being developed by multiple teams, having a security 
pipeline is essential.
What tends to happen is a pipeline or DevOps team will build CI/CD configuration templates that can be consumed and 
customized by the Dev Teams to fit their needs.
Pipelines are great for that because they are basically used for arbitrary code execution in a slightly more secure way
than running commands on a non-ephemeral server.

The natural thought process when the need for security controls arise is to just run the tools in the CI/CD pipeline.
This is a great strategy when there are only a handful of projects but one that can quickly get out of hand and turn 
into what I call _template hell_.
You have to accommodate every edge, every permutation, for every conceivable situation and application stack.
Centralizing the CI/CD configurations keeps code duplication down and provides a single place to modify thing, but 
it inevitably turns a mess of variable overrides, custom solutions, and edge cases.
Part of the reason is that most configuration files are defined in something like YAML which... has its problems to say
the least.

Guardians seeks to provide an alternative solution specifically for the security portions of the pipeline with the 
end goal of becoming a Security Pipeline as a Service.
However, in the interim there are many problems that can be solved with our existing solutions.
YAML makes it easy to read and customize the pipeline and execution flow but gets messy once you're in the script 
block.
The current solution is to go fetch a Bash wizard if you want to customize how the command is run.

Guardian CLI's `exec` command takes advantage of the same templating language used in tools like Helm for Helm charts.
It's incredibly flexible and expressive which gives developers the ability to define commands in a more declarative 
fashion than the traditional script block.
It's also a bit more readable than bash, and you can test and validate your assumptions locally without having to run
a job and read through the pipeline log.

See the reference documentation for command usage and how it can dramatically improve the pipeline development process.
This is just one of many utilities planned for Guardian so stay tuned for future features!

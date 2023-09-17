# Guardian

An effortless Security Pipeline.

## Origin

CI/CD Pipelines are easy with the multitude of tools and platforms out there like GitHub Actions, GitLab CI/CD,
Circle CI, and of course he who shall not be named.
Essentially these tools are arbitrary code execution engines.
You tell it what you want to run, how it should be run, and define the environment in which the jobs will be run.

This is a great model for building, unit testing, etc., things that can vary greatly depending on the technology
stack your team has decided on.
CI/CD Pipelines are to the developer community what Excel is to offices around the world;
A tool that is powerful enough to handle pretty much any task.
It's a hammer at your fingertips and with a hammer, everything begins to look like a nail.

This is not a great model for a security pipeline for several reasons, see
[Issues With Security Pipelines](issues-with-security-pipelines.md).
Most security pipelines have predefined requirements with tools that aren't opinionated about the technology stack.
Static code analysis, container build, virus scan, SBOM, vulnerability analysis, and threshold verification should
run exactly the same (give or take a few parameters), regardless of the application.

Guardian aims to be that security pipeline with the ability to run anywhere as a single binary.
It won't be recreated the wheel, there are far smarter people in the FOSS space solving problems like package
dependency vulnerability analysis.
Guardian will provide a framework to define the _what_ not the _how_, leveraging WASM in a plugin-like fashion.

## Concepts

Guardian will have a single command:

```shell
guardian run
```

Using a simple configuration file, it will run the security pipeline on the target code base.

### Pipeline

| Tasks                  | Default Tool                                             |
| ---------------------- | -------------------------------------------------------- |
| Static Code Analysis   | [Semgrep](https://github.com/returntocorp/semgrep)       |
| Container Build        | [Kaniko](https://github.com/GoogleContainerTools/kaniko) |
| Virus Scan             | [ClamAV](https://github.com/Cisco-Talos/clamav)          |
| SBOM Generation        | [Syft](https://github.com/anchore/syft)                  |
| Vulnerability Analysis | [Grype](https://github.com/anchore/grype)                |
| Threshold Validation   | [Gatecheck](https://github/com/gatecheckdev/gatecheck)   |

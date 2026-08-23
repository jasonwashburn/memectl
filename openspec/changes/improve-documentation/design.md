## Context

See `proposal.md` for motivation. The current README combines end-user instructions with the local toolchain and validation workflow required to develop memectl. Releases now publish `memectl` archives for macOS and Linux on amd64 and arm64, so the user-facing guide can document the supported artifact distribution directly.

## Goals / Non-Goals

**Goals:**

- Give first-time users a clear, release-based path from download to a working `memectl` command.
- Separate contributor workflow from product usage without losing project context.
- Retain the project’s OpenSpec experiment framing in both documents at an appropriate level of detail.
- Keep the writing concise, accurate, and lightly meme-aware where it fits naturally.

**Non-Goals:**

- Add an installation script, package-manager distribution, Windows support, or release signing.
- Change the CLI, release pipeline, toolchain, or OpenSpec workflow.
- Turn the documentation into a Kubernetes parody; the kubectl-inspired command structure remains an interface convention only.

## Decisions

### Document manual archive installation as the supported user path

The README will direct users to download the release archive matching their operating system and architecture, extract it, and place `memectl` in a directory on `PATH`. It will verify the result with `memectl --version` before showing CLI examples.

This matches the artifacts already published without creating a dependency on GitHub CLI or unimplemented package-manager or shell-installer infrastructure. The instructions will name only the currently released macOS and Linux architectures and link to Releases rather than hard-coding a versioned asset URL.

### Split documentation by reader intent

The README will contain the product overview, release installation, command examples, credential guidance, and a link to `CONTRIBUTING.md`. `CONTRIBUTING.md` will contain mise setup, hook installation, quality checks, source execution through `go run`, and local release snapshots.

Keeping source execution exclusively in contributor documentation prevents `go run` from being presented as the normal way to use a released CLI, while preserving all existing development guidance.

### Mention OpenSpec briefly in both documents

The README will preserve a short statement that the project is an experiment in OpenSpec and spec-driven development. The contributing guide will explain that this experiment is a driving project goal, without attempting to teach or prescribe the OpenSpec workflow.

This gives colleagues and prospective contributors the intended context while avoiding duplicated process documentation that can become stale.

### Use restrained, meme-native personality

Documentation must prioritize scanability and commands users can copy. Humor may appear where it arises from the tool’s meme domain, including existing imagery and suitable example content, but must not be forced into instructional copy. Kubernetes terminology and cluster jokes are excluded despite the kubectl-inspired command grammar.

Regular-text mentions of `memectl` in `README.md` and `CONTRIBUTING.md` will use inline code formatting. Headings and image alt text are excluded where code formatting would reduce readability or be unsupported.

## Risks / Trade-offs

- [Manual installation has more steps than a package manager] -> State the archive and `PATH` steps plainly, and limit claims to the distribution methods that exist.
- [Release instructions can become inaccurate when supported targets change] -> Describe the currently configured macOS/Linux amd64/arm64 support and link to Releases for current artifacts.
- [Splitting docs can hide essential contributor information] -> Link prominently from the README and retain focused contributor command examples.
- [Humor can make instructions unclear or dated] -> Keep headings and procedural steps literal; use personality only where it does not alter meaning.
